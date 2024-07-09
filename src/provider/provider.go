// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later
package provider

import (
	"fmt"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"
	"opendatahub/sta-nap-export/siri"
)

func FetchOdhStations[T any](tp string, origin string) (T, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{tp}
	req.Where = `sorigin.eq.` + origin + `,sactive.eq.true`
	var res ninja.NinjaResponse[T]
	err := ninja.StationType(req, &res)
	return res.Data, err
}

var ParkingStatic = []netex.StParking{&ParkingGeneric{}, ParkingEcharging{}}
var ParkingRt = []siri.FMProvider{&ParkingGeneric{}, ParkingEcharging{}}
var SharingBikesStatic = []netex.StSharing{NewBikeBz(), NewBikeMe(), &BikePapin{}}
var SharingCarsStatic = []netex.StSharing{NewCarSharingHal()}
var SharingRt = []siri.FMProvider{NewBikeBz(), NewBikeMe(), NewCarSharingHal()}

func filterIDs(ids []string, idPrefix string, odhField string) string {
	ret := ""
	for _, f := range ids {
		var scode string
		_, err := fmt.Sscanf(f, fmt.Sprintf("%s:%%s", idPrefix), &scode)
		if err == nil {
			ret += fmt.Sprintf(",%s.eq.\"%s\"", odhField, scode)
		}
	}
	return ret
}
