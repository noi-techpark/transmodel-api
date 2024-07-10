// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package provider

import (
	"fmt"
	"log/slog"
	"opendatahub/sta-nap-export/config"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"
	"opendatahub/sta-nap-export/siri"
	"strings"

	"golang.org/x/exp/maps"
)

type OdhParkingEcharging struct {
	Scode       string
	Sname       string
	Sorigin     string
	Scoordinate struct {
		X    float32
		Y    float32
		Srid uint32
	}
	Smetadata struct {
		State    string
		Capacity int32
	}
}

type ParkingEcharging struct{}

func (ParkingEcharging) origins() string {
	origins := []string{
		"ALPERIA",
		"route220",
		"DRIWE",
	}
	quoted := []string{}
	for _, o := range origins {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", o))
	}
	return strings.Join(quoted, ",")
}

func (p ParkingEcharging) odhStatic() ([]OdhParkingEcharging, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"EChargingStation"}
	req.Where = "sactive.eq.true"
	// Rudimentary geographical limit
	// req.Where += ",scoordinate.bbi.(10.368347,46.185535,12.551880,47.088826,4326)"
	req.Where += fmt.Sprintf(",sorigin.in.(%s)", p.origins())
	var res ninja.NinjaResponse[[]OdhParkingEcharging]
	err := ninja.StationType(req, &res)
	return res.Data, err
}

func (ParkingEcharging) mapNetex(os []OdhParkingEcharging) ([]netex.Parking, []netex.Operator) {
	ops := make(map[string]netex.Operator)

	var ps []netex.Parking
	for _, o := range os {
		var p netex.Parking

		p.Id = netex.CreateID("Parking", o.Scode)
		p.Version = "1"
		p.ShortName = o.Sname
		p.Centroid.Location.Longitude = o.Scoordinate.X
		p.Centroid.Location.Latitude = o.Scoordinate.Y
		p.GmlPolygon = nil
		op := netex.GetOperator(&config.Cfg, o.Sorigin)
		ops[op.Id] = op
		p.OperatorRef = netex.MkRef("Operator", op.Id)

		p.Entrances = nil
		p.ParkingType = "roadside"
		p.ParkingVehicleTypes = "car"
		p.ParkingLayout = "undefined"
		p.ProhibitedForHazardousMaterials.Ignore()
		p.RechargingAvailable.Set(true)
		p.Secure.Ignore()
		p.ParkingReservation = "reservationAllowed"
		p.ParkingProperties = nil

		p.Name = o.Sname
		p.PrincipalCapacity = o.Smetadata.Capacity
		p.TotalCapacity = o.Smetadata.Capacity

		ps = append(ps, p)
	}
	return ps, maps.Values(ops)
}

func (p ParkingEcharging) StParking() (netex.StParkingData, error) {
	odh, err := p.odhStatic()
	if err != nil {
		return netex.StParkingData{}, err
	}
	parkings, operators := p.mapNetex(odh)
	return netex.StParkingData{Parkings: parkings, Operators: operators}, nil
}

func (p ParkingEcharging) odhLatest(q siri.Query) ([]OdhParkingLatest, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = q.MaxSize()
	req.Repr = ninja.FlatNode
	req.StationTypes = []string{"EChargingStation"}
	req.DataTypes = []string{"number-available"}
	req.Select = "mperiod,mvalue,mvalidtime,scode,stype,smetadata.capacity"
	req.Where = "sactive.eq.true"
	req.Where += fmt.Sprintf(",sorigin.in.(%s)", p.origins())
	var res ninja.NinjaResponse[[]OdhParkingLatest]
	if err := ninja.Latest(req, &res); err != nil {
		slog.Error("Error retrieving parking state", "err", err)
		return res.Data, err
	}
	return res.Data, nil
}

func (p ParkingEcharging) mapSiri(latest []OdhParkingLatest) []siri.FacilityCondition {
	ret := []siri.FacilityCondition{}

	for _, o := range latest {
		fc := siri.FacilityCondition{}
		fc.FacilityRef = netex.CreateID("Parking", o.Scode)
		fc.MonitoredCounting = &siri.MonitoredCounting{}
		fc.MonitoredCounting.CountingType = "presentCount"

		fc.FacilityStatus.Status = siri.MapFacilityStatus(o.MValue, 1)
		fc.MonitoredCounting.CountedFeatureUnit = "devices"
		fc.MonitoredCounting.Count = o.Capacity - o.MValue

		ret = append(ret, fc)
	}

	return ret
}

func (p ParkingEcharging) SiriFM(query siri.Query) (siri.FMData, error) {
	ret := siri.FMData{}
	idFilter := maybeIdMatch(query.FacilityRef(), netex.CreateID("Parking"))
	if len(query.FacilityRef()) > 0 && len(idFilter) == 0 {
		return ret, nil
	}

	l, err := p.odhLatest(query)
	if err != nil {
		return ret, err
	}
	ret.Conditions = filterFacilityConditions(p.mapSiri(l), idFilter)
	return ret, nil
}
