// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"encoding/json"
	"log/slog"
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
	Fleets      []Fleet
	Vehicles    []Vehicle
	CarModels   []CarModelProfile
	CycleModels []CycleModelProfile
	Operators   []Operator
	Modes       []VehicleSharing
	Services    []VehicleSharingService
	Constraints []MobilityServiceConstraintZone
	Parkings    []Parking
}

type StSharing interface {
	StSharing() (StSharingData, error)
}

func GetSharing(bikeProviders []StSharing, carProviders []StSharing) ([]CompositeFrame, error) {
	ret := []CompositeFrame{}

	c, err := compBikeSharing(bikeProviders)
	if err != nil {
		return ret, err
	}
	ret = append(ret, c)

	c, err = compCarSharing(carProviders)
	if err != nil {
		return ret, err
	}
	ret = append(ret, c)

	return ret, nil
}
func compBikeSharing(ps []StSharing) (CompositeFrame, error) {
	mob := MobilityServiceFrame{}
	mob.Id = CreateFrameId("MobilityServiceFrame_EU_PI_MOBILITY", "BikeSharing")
	mob.Version = "1"
	mob.FrameDefaults.DefaultCurrency = "EUR"

	res := ResourceFrame{}
	res.Id = CreateFrameId("ResourceFrame_EU_PI_COMMON", "BikeSharing")
	res.Version = "1"
	res.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_COMMON")

	site := SiteFrame{}
	site.Id = CreateFrameId("SiteFrame_EU_PI_STOP", "BikeSharing")
	site.Version = "1"
	site.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_STOP")

	for _, p := range ps {
		d, err := p.StSharing()
		if err != nil {
			return CompositeFrame{}, err
		}

		mob.Fleets = append(mob.Fleets, d.Fleets...)
		mob.ModesOfOperation = append(mob.ModesOfOperation, d.Modes...)
		mob.MobilityServices = append(mob.MobilityServices, d.Services...)
		mob.MobilityServiceConstraintZones = append(mob.MobilityServiceConstraintZones, d.Constraints...)

		res.Vehicles = AppendSafe(res.Vehicles, d.Vehicles...)
		res.CarModels = AppendSafe(res.CarModels, d.CarModels...)
		res.CycleModels = AppendSafe(res.CycleModels, d.CycleModels...)
		res.Operators = AppendSafe(res.Operators, d.Operators...)

		site.Parkings.Parkings = append(site.Parkings.Parkings, d.Parkings...)
	}

	comp := CompositeFrame{}
	comp.Defaults()
	comp.Id = CreateFrameId("CompositeFrame_EU_PI_STOP_OFFER", "SHARING", "BikeSharing")
	comp.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_LINE_OFFER")
	comp.Frames.Frames = append(comp.Frames.Frames, mob, res, site)

	return comp, nil
}

func compCarSharing(ps []StSharing) (CompositeFrame, error) {
	mob := MobilityServiceFrame{}
	mob.Id = CreateFrameId("MobilityServiceFrame_EU_PI_MOBILITY", "CarSharing")
	mob.Version = "1"
	mob.FrameDefaults.DefaultCurrency = "EUR"

	res := ResourceFrame{}
	res.Id = CreateFrameId("ResourceFrame_EU_PI_COMMON", "CarSharing")
	res.Version = "1"
	res.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_COMMON")

	site := SiteFrame{}
	site.Id = CreateFrameId("SiteFrame_EU_PI_STOP", "CarSharing")
	site.Version = "1"
	site.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_STOP")

	for _, p := range ps {
		d, err := p.StSharing()
		if err != nil {
			return CompositeFrame{}, err
		}

		mob.Fleets = append(mob.Fleets, d.Fleets...)
		mob.ModesOfOperation = append(mob.ModesOfOperation, d.Modes...)
		mob.MobilityServices = append(mob.MobilityServices, d.Services...)
		mob.MobilityServiceConstraintZones = append(mob.MobilityServiceConstraintZones, d.Constraints...)

		res.Vehicles = AppendSafe(res.Vehicles, d.Vehicles...)
		res.CarModels = AppendSafe(res.CarModels, d.CarModels...)
		res.CycleModels = AppendSafe(res.CycleModels, d.CycleModels...)
		res.Operators = AppendSafe(res.Operators, d.Operators...)

		site.Parkings.Parkings = append(site.Parkings.Parkings, d.Parkings...)
	}

	comp := CompositeFrame{}
	comp.Defaults()
	comp.Id = CreateFrameId("CompositeFrame_EU_PI_STOP_OFFER", "SHARING", "CarSharing")
	comp.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_LINE_OFFER")
	comp.Frames.Frames = append(comp.Frames.Frames, mob, res, site)

	return comp, nil
}
