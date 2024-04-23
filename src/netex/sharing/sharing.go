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

// type SharingMeta struct {
// 	Company      Company `json:"company"`
// 	LicensePlate string  `json:"licensePlate"`
// 	Electric     bool    `json:"electric"`
// 	Brand        bool    `json:"brand"`
// }

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
	// type with same underlying type as PrimaryContact and
	// no methods.
	type x Company
	return json.Unmarshal(p, (*x)(pc))
}

func odhMob[T any](tp string, origin string) (T, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{tp}
	req.Where = `sorigin.eq.` + origin
	var res ninja.NinjaResponse[T]
	err := ninja.StationType(req, &res)
	return res.Data, err
}

type SharingData struct {
	Fleets        []Fleet
	Vehicles      []Vehicle
	VehicleModels []VehicleModel
	CarModels     []CarModelProfile
	CycleModels   []CycleModelProfile
	Operators     []Operator
	Modes         []VehicleSharing
	Services      []VehicleSharingService
	Constraints   []MobilityServiceConstraintZone
}

type SharingProvider interface {
	get() (SharingData, error)
}

func GetSharing() (*netex.CompositeFrame, error) {
	return frame([]SharingProvider{&Bz{}})
}

func frame(ps []SharingProvider) (*netex.CompositeFrame, error) {
	mob := MobilityServiceFrame{}
	mob.Id = netex.CreateFrameId("MobilityServiceFrame_EU_PI_MOBILITY", "BikeSharing", "ita")
	mob.Version = "1"
	mob.FrameDefaults.DefaultCurrency = "EUR"

	res := ResourceFrame{}
	res.Id = netex.CreateFrameId("ResourceFrame_EU_PI_MOBILITY", "ita")
	res.Version = "1"
	res.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_COMMON")

	for _, p := range ps {
		d, err := p.get()
		if err != nil {
			return nil, err
		}
		mob.Fleets = append(mob.Fleets, d.Fleets...)
		mob.ModesOfOperation = append(mob.ModesOfOperation, d.Modes...)
		mob.MobilityServices = append(mob.MobilityServices, d.Services...)
		mob.MobilityServiceConstraintZones = append(mob.MobilityServiceConstraintZones, d.Constraints...)

		res.Vehicles = append(res.Vehicles, d.Vehicles...)
		res.VehicleModels = append(res.VehicleModels, d.VehicleModels...)
		res.CarModels = append(res.CarModels, d.CarModels...)
		res.CycleModels = append(res.CycleModels, d.CycleModels...)
		res.Operators = append(res.Operators, d.Operators...)
	}

	comp := netex.CompositeFrame{}
	comp.Defaults()
	comp.Id = netex.CreateFrameId("CompositeFrame_EU_PI_STOP_OFFER", "SHARING", "ita")
	comp.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_LINE_OFFER")
	comp.Frames.Frames = append(comp.Frames.Frames, mob, res)

	return &comp, nil
}
