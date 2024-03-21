// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package sharing

import (
	"encoding/json"
	"log/slog"
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
	fleets        []Fleet
	vehicles      []Vehicle
	vehicleModels []VehicleModel
	carModels     []CarModelProfile
	cycleModels   []CycleModelProfile
	operators     []Operator
	modes         []VehicleSharing
	services      []VehicleSharingService
	constraints   []MobilityServiceConstraintZone
}

type SharingProvider interface {
	get() (SharingData, error)
}

func GetSharing() (MobilityServiceFrame, error) {
	return frame([]SharingProvider{&Bz{}})
}

func frame(ps []SharingProvider) (MobilityServiceFrame, error) {
	f := MobilityServiceFrame{}
	f.Id = "edp:IT:ITH10:ASDASDASDASDASD"
	f.Version = "1"

	for _, p := range ps {
		d, err := p.get()
		if err != nil {
			return f, err
		}
		f.Fleets = append(f.Fleets, d.fleets...)
		f.ModesOfOperation = append(f.ModesOfOperation, d.modes...)
		f.MobilityServices = append(f.MobilityServices, d.services...)
		f.MobilityServiceConstraintZones = append(f.MobilityServiceConstraintZones, d.constraints...)
	}

	return f, nil
}

func MkRef(tp string, id string) Ref {
	r := Ref{}
	r.Ref = id
	r.Version = "1"
	r.XMLName.Local = tp + "Ref"
	return r
}
