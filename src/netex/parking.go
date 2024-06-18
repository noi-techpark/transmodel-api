// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"fmt"
	"opendatahub/sta-nap-export/config"
	"strings"
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

func GetParking(ps []StParking) ([]CompositeFrame, error) {
	ret := []CompositeFrame{}

	apd := StParkingData{}

	for _, p := range ps {
		pd, err := p.StParking()
		if err != nil {
			return ret, err
		}
		apd.Parkings = append(apd.Parkings, pd.Parkings...)
		apd.Operators = append(apd.Operators, pd.Operators...)
	}

	ret = append(ret, compFrame(apd))

	return ret, nil
}

func siteFrame() SiteFrame {
	var site SiteFrame
	site.Id = CreateFrameId("SiteFrame_EU_PI_STOP", "ita")
	site.Version = "1"
	site.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_STOP")
	return site
}
