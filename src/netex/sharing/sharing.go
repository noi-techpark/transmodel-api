// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package sharing

import (
	"encoding/json"
	"log/slog"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"
)

type Company struct {
	UID       string `json:"uid"`
	ShortName string `json:"shortName"`
	FullName  string `json:"fullName"`
}

type OdhMobility[T any] struct {
	Scode   string `json:"scode"`
	Sname   string `json:"sname"`
	Sorigin string `json:"sorigin"`
	Scoord  struct {
		X    float32 `json:"x"`
		Y    float32 `json:"y"`
		Srid uint32  `json:"srid"`
	} `json:"scoordinate"`
	Smeta T `json:"smetadata"`
}

type metaAny map[string]any

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

func odhMob[T any](tp string, origin string) (T, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{tp}
	req.Where = `sorigin.eq.` + origin + `,sactive.eq.true`
	var res ninja.NinjaResponse[T]
	err := ninja.StationType(req, &res)
	return res.Data, err
}

type SharingData struct {
	Fleets   []netex.Fleet
	Vehicles []netex.Vehicle
	// VehicleTypes  []netex.VehicleType
	// VehicleModels []netex.VehicleModel
	CarModels   []netex.CarModelProfile
	CycleModels []netex.CycleModelProfile
	Operators   []netex.Operator
	Modes       []netex.VehicleSharing
	Services    []netex.VehicleSharingService
	Constraints []netex.MobilityServiceConstraintZone
}

type SharingProvider interface {
	get() (SharingData, error)
}

func GetSharing() (netex.Root, error) {
	var ret netex.Root

	c, err := compBikeSharing([]SharingProvider{&BikeBz{}, &BikeMe{}, &BikePapin{}})
	if err != nil {
		return ret, err
	}
	ret.CompositeFrame = append(ret.CompositeFrame, c)

	c, err = compCarSharing([]SharingProvider{&CarHAL{}})
	if err != nil {
		return ret, err
	}
	ret.CompositeFrame = append(ret.CompositeFrame, c)

	return ret, nil
}
func compBikeSharing(ps []SharingProvider) (netex.CompositeFrame, error) {
	mob := netex.MobilityServiceFrame{}
	mob.Id = netex.CreateFrameId("MobilityServiceFrame_EU_PI_MOBILITY", "BikeSharing")
	mob.Version = "1"
	mob.FrameDefaults.DefaultCurrency = "EUR"

	res := netex.ResourceFrame{}
	res.Id = netex.CreateFrameId("ResourceFrame_EU_PI_COMMON", "BikeSharing")
	res.Version = "1"
	res.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_COMMON")

	for _, p := range ps {
		d, err := p.get()
		if err != nil {
			return netex.CompositeFrame{}, err
		}

		mob.Fleets = append(mob.Fleets, d.Fleets...)
		mob.ModesOfOperation = append(mob.ModesOfOperation, d.Modes...)
		mob.MobilityServices = append(mob.MobilityServices, d.Services...)
		mob.MobilityServiceConstraintZones = append(mob.MobilityServiceConstraintZones, d.Constraints...)

		res.Vehicles = netex.AppendSafe(res.Vehicles, d.Vehicles...)
		res.CarModels = netex.AppendSafe(res.CarModels, d.CarModels...)
		res.CycleModels = netex.AppendSafe(res.CycleModels, d.CycleModels...)
		res.Operators = netex.AppendSafe(res.Operators, d.Operators...)
	}

	comp := netex.CompositeFrame{}
	comp.Defaults()
	comp.Id = netex.CreateFrameId("CompositeFrame_EU_PI_STOP_OFFER", "SHARING", "BikeSharing")
	comp.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_LINE_OFFER")
	comp.Frames.Frames = append(comp.Frames.Frames, mob, res)

	return comp, nil
}

func compCarSharing(ps []SharingProvider) (netex.CompositeFrame, error) {
	mob := netex.MobilityServiceFrame{}
	mob.Id = netex.CreateFrameId("MobilityServiceFrame_EU_PI_MOBILITY", "CarSharing")
	mob.Version = "1"
	mob.FrameDefaults.DefaultCurrency = "EUR"

	res := netex.ResourceFrame{}
	res.Id = netex.CreateFrameId("ResourceFrame_EU_PI_COMMON", "CarSharing")
	res.Version = "1"
	res.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_COMMON")

	for _, p := range ps {
		d, err := p.get()
		if err != nil {
			return netex.CompositeFrame{}, err
		}

		mob.Fleets = append(mob.Fleets, d.Fleets...)
		mob.ModesOfOperation = append(mob.ModesOfOperation, d.Modes...)
		mob.MobilityServices = append(mob.MobilityServices, d.Services...)
		mob.MobilityServiceConstraintZones = append(mob.MobilityServiceConstraintZones, d.Constraints...)

		res.Vehicles = netex.AppendSafe(res.Vehicles, d.Vehicles...)
		res.CarModels = netex.AppendSafe(res.CarModels, d.CarModels...)
		res.CycleModels = netex.AppendSafe(res.CycleModels, d.CycleModels...)
		res.Operators = netex.AppendSafe(res.Operators, d.Operators...)
	}

	comp := netex.CompositeFrame{}
	comp.Defaults()
	comp.Id = netex.CreateFrameId("CompositeFrame_EU_PI_STOP_OFFER", "SHARING", "CarSharing")
	comp.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_LINE_OFFER")
	comp.Frames.Frames = append(comp.Frames.Frames, mob, res)

	return comp, nil
}
