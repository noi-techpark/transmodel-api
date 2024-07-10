// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later
package provider

import (
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"
	"opendatahub/sta-nap-export/siri"
	"slices"
	"strings"
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

func maybeIdMatch(ids []string, prefix string) []string {
	return slices.DeleteFunc(ids, func(id string) bool { return !strings.HasPrefix(id, prefix) })
}

func filterFacilityConditions(conds []siri.FacilityCondition, ids []string) []siri.FacilityCondition {
	if len(ids) == 0 {
		return conds
	}
	return slices.DeleteFunc(conds, func(c siri.FacilityCondition) bool {
		for _, id := range ids {
			if c.FacilityRef == id {
				return false
			}
		}
		return true
	})
}

func apiBoundingBox(q siri.Query) string {
	lat := q.Lat()
	long := q.Lon()
	rad := q.MaxDistance() // we assume this is radius in meters
	if lat == 0 || long == 0 || rad == 0 {
		return ""
	}
	return ""

}
