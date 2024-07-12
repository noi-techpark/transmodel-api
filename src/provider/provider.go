// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later
package provider

import (
	"fmt"
	"opendatahub/transmodel-api/config"
	"opendatahub/transmodel-api/netex"
	"opendatahub/transmodel-api/ninja"
	"opendatahub/transmodel-api/siri"
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

var ParkingStatic = []netex.StParking{NewParkingGeneric(), NewParkingEcharging()}
var ParkingRt = []siri.FMProvider{NewParkingGeneric(), NewParkingEcharging()}
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
	return fmt.Sprintf(",scoordinate.dlt.(%f,%f,%f)", rad, lat, long)
}

func intersect(a []string, b []string) []string {
	return slices.DeleteFunc(a, func(s string) bool { return !slices.Contains(b, s) })
}
func quotedList(origins []string) string {
	quoted := []string{}
	for _, o := range origins {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", o))
	}
	return strings.Join(quoted, ",")
}

func filterOpOrigins(ops []string, origins []string) []string {
	if len(ops) == 0 {
		return origins
	}
	ret := []string{}
	for _, op := range ops {
		ret = append(ret, intersect(netex.GetOperatorOrigins(&config.Cfg, op), origins)...)
	}
	return ret
}
