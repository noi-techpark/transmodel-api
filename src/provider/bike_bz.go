// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package provider

import (
	"opendatahub/sta-nap-export/config"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"

	"golang.org/x/exp/maps"
)

type odhBzBike []ninja.OdhStation[BzCycleMeta]

type BikeBz struct {
	cycles odhBzBike
	origin string
}

type BzCycleMeta struct {
	Model    string
	Electric bool
	Lamp     bool
	Lock     bool
	Basket   bool
}

const ORIGIN_BIKE_SHARING_BOLZANO = "BIKE_SHARING_BOLZANO"

func (b *BikeBz) init() error {
	b.origin = ORIGIN_BIKE_SHARING_BOLZANO

	bk, err := bzBike()
	if err != nil {
		return err
	}
	b.cycles = bk
	return nil
}

func (b *BikeBz) StSharing() (netex.StSharingData, error) {
	ret := netex.StSharingData{}
	if err := b.init(); err != nil {
		return ret, err
	}

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

	for _, c := range b.cycles {
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

	return ret, nil
}
func bzBike() (odhBzBike, error) {
	return odhMob[odhBzBike]("Bicycle", ORIGIN_BIKE_SHARING_BOLZANO)
}
