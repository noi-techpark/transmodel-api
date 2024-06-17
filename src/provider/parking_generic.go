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

type OdhParkingGeneric struct {
	Scode       string
	Sname       string
	Sorigin     string
	Stype       string
	Scoordinate struct {
		X    float32
		Y    float32
		Srid uint32
	}
	Smetadata struct {
		StandardName string
		Capacity     int32
		TotalPlaces  int32 // Bikeparking specific capacity
		Municipality string
		Netex        struct {
			Type             string
			VehicleTypes     string
			Layout           string
			HazardProhibited bool `json:"hazard_prohibited"`
			Charging         bool
			Surveillance     bool
			Reservation      string `json:"reservation"`
		} `json:"parking"`
	}
}

type ParkingGeneric struct {
}

func parkingOrigins() string {
	origins := config.Cfg.ParkingOrigins()
	quoted := []string{}
	for _, o := range origins {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", o))
	}
	return strings.Join(quoted, ",")
}

func (p ParkingGeneric) odhStatic() ([]OdhParkingGeneric, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"ParkingStation", "BikeParking"}
	req.Where = "sactive.eq.true"
	req.Where += fmt.Sprintf(",sorigin.in.(%s)", parkingOrigins())
	// TODO: limit bounding box / polygon
	var res ninja.NinjaResponse[[]OdhParkingGeneric]
	err := ninja.StationType(req, &res)
	return res.Data, err
}

func defEmpty(s string, d string) string {
	if s == "" {
		return d
	} else {
		return s
	}
}

func (ParkingGeneric) originButWithHacks(p OdhParkingGeneric) string {
	// There are two different Operators that both have origin FBK (Trento and Rovereto)
	if p.Sorigin == "FBK" {
		return fmt.Sprintf("%s-%s", p.Sorigin, p.Smetadata.Municipality)
	}
	return p.Sorigin
}

func (pg ParkingGeneric) mapNetex(os []OdhParkingGeneric) ([]netex.Parking, []netex.Operator) {
	ops := make(map[string]netex.Operator)

	var ps []netex.Parking
	for _, o := range os {
		var p netex.Parking

		p.Id = netex.CreateID("Parking", o.Scode)
		p.Version = "1"
		p.ShortName = o.Sname
		// p.Centroid.Location.Precision = 1  not sure what this actually does, according to specification not needed?
		p.Centroid.Location.Longitude = o.Scoordinate.X
		p.Centroid.Location.Latitude = o.Scoordinate.Y
		p.GmlPolygon = nil
		op := netex.GetOperator(&config.Cfg, pg.originButWithHacks(o))
		ops[op.Id] = op
		p.OperatorRef = netex.MkRef("Operator", op.Id)

		p.Entrances = nil
		p.ParkingType = defEmpty(o.Smetadata.Netex.Type, "undefined")
		p.ParkingVehicleTypes = o.Smetadata.Netex.VehicleTypes
		p.ParkingLayout = defEmpty(o.Smetadata.Netex.Layout, "undefined")
		p.ProhibitedForHazardousMaterials.Set(o.Smetadata.Netex.HazardProhibited)
		p.RechargingAvailable = o.Smetadata.Netex.Charging
		p.Secure.Set(o.Smetadata.Netex.Surveillance)
		p.ParkingReservation = defEmpty(o.Smetadata.Netex.Reservation, "noReservations")
		p.ParkingProperties = nil

		if o.Stype == "BikeParking" {
			p.Name = o.Sname
			p.PrincipalCapacity = o.Smetadata.TotalPlaces
			p.TotalCapacity = o.Smetadata.TotalPlaces
		} else {
			p.Name = o.Smetadata.StandardName
			p.PrincipalCapacity = o.Smetadata.Capacity
			p.TotalCapacity = o.Smetadata.Capacity
		}

		ps = append(ps, p)
	}
	return ps, maps.Values(ops)
}

func (p ParkingGeneric) StParking() (netex.StParkingData, error) {
	odh, err := p.odhStatic()
	if err != nil {
		return netex.StParkingData{}, err
	}
	parkings, operators := p.mapNetex(odh)
	return netex.StParkingData{Parkings: parkings, Operators: operators}, nil
}

type OdhParkingLatest struct {
	ninja.OdhLatest
	Capacity    int `json:"smetadata.capacity"`
	TotalPlaces int `json:"smetadata.totalPlaces"`
}

func (p ParkingGeneric) odhLatest() ([]OdhParkingLatest, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.Repr = ninja.FlatNode
	req.StationTypes = []string{"ParkingStation", "BikeParking"}
	req.DataTypes = []string{"free", "number-available"}
	req.Select = "mperiod,mvalue,mvalidtime,scode,stype,smetadata.capacity,smetadata.totalPlaces"
	req.Where = "sactive.eq.true"
	req.Where += fmt.Sprintf(",sorigin.in.(%s)", netex.ParkingOrigins())
	var res ninja.NinjaResponse[[]OdhParkingLatest]
	err := ninja.Latest(req, &res)
	if err != nil {
		slog.Error("Error retrieving parking state", "err", err)
	}
	return res.Data, err
}

func (p ParkingGeneric) mapSiri(latest []OdhParkingLatest) []siri.FacilityCondition {
	ret := []siri.FacilityCondition{}

	for _, o := range latest {
		fc := siri.FacilityCondition{}
		fc.FacilityRef = netex.CreateID("Parking", o.Scode)
		fc.MonitoredCounting.CountingType = "presentCount"

		switch o.Stype {
		case "BikeParking":
			fc.FacilityStatus.Status = siri.MapFacilityStatus(o.MValue, 10)
			fc.MonitoredCounting.CountedFeatureUnit = "otherSpaces"
			fc.MonitoredCounting.Count = o.TotalPlaces - o.MValue
		case "ParkingStation":
			fc.FacilityStatus.Status = siri.MapFacilityStatus(o.MValue, 10)
			fc.MonitoredCounting.CountedFeatureUnit = "bays"
			fc.MonitoredCounting.Count = o.Capacity - o.MValue
		}

		ret = append(ret, fc)
	}

	return ret
}

func (p ParkingGeneric) RtSharing() (siri.FMData, error) {
	ret := siri.FMData{}
	l, err := p.odhLatest()
	if err != nil {
		return ret, err
	}
	ret.Conditions = p.mapSiri(l)
	return ret, nil
}
