// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package parking

import (
	"encoding/xml"
	"fmt"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"

	"golang.org/x/exp/maps"
)

type Parkings struct {
	XMLName  xml.Name `xml:"parkings"`
	Parkings []Parking
}

type Parking struct {
	XMLName   xml.Name `xml:"Parking"`
	Id        string   `xml:"id,attr"`
	Version   string   `xml:"version,attr"`
	Name      string
	ShortName string
	Centroid  struct {
		Location struct {
			Longitude float32
			Latitude  float32
			//Precision int8
		}
	}
	GmlPolygon                      any `xml:"gml:Polygon"`
	OperatorRef                     netex.Ref
	Entrances                       any `xml:"entrances"`
	ParkingType                     string
	ParkingVehicleTypes             string
	ParkingLayout                   string
	PrincipalCapacity               int32
	TotalCapacity                   int32
	ProhibitedForHazardousMaterials bool
	RechargingAvailable             bool
	Secure                          bool
	ParkingReservation              string
	ParkingProperties               any
}

type OdhParking struct {
	Scode   string `json:"scode"`
	Sname   string `json:"sname"`
	Sorigin string `json:"sorigin"`
	Scoord  struct {
		X    float32 `json:"x"`
		Y    float32 `json:"y"`
		Srid uint32  `json:"srid"`
	} `json:"scoordinate"`
	Smeta struct {
		StandardName string `json:"standard_name"`
		Capacity     int32  `json:"capacity"`
		Municipality string `json:"municipality"`
		Netex        struct {
			Type             string `json:"type"`
			VehicleTypes     string `json:"vehicletypes"`
			Layout           string `json:"layout"`
			HazardProhibited bool   `json:"hazard_prohibited"`
			Charging         bool   `json:"charging"`
			Surveillance     bool   `json:"surveillance"`
			Reservation      string `json:"reservation"`
		} `json:"netex_parking"`
	} `json:"smetadata"`
}

func getOdhParking() ([]OdhParking, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"ParkingStation"}
	req.Where = "sactive.eq.true"
	// TODO: limit bounding box / polygon
	var res ninja.NinjaResponse[[]OdhParking]
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

func originButWithHacks(p OdhParking) string {
	// There are two different Operators that both have origin FBK (Trento and Rovereto)
	if p.Sorigin == "FBK" {
		return fmt.Sprintf("%s-%s", p.Sorigin, p.Smeta.Municipality)
	}
	return p.Sorigin
}

func mapToNetex(os []OdhParking) ([]Parking, []netex.Operator) {
	ops := make(map[string]netex.Operator)

	var ps []Parking
	for _, o := range os {
		var p Parking

		p.Id = netex.CreateID("Parking", o.Scode)
		p.Version = "1"
		p.Name = o.Smeta.StandardName
		p.ShortName = o.Sname
		// p.Centroid.Location.Precision = 1  not sure what this actually does, according to specification not needed?
		p.Centroid.Location.Longitude = o.Scoord.X
		p.Centroid.Location.Latitude = o.Scoord.Y
		p.GmlPolygon = nil
		op := netex.GetOperator(originButWithHacks(o))
		ops[op.Id] = op
		p.OperatorRef = netex.MkRef("Operator", op.Id)

		p.Entrances = nil
		p.ParkingType = defEmpty(o.Smeta.Netex.Type, "undefined")
		p.ParkingVehicleTypes = o.Smeta.Netex.VehicleTypes
		p.ParkingLayout = defEmpty(o.Smeta.Netex.Layout, "undefined")
		p.PrincipalCapacity = o.Smeta.Capacity
		p.TotalCapacity = o.Smeta.Capacity
		p.ProhibitedForHazardousMaterials = o.Smeta.Netex.HazardProhibited
		p.RechargingAvailable = o.Smeta.Netex.Charging
		p.Secure = o.Smeta.Netex.Surveillance
		p.ParkingReservation = defEmpty(o.Smeta.Netex.Reservation, "noReservations")
		p.ParkingProperties = nil

		ps = append(ps, p)
	}
	return ps, maps.Values(ops)
}

func validateXml(p any) error {
	// TODO: everything
	return nil
}

func compFrame(ps []Parking, os []netex.Operator) netex.CompositeFrame {
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

	site.Parkings = Parkings{Parkings: ps}
	res.Operators = &os

	return ret
}

func GetParking() (netex.CompositeFrame, error) {
	var ret netex.CompositeFrame

	odh, err := getOdhParking()
	if err != nil {
		return ret, err
	}

	parkings, operators := mapToNetex(odh)
	ret = compFrame(parkings, operators)

	err = validateXml(ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func siteFrame() netex.SiteFrame {
	var site netex.SiteFrame
	site.Id = netex.CreateFrameId("SiteFrame_EU_PI_STOP", "ita")
	site.Version = "1"
	site.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_STOP")
	return site
}
