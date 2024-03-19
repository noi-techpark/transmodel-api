// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import "opendatahub/sta-nap-export/ninja"

type OdhSharing struct {
	Scode   string `json:"scode"`
	Sname   string `json:"sname"`
	Sorigin string `json:"sorigin"`
	Scoord  struct {
		X    float32 `json:"x"`
		Y    float32 `json:"y"`
		Srid uint32  `json:"srid"`
	} `json:"scoordinate"`
	Smeta struct {
		Company struct {
			ShortName string `json:"shortName"`
		} `json:"company"`
		LicensePlate string `json:"licensePlate"`
		Electric     bool   `json:"electric"`
		Brand        bool   `json:"brand"`
	} `json:"smetadata"`
}

func getOdhSharingBz() ([]OdhSharing, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"BikesharingStation"}
	req.Where = `sorigin.eq.BIKE_SHARING_BOLZANO`
	var res ninja.NinjaResponse[[]OdhSharing]
	err := ninja.StationType(req, &res)
	return res.Data, err
}
func getOdhSharingMe() ([]OdhSharing, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"Bicycle"}
	req.Where = `sorigin.eq.BIKE_SHARING_MERANO`
	var res ninja.NinjaResponse[[]OdhSharing]
	err := ninja.StationType(req, &res)
	return res.Data, err
}
func getOdhSharingPapin() ([]OdhSharing, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"BikesharingStation"}
	req.Where = `sorigin.eq.BIKE_SHARING_PAPIN`
	var res ninja.NinjaResponse[[]OdhSharing]
	err := ninja.StationType(req, &res)
	return res.Data, err
}
func getOdhSharingHal() ([]OdhSharing, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"CarsharingStation"}
	req.Where = `sorigin.eq.HAL-API`
	var res ninja.NinjaResponse[[]OdhSharing]
	err := ninja.StationType(req, &res)
	return res.Data, err
}
