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

	"golang.org/x/exp/maps"
)

type bzCycleMeta struct {
	Model    string
	Electric bool
	Lamp     bool
	Lock     bool
	Basket   bool
}

type bzSharingMeta struct {
	Address   string
	TotalBays int `json:"total-bays"`
}

type bikeBzCycles []ninja.OdhStation[bzCycleMeta]
type bikeBzSharing []ninja.OdhStation[bzSharingMeta]

type BikeBz struct {
	origin  string
	cycles  func(string) (bikeBzCycles, error)
	sharing func(string) (bikeBzSharing, error)
}

const ORIGIN_BIKE_SHARING_BOLZANO = "BIKE_SHARING_BOLZANO"

func NewBikeBz() *BikeBz {
	b := BikeBz{}
	b.origin = ORIGIN_BIKE_SHARING_BOLZANO
	b.cycles = func(origin string) (bikeBzCycles, error) {
		return FetchOdhStations[bikeBzCycles]("Bicycle", origin)
	}
	b.sharing = func(origin string) (bikeBzSharing, error) {
		return FetchOdhStations[bikeBzSharing]("BikesharingStation", origin)
	}
	return &b
}

func (b *BikeBz) StSharing() (netex.StSharingData, error) {
	ret := netex.StSharingData{}

	// Operators
	o := netex.GetOperator(&config.Cfg, ORIGIN_BIKE_SHARING_BOLZANO)
	ret.Operators = append(ret.Operators, o)

	// Modes of Operation
	m := netex.VehicleSharing{}
	m.Id = netex.CreateID("VehicleSharing", b.origin)
	m.Version = "1"
	sub := netex.Submode{}
	sub.Id = netex.CreateID("Submode", b.origin)
	sub.Version = "1"
	sub.TransportMode = "bicycle"
	sub.SelfDriveSubmode = "hireCycle"
	m.Submodes = append(m.Submodes, sub)
	ret.Modes = append(ret.Modes, m)

	models := make(map[string]netex.CycleModelProfile)

	cycles, err := b.cycles(b.origin)
	if err != nil {
		return ret, err
	}

	for _, c := range cycles {
		p, found := models[c.Smeta.Model]
		if !found {
			// Cycle model profile
			p = netex.CycleModelProfile{}
			p.Id = netex.CreateID("CycleModelProfile", b.origin, c.Smeta.Model)
			p.Version = "1"
			p.ChildSeat = "none"
			// Assume model and features map correctly, just take the first one we encounter
			p.Battery = c.Smeta.Electric
			p.Lamps = c.Smeta.Lamp
			p.Pump = false
			p.Basket = c.Smeta.Basket
			p.Lock = false
			models[c.Smeta.Model] = p
		}

		// Vehicles
		v := netex.Vehicle{}
		v.Id = netex.CreateID("Vehicle", b.origin, c.Scode)
		v.Version = "1"
		v.ValidBetween.AYear()
		v.Name = c.Sname
		v.ShortName = c.Sname
		v.PrivateCode = c.Scode
		v.OperatorRef = netex.MkRef("Operator", o.Id)
		v.VehicleTypeRef = netex.MkRef("CycleModelProfile", p.Id)
		ret.Vehicles = append(ret.Vehicles, v)
	}
	ret.CycleModels = maps.Values(models)

	// Fleets = all Vehicles + operator
	f := netex.Fleet{}
	f.Id = netex.CreateID("Fleet", b.origin)
	f.Version = "1"
	f.ValidBetween.AYear()
	for _, v := range ret.Vehicles {
		f.Members = append(f.Members, netex.MkRef("Vehicle", v.Id))
	}
	f.OperatorRef = netex.MkRef("Operator", o.Id)
	ret.Fleets = append(ret.Fleets, f)

	// Mobility services = Fleet + mode
	s := netex.VehicleSharingService{}
	s.Id = netex.CreateID("VehicleSharingService", b.origin)
	s.Version = "1"
	s.VehicleSharingRef = netex.MkRef("VehicleSharing", m.Id)
	s.FloatingVehicles = false
	for _, fl := range ret.Fleets {
		s.Fleets = append(s.Fleets, netex.MkRef("Fleet", fl.Id))
	}
	ret.Services = append(ret.Services, s)

	// Constraint zone
	c := netex.MobilityServiceConstraintZone{}
	c.Id = netex.CreateID("MobilityServiceConstraintZone", b.origin)
	c.Version = "1"
	c.GmlPolygon.Id = b.origin
	c.GmlPolygon.SetPoly(config.GML_MUNICIPALITY_BZ)
	c.VehicleSharingRef = netex.MkRef("VehicleSharingService", s.Id)
	ret.Constraints = append(ret.Constraints, c)

	// Sharing ss as Parking (for SIRI reference)
	ss, err := b.sharing(b.origin)
	if err != nil {
		return ret, err
	}
	for _, s := range ss {
		p := netex.Parking{}
		p.Id = netex.CreateID("Parking", b.origin, s.Sname)
		p.Version = "1"
		p.ShortName = s.Sname
		p.Centroid.Location.Longitude = s.Scoord.X
		p.Centroid.Location.Latitude = s.Scoord.Y
		p.OperatorRef = netex.MkRef("Operator", o.Id)
		p.GmlPolygon = nil
		p.Entrances = nil
		p.ParkingType = "cycleRental"
		p.ParkingVehicleTypes = "cycle"
		p.ParkingLayout = "cycleHire"
		p.ProhibitedForHazardousMaterials.Ignore()
		p.RechargingAvailable.Set(true)
		p.Secure.Set(false)
		p.ParkingReservation = "registrationRequired"
		p.ParkingProperties = nil

		p.Name = s.Sname
		p.PrincipalCapacity = int32(s.Smeta.TotalBays)
		p.TotalCapacity = int32(s.Smeta.TotalBays)
		ret.Parkings = append(ret.Parkings, p)
	}

	return ret, nil
}

type OdhBzSharingLatest struct {
	ninja.OdhLatest
	Sname string
}

func (p BikeBz) odhLatest(q siri.Query) ([]OdhBzSharingLatest, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = q.MaxSize()
	req.Repr = ninja.FlatNode
	req.StationTypes = []string{"BikesharingStation"}
	req.DataTypes = []string{"free-bays,number-available"}
	req.Select = "mperiod,mvalue,mvalidtime,scode,sname,tname"
	req.Where = "sactive.eq.true"
	req.Where += fmt.Sprintf(",sorigin.eq.%s", p.origin)
	req.Where += filterIDs(q.FacilityRef(), netex.CreateID("Parking", p.origin), "sname")

	var res ninja.NinjaResponse[[]OdhBzSharingLatest]
	err := ninja.Latest(req, &res)
	if err != nil {
		slog.Error("Error retrieving parking state", "err", err)
	}
	return res.Data, err
}

func (p BikeBz) mapSiri(latest []OdhBzSharingLatest) []siri.FacilityCondition {
	ret := []siri.FacilityCondition{}

	type station struct {
		name      string
		free      int
		available int
	}

	stations := map[string]station{}

	// group data types by station
	for _, l := range latest {
		s, found := stations[l.Scode]
		if !found {
			s = station{name: l.Sname}
		}
		switch l.Tname {
		case "free-bays":
			s.free = l.MValue
		case "number-available":
			s.available = l.MValue
		}
		stations[l.Scode] = s
	}

	for _, o := range stations {
		fc := siri.FacilityCondition{}
		fc.FacilityRef = netex.CreateID("Parking", p.origin, o.name)
		fc.FacilityStatus.Status = siri.MapFacilityStatus(o.free, 1)
		fc.MonitoredCounting = &siri.MonitoredCounting{}
		fc.MonitoredCounting.CountingType = "availabilityCount"
		fc.MonitoredCounting.CountedFeatureUnit = "otherSpaces"
		fc.MonitoredCounting.Count = o.available

		ret = append(ret, fc)
	}

	return ret
}

func (p BikeBz) SiriFM(query siri.Query) (siri.FMData, error) {
	ret := siri.FMData{}
	l, err := p.odhLatest(query)
	if err != nil {
		return ret, err
	}
	ret.Conditions = p.mapSiri(l)
	return ret, nil
}
