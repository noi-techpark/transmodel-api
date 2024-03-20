// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

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

func bikeSharingBz() ([]OdhMobility[metaAny], error) {
	return odhMob[[]OdhMobility[metaAny]]("BikesharingStation", "BIKE_SHARING_BOLZANO")
}
func bikeBz() ([]OdhMobility[metaAny], error) {
	return odhMob[[]OdhMobility[metaAny]]("Bicycle", "BIKE_SHARING_BOLZANO")
}
func bikeMe() ([]OdhMobility[metaAny], error) {
	return odhMob[[]OdhMobility[metaAny]]("Bicycle", "BIKE_SHARING_MERANO")
}
func bikeSharingPapin() ([]OdhMobility[metaAny], error) {
	return odhMob[[]OdhMobility[metaAny]]("BikesharingStation", "BIKE_SHARING_PAPIN")
}
func carSharingHal() ([]OdhMobility[metaAny], error) {
	return odhMob[[]OdhMobility[metaAny]]("CarsharingStation", "HAL-API")
}

func getFleet() (Fleet, error) {
	ret := Fleet{}
	bz, err := bikeSharingBz()
	if err != nil {
		return ret, err
	}
	me, err := bikeMe()
	if err != nil {
		return ret, err
	}
	// carsharing HAL CarsharingCar
	car, err := bikeMe()
	if err != nil {
		return ret, err
	}
	_ = bz
	_ = me
	_ = car

	return ret, nil
}
