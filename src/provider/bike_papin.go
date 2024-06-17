// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package provider

import (
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"
)

type odhPapinBike []ninja.OdhStation[any]

type BikePapin struct {
	cycles odhPapinBike
	origin string
}

const ORIGIN_BIKE_SHARING_PAPIN = "BIKE_SHARING_PAPIN"

func (b *BikePapin) init() error {
	b.origin = ORIGIN_BIKE_SHARING_PAPIN

	return b.fetch()
}

func (b *BikePapin) StSharing() (netex.StSharingData, error) {
	ret := netex.StSharingData{}
	if err := b.init(); err != nil {
		return ret, err
	}

	// Operators
	o := netex.Cfg.GetOperator(b.origin)
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

	// Fleets = all Vehicles + operator
	f := netex.Fleet{}
	f.Id = netex.CreateID("Fleet", b.origin)
	f.Version = "1"
	f.ValidBetween.AYear()
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

	return ret, nil
}
func (b *BikePapin) fetch() error {
	bk, err := odhMob[odhPapinBike]("Bicycle", b.origin)
	b.cycles = bk
	return err
}
