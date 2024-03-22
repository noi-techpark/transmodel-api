// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package sharing

import "opendatahub/sta-nap-export/netex"

type odhBzShare []OdhMobility[metaAny]
type odhBzBike []OdhMobility[metaAny]

type Bz struct {
	sharing odhBzShare
	cycles  odhBzBike
	origin  string
}

const ORIGIN_BIKE_SHARING_BOLZANO = "BIKE_SHARING_BOLZANO"

func (b *Bz) init() error {
	b.origin = ORIGIN_BIKE_SHARING_BOLZANO

	s, err := bzSharing()
	if err != nil {
		return err
	}
	b.sharing = s

	bk, err := bzBike()
	if err != nil {
		return err
	}
	b.cycles = bk
	return nil
}

func (b *Bz) get() (SharingData, error) {
	ret := SharingData{}
	if err := b.init(); err != nil {
		return ret, err
	}

	// Operators
	o := Operator{}
	o.Id = netex.CreateID("Operator", b.origin)
	o.Version = "1"
	o.PrivateCode = b.origin
	o.Name = b.origin
	o.ShortName = b.origin
	o.LegalName = b.origin
	o.OrganizationType = "operator"
	ret.operators = append(ret.operators, o)

	// Modes of Operation
	m := VehicleSharing{}
	m.Id = netex.CreateID("VehicleSharing", b.origin)
	m.Version = "1"
	sub := Submode{}
	sub.Id = netex.CreateID("Submode", b.origin)
	sub.Version = "1"
	sub.TransportMode = "bicycle"
	sub.SelfDriveMode = "hireCycle"
	m.Submodes = append(m.Submodes, sub)
	ret.modes = append(ret.modes, m)

	// Cycle model profile
	p := CycleModelProfile{}
	p.Id = netex.CreateID("CycleModelProfile", b.origin, "default")
	p.Version = "1"
	p.ChildSeat = "none"
	p.Battery = true // todo map
	p.Lamps = true   // todo map
	p.Pump = false
	p.Basket = true // todo map
	p.Lock = false  // true for merano
	ret.cycleModels = append(ret.cycleModels, p)

	// Vehicles
	for _, c := range b.cycles {
		v := Vehicle{}
		v.Id = netex.CreateID("Vehicle", b.origin, c.Scode)
		v.Version = "1"
		v.Name = c.Sname
		v.ShortName = c.Sname
		v.PrivateCode = c.Scode
		v.OperatorRef = MkRef("Operator", o.Id)
		v.VehicleTypeRef = MkRef("VehicleType", p.Id)
		ret.vehicles = append(ret.vehicles, v)
	}

	// Fleets = all Vehicles + operator
	f := Fleet{}
	f.Id = netex.CreateID("Fleet", b.origin)
	f.Version = "1"
	for _, v := range ret.vehicles {
		f.Members = append(f.Members, MkRef("Vehicle", v.Id))
	}
	f.OperatorRef = MkRef("Operator", o.Id)
	ret.fleets = append(ret.fleets, f)

	// Mobility services = Fleet + mode

	s := VehicleSharingService{}
	s.Id = netex.CreateID("VehicleSharingService", b.origin)
	s.Version = "1"
	s.VehicleSharingRef = MkRef("VehicleSharing", m.Id)
	s.FloatingVehicles = false
	for _, fl := range ret.fleets {
		s.Fleets = append(s.Fleets, MkRef("Fleet", fl.Id))
	}
	ret.services = append(ret.services, s)

	// Constraint zone
	c := MobilityServiceConstraintZone{}
	c.Id = netex.CreateID("MobilityServiceConstraintZone", b.origin)
	c.Version = "1"
	c.GmlPolygon = ""
	c.VehicleSharingRef = MkRef("VehicleSharingService", s.Id)
	ret.constraints = append(ret.constraints, c)

	return ret, nil
}

func bzSharing() (odhBzShare, error) {
	return odhMob[odhBzShare]("BikesharingStation", ORIGIN_BIKE_SHARING_BOLZANO)
}
func bzBike() (odhBzBike, error) {
	return odhMob[odhBzBike]("Bicycle", ORIGIN_BIKE_SHARING_BOLZANO)
}
