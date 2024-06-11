// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package parking

import (
	"fmt"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"
	"strings"

	"golang.org/x/exp/maps"
)

type OdhParking struct {
	Scode       string
	Sname       string
	Sorigin     string
	Stype       string
	Scoordinate struct {
		X    float32
		Y    float32
		Srid uint32
	}
	Smetadata struct {
		StandardName string
		Capacity     int32
		TotalPlaces  int32 // Bikeparking specific capacity
		Municipality string
		Netex        struct {
			Type             string
			VehicleTypes     string
			Layout           string
			HazardProhibited bool `json:"hazard_prohibited"`
			Charging         bool
			Surveillance     bool
			Reservation      string `json:"reservation"`
		} `json:"netex_parking"`
	}
}
type OdhEcharging struct {
	Scode       string
	Sname       string
	Sorigin     string
	Scoordinate struct {
		X    float32
		Y    float32
		Srid uint32
	}
	Smetadata struct {
		State    string
		Capacity int32
	}
}

func originList() string {
	origins := netex.Cfg.ParkingOrigins()
	quoted := []string{}
	for _, o := range origins {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", o))
	}
	return strings.Join(quoted, ",")
}

func getOdhParking() ([]OdhParking, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"ParkingStation", "BikeParking"}
	req.Where = "sactive.eq.true"
	req.Where += fmt.Sprintf(",sorigin.in.(%s)", originList())
	// TODO: limit bounding box / polygon
	var res ninja.NinjaResponse[[]OdhParking]
	err := ninja.StationType(req, &res)
	return res.Data, err
}
func getOdhEcharging() ([]OdhEcharging, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"EchargingStation"}
	req.Where = "sactive.eq.true"
	//req.Where += fmt.Sprintf(",sorigin.in.(%s)", originList())
	// Rudimentary geographical limit
	req.Where += fmt.Sprintf(",scoordinate.bbi.(%s)", bboxSouthTyrol)
	var res ninja.NinjaResponse[[]OdhEcharging]
	err := ninja.StationType(req, &res)
	return res.Data, err
}

const bboxSouthTyrol = "10.368347,46.185535,12.551880,47.088826,4326"

func defEmpty(s string, d string) string {
	if s == "" {
		return d
	} else {
		return s
	}
}

func originButWithHacks(p OdhParking) string {
	// There are two different Operators that both have origin FBK (Trento and Rovereto)
	if p.Sorigin == "FBK" {
		return fmt.Sprintf("%s-%s", p.Sorigin, p.Smetadata.Municipality)
	}
	return p.Sorigin
}

func mapParking(os []OdhParking) ([]netex.Parking, []netex.Operator) {
	ops := make(map[string]netex.Operator)

	var ps []netex.Parking
	for _, o := range os {
		var p netex.Parking

		p.Id = netex.CreateID("Parking", o.Scode)
		p.Version = "1"
		p.ShortName = o.Sname
		// p.Centroid.Location.Precision = 1  not sure what this actually does, according to specification not needed?
		p.Centroid.Location.Longitude = o.Scoordinate.X
		p.Centroid.Location.Latitude = o.Scoordinate.Y
		p.GmlPolygon = nil
		op := netex.Cfg.GetOperator(originButWithHacks(o))
		ops[op.Id] = op
		p.OperatorRef = netex.MkRef("Operator", op.Id)

		p.Entrances = nil
		p.ParkingType = defEmpty(o.Smetadata.Netex.Type, "undefined")
		p.ParkingVehicleTypes = o.Smetadata.Netex.VehicleTypes
		p.ParkingLayout = defEmpty(o.Smetadata.Netex.Layout, "undefined")
		p.ProhibitedForHazardousMaterials.Set(o.Smetadata.Netex.HazardProhibited)
		p.RechargingAvailable = o.Smetadata.Netex.Charging
		p.Secure.Set(o.Smetadata.Netex.Surveillance)
		p.ParkingReservation = defEmpty(o.Smetadata.Netex.Reservation, "noReservations")
		p.ParkingProperties = nil

		if o.Stype == "BikeParking" {
			p.Name = o.Sname
			p.PrincipalCapacity = o.Smetadata.TotalPlaces
			p.TotalCapacity = o.Smetadata.TotalPlaces
		} else {
			p.Name = o.Smetadata.StandardName
			p.PrincipalCapacity = o.Smetadata.Capacity
			p.TotalCapacity = o.Smetadata.Capacity
		}

		ps = append(ps, p)
	}
	return ps, maps.Values(ops)
}
func mapEcharging(os []OdhEcharging) ([]netex.Parking, []netex.Operator) {
	ops := make(map[string]netex.Operator)

	var ps []netex.Parking
	for _, o := range os {
		var p netex.Parking

		p.Id = netex.CreateID("Parking", o.Scode)
		p.Version = "1"
		p.ShortName = o.Sname
		p.Centroid.Location.Longitude = o.Scoordinate.X
		p.Centroid.Location.Latitude = o.Scoordinate.Y
		p.GmlPolygon = nil
		op := netex.Cfg.GetOperator(o.Sorigin)
		ops[op.Id] = op
		p.OperatorRef = netex.MkRef("Operator", op.Id)

		p.Entrances = nil
		p.ParkingType = "roadside"
		p.ParkingVehicleTypes = ""
		p.ParkingLayout = "undefined"
		p.ProhibitedForHazardousMaterials.Ignore()
		p.RechargingAvailable = true
		p.Secure.Ignore()
		p.ParkingReservation = "reservationAllowed"
		p.ParkingProperties = nil

		p.Name = o.Sname
		p.PrincipalCapacity = o.Smetadata.Capacity
		p.TotalCapacity = o.Smetadata.Capacity

		ps = append(ps, p)
	}
	return ps, maps.Values(ops)
}

func compFrame(ps []netex.Parking, os []netex.Operator) netex.CompositeFrame {
	var ret netex.CompositeFrame
	ret.Defaults()
	ret.Id = netex.CreateFrameId("CompositeFrame_EU_PI_STOP_OFFER", "PARKING", "ita")
	ret.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_LINE_OFFER")

	site := siteFrame()
	ret.Frames.Frames = append(ret.Frames.Frames, &site)

	res := netex.ResourceFrame{}
	res.Id = netex.CreateFrameId("ResourceFrame_EU_PI_MOBILITY", "ita")
	res.Version = "1"
	res.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_COMMON")
	ret.Frames.Frames = append(ret.Frames.Frames, &res)

	site.Parkings = netex.Parkings{Parkings: ps}
	res.Operators = &os

	return ret
}

func GetParking() (netex.CompositeFrame, error) {
	var ret netex.CompositeFrame

	odh, err := getOdhParking()
	if err != nil {
		return ret, err
	}
	parkings, operators := mapParking(odh)

	eodh, err := getOdhEcharging()
	if err != nil {
		return ret, err
	}
	eparkings, eoperators := mapEcharging(eodh)

	parkings = append(parkings, eparkings...)
	operators = append(operators, eoperators...)

	ret = compFrame(parkings, operators)

	return ret, nil
}

func siteFrame() netex.SiteFrame {
	var site netex.SiteFrame
	site.Id = netex.CreateFrameId("SiteFrame_EU_PI_STOP", "ita")
	site.Version = "1"
	site.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_STOP")
	return site
}
