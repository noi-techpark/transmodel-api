// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"fmt"
	"opendatahub/sta-nap-export/config"
	"opendatahub/sta-nap-export/ninja"
	"strings"

	"golang.org/x/exp/maps"
)

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

func ParkingOrigins() string {
	origins := config.Cfg.ParkingOrigins()
	quoted := []string{}
	for _, o := range origins {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", o))
	}
	return strings.Join(quoted, ",")
}
func getOdhEcharging() ([]OdhEcharging, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"EChargingStation"}
	req.Where = "sactive.eq.true"
	// Rudimentary geographical limit
	// req.Where += ",scoordinate.bbi.(10.368347,46.185535,12.551880,47.088826,4326)"
	req.Where += fmt.Sprintf(",sorigin.in.(%s)", ParkingOrigins())
	var res ninja.NinjaResponse[[]OdhEcharging]
	err := ninja.StationType(req, &res)
	return res.Data, err
}

func defEmpty(s string, d string) string {
	if s == "" {
		return d
	} else {
		return s
	}
}

func mapEcharging(os []OdhEcharging) ([]Parking, []Operator) {
	ops := make(map[string]Operator)

	var ps []Parking
	for _, o := range os {
		var p Parking

		p.Id = CreateID("Parking", o.Scode)
		p.Version = "1"
		p.ShortName = o.Sname
		p.Centroid.Location.Longitude = o.Scoordinate.X
		p.Centroid.Location.Latitude = o.Scoordinate.Y
		p.GmlPolygon = nil
		op := GetOperator(&config.Cfg, o.Sorigin)
		ops[op.Id] = op
		p.OperatorRef = MkRef("Operator", op.Id)

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

func compFrame(pd StParkingData) CompositeFrame {
	var ret CompositeFrame
	ret.Defaults()
	ret.Id = CreateFrameId("CompositeFrame_EU_PI_STOP_OFFER", "PARKING", "ita")
	ret.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_LINE_OFFER")

	site := siteFrame()
	ret.Frames.Frames = append(ret.Frames.Frames, &site)

	res := ResourceFrame{}
	res.Id = CreateFrameId("ResourceFrame_EU_PI_MOBILITY", "ita")
	res.Version = "1"
	res.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_COMMON")
	ret.Frames.Frames = append(ret.Frames.Frames, &res)

	site.Parkings = Parkings{Parkings: pd.Parkings}
	res.Operators = &pd.Operators

	return ret
}

type StParkingData struct {
	Parkings  []Parking
	Operators []Operator
}

type StParking interface {
	StParking() (StParkingData, error)
}

func GetParking(ps []StParking) (Root, error) {
	var ret Root

	apd := StParkingData{}

	for _, p := range ps {
		pd, err := p.StParking()
		if err != nil {
			return ret, err
		}
		apd.Parkings = append(apd.Parkings, pd.Parkings...)
		apd.Operators = append(apd.Operators, pd.Operators...)
	}

	ret.CompositeFrame = append(ret.CompositeFrame, compFrame(apd))

	return ret, nil
}

func siteFrame() SiteFrame {
	var site SiteFrame
	site.Id = CreateFrameId("SiteFrame_EU_PI_STOP", "ita")
	site.Version = "1"
	site.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_STOP")
	return site
}
