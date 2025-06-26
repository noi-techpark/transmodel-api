// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package comp

import (
	"encoding/json"
	"log/slog"

	"github.com/noi-techpark/go-netex"
)

type Company struct {
	UID       string `json:"uid"`
	ShortName string `json:"shortName"`
	FullName  string `json:"fullName"`
}

func (pc *Company) UnmarshalJSON(p []byte) error {
	if string(p) == "" {
		// empty string, do nothing
		return nil
	}
	slog.Debug("unmarshalling ", "str", string(p))
	// Prevent recursion to this method by declaring a new
	// type with same underlying type as Company and
	// no methods.
	type x Company
	return json.Unmarshal(p, (*x)(pc))
}

type StSharingData struct {
	Fleets      []netex.Fleet
	Vehicles    []netex.Vehicle
	CarModels   []netex.CarModelProfile
	CycleModels []netex.CycleModelProfile
	Operators   []netex.Operator
	Modes       []netex.VehicleSharing
	Services    []netex.VehicleSharingService
	Constraints []netex.MobilityServiceConstraintZone
	Parkings    []netex.Parking
}

type StSharing interface {
	StSharing() (StSharingData, error)
}

func GetSharing(bikeProviders []StSharing, carProviders []StSharing) ([]netex.CompositeFrame, error) {
	ret := []netex.CompositeFrame{}

	c, err := compSharing("BikeSharing", bikeProviders)
	if err != nil {
		return ret, err
	}
	ret = append(ret, c)

	c, err = compSharing("CarSharing", carProviders)
	if err != nil {
		return ret, err
	}
	ret = append(ret, c)

	return ret, nil
}
func compSharing(serviceName string, ps []StSharing) (netex.CompositeFrame, error) {
	mob := netex.MobilityServiceFrame{}
	mob.Id = CreateFrameId(netex.TypeMobilityServiceFrameMobility, serviceName)
	mob.Version = "1"
	mob.FrameDefaults.DefaultCurrency = "EUR"
	mob.Fleets = new([]netex.Fleet)
	mob.ModesOfOperation = new([]netex.VehicleSharing)
	mob.MobilityServices = new([]netex.VehicleSharingService)
	mob.MobilityServiceConstraintZones = new([]netex.MobilityServiceConstraintZone)

	res := netex.ResourceFrame{}
	res.Id = CreateFrameId(netex.TypeResourceFrameCommon, serviceName)
	res.Version = "1"
	res.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeCommon)

	site := netex.SiteFrame{}
	site.Id = CreateFrameId(netex.TypeSiteFrameStop, serviceName)
	site.Version = "1"
	site.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeStop)

	for _, p := range ps {
		d, err := p.StSharing()
		if err != nil {
			return netex.CompositeFrame{}, err
		}

		mob.Fleets = netex.AppendMaybe(mob.Fleets, d.Fleets...)
		mob.ModesOfOperation = netex.AppendMaybe(mob.ModesOfOperation, d.Modes...)
		mob.MobilityServices = netex.AppendMaybe(mob.MobilityServices, d.Services...)
		mob.MobilityServiceConstraintZones = netex.AppendMaybe(mob.MobilityServiceConstraintZones, d.Constraints...)

		res.Vehicles = netex.AppendMaybe(res.Vehicles, d.Vehicles...)
		res.CarModels = netex.AppendMaybe(res.CarModels, d.CarModels...)
		res.CycleModels = netex.AppendMaybe(res.CycleModels, d.CycleModels...)
		res.Operators = netex.AppendMaybe(res.Operators, d.Operators...)

		site.Parkings = netex.AppendMaybe(site.Parkings, d.Parkings...)
	}

	comp := DefaultCompositFrame()
	comp.Id = CreateFrameId(netex.TypeCompositeFrameStopOffer, serviceName)
	comp.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeStopOffer)
	comp.Frames.Frames = append(comp.Frames.Frames, mob, res, site)

	return comp, nil
}
