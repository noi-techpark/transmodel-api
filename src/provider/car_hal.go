// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package provider

import (
	"opendatahub/sta-nap-export/config"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"

	"golang.org/x/exp/maps"
)

type odhHALCar []struct {
	ninja.OdhStation[HALCarMeta]
	Pmetadata struct {
		Company struct {
			Uid       string
			ShortName string
			FullName  string
		}
	}
}

type CarHAL struct {
	cars     odhHALCar
	origin   string
	provider string
}

type HALCarMeta struct {
	Brand        string
	Model        string
	LicensePlate string
	Features     *struct {
		Doors           uint8
		Seats           uint8
		Chains          bool
		Satnav          bool
		Skirack         bool
		Roofrack        bool
		Childseat       string
		Cyclerack       bool
		Wintertyres     bool
		Transmission    string
		Cruisecontrol   bool
		Trailerhitch    bool
		Usbpowersockets bool
	}
}

const ORIGIN_CAR_SHARING_HAL_API = "HAL-API"

func (b *CarHAL) init() error {
	b.origin = ORIGIN_CAR_SHARING_HAL_API
	if err := b.fetch(); err != nil {
		return err
	}
	b.provider = b.cars[0].Pmetadata.Company.ShortName // As soon as there is more than one, we have to change some more stuff anyways
	return nil
}

func (b *CarHAL) StSharing() (netex.StSharingData, error) {
	ret := netex.StSharingData{}
	if err := b.init(); err != nil {
		return ret, err
	}

	// Operators
	o := netex.Cfg.GetOperator(b.origin)
	ret.Operators = append(ret.Operators, o)

	// Modes of Operation
	m := netex.VehicleSharing{}
	m.Id = netex.CreateID("VehicleSharing", b.provider)
	m.Version = "1"
	sub := netex.Submode{}
	sub.Id = netex.CreateID("Submode", b.provider)
	sub.Version = "1"
	sub.TransportMode = "car"
	sub.SelfDriveSubmode = "hireCar"
	m.Submodes = append(m.Submodes, sub)
	ret.Modes = append(ret.Modes, m)

	models := make(map[string]netex.CarModelProfile)

	for _, c := range b.cars {
		modelname := c.Smeta.Brand
		p, found := models[modelname]
		if !found {
			// Car model profile
			p = netex.CarModelProfile{}
			p.Id = netex.CreateID("CarModelProfile", b.provider, modelname)
			p.Version = "1"
			p.ChildSeat = c.Smeta.Features.Childseat
			p.Seats = c.Smeta.Features.Seats
			p.Doors = c.Smeta.Features.Doors
			p.Transmission = c.Smeta.Features.Transmission
			p.CruiseControl = c.Smeta.Features.Cruisecontrol
			p.SatNav = c.Smeta.Features.Satnav
			p.AirConditioning = true
			p.Convertible = false
			p.UsbPowerSockets = c.Smeta.Features.Usbpowersockets
			p.WinterTyres = c.Smeta.Features.Wintertyres
			p.Chains = c.Smeta.Features.Chains
			p.TrailerHitch = c.Smeta.Features.Trailerhitch
			p.RoofRack = c.Smeta.Features.Roofrack
			p.CycleRack = c.Smeta.Features.Cyclerack
			p.SkiRack = c.Smeta.Features.Skirack
			models[modelname] = p
		}

		// Vehicles
		v := netex.Vehicle{}
		v.Id = netex.CreateID("Vehicle", b.provider, c.Scode)
		v.Version = "1"
		v.ValidBetween.AYear()
		v.Name = c.Sname
		v.ShortName = c.Sname
		v.PrivateCode = c.Scode
		v.RegistrationNumber = c.Smeta.LicensePlate
		v.OperatorRef = netex.MkRef("Operator", o.Id)
		v.VehicleTypeRef = netex.MkRef("CarModelProfile", p.Id)
		ret.Vehicles = append(ret.Vehicles, v)
	}
	ret.CarModels = maps.Values(models)

	// Fleets = all Vehicles + operator
	f := netex.Fleet{}
	f.Id = netex.CreateID("Fleet", b.provider)
	f.Version = "1"
	f.ValidBetween.AYear()
	for _, v := range ret.Vehicles {
		f.Members = append(f.Members, netex.MkRef("Vehicle", v.Id))
	}
	f.OperatorRef = netex.MkRef("Operator", o.Id)
	ret.Fleets = append(ret.Fleets, f)

	// Mobility services = Fleet + mode
	s := netex.VehicleSharingService{}
	s.Id = netex.CreateID("VehicleSharingService", b.provider)
	s.Version = "1"
	s.VehicleSharingRef = netex.MkRef("VehicleSharing", m.Id)
	s.FloatingVehicles = false
	for _, fl := range ret.Fleets {
		s.Fleets = append(s.Fleets, netex.MkRef("Fleet", fl.Id))
	}
	ret.Services = append(ret.Services, s)

	// Constraint zone
	c := netex.MobilityServiceConstraintZone{}
	c.Id = netex.CreateID("MobilityServiceConstraintZone", b.provider)
	c.Version = "1"
	c.GmlPolygon.Id = b.provider
	c.GmlPolygon.SetPoly(config.GML_PROVINCE_BZ)
	c.VehicleSharingRef = netex.MkRef("VehicleSharingService", s.Id)
	ret.Constraints = append(ret.Constraints, c)

	return ret, nil
}

func (b *CarHAL) fetch() error {
	cs, err := odhMob[odhHALCar]("CarsharingCar", b.origin)
	b.cars = cs
	return err
}
