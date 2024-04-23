// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package parking

import (
	"encoding/xml"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"
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
	OperatorRef                     string
	Entrances                       any `xml:"entrances"`
	ParkingType                     string
	ParkingVehicleTypes             string
	ParkingLayout                   string
	PrincipalCapacity               string
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
		StandardName        string `json:"standard_name"`
		ParkingType         string `json:"parkingtype"`
		ParkingVehicleTypes string `json:"parkingvehicletypes"`
		ParkingLayout       string `json:"parkinglayout"`
		Capacity            int32  `json:"capacity"`
		ParkingProhibitions bool   `json:"parkingprohibitions"`
		ParkingCharging     bool   `json:"parkingcharging"`
		ParkingSurveillance bool   `json:"parkingsurveillance"`
		ParkingReservation  string `json:"parkingreservation"`
	} `json:"smetadata"`
}

func getOdhParking() ([]OdhParking, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"ParkingStation"}
	// TODO: limit bounding box / polygon
	var res ninja.NinjaResponse[[]OdhParking]
	err := ninja.StationType(req, &res)
	return res.Data, err
}

func mapToNetex(os []OdhParking) []Parking {
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
		p.OperatorRef = o.Sorigin
		p.Entrances = nil
		p.ParkingType = o.Smeta.ParkingType
		p.ParkingVehicleTypes = o.Smeta.ParkingVehicleTypes
		p.ParkingLayout = o.Smeta.ParkingLayout
		p.PrincipalCapacity = ""
		p.TotalCapacity = o.Smeta.Capacity
		p.ProhibitedForHazardousMaterials = o.Smeta.ParkingProhibitions
		p.RechargingAvailable = o.Smeta.ParkingProhibitions
		p.Secure = o.Smeta.ParkingSurveillance
		p.ParkingReservation = o.Smeta.ParkingReservation
		p.ParkingProperties = nil

		ps = append(ps, p)
	}
	return ps
}

func validateXml(p any) error {
	// TODO: everything
	return nil
}

func GetParking() (netex.CompositeFrame, error) {
	var ret netex.CompositeFrame
	ret.Defaults()
	ret.Id = netex.CreateFrameId("CompositeFrame_EU_PI_STOP_OFFER", "EPIP", "ita")
	ret.TypeOfFrameRef = netex.MkTypeOfFrameRef("EU_PI_LINE_OFFER")

	site := siteFrame()
	ret.Frames.Frames = append(ret.Frames.Frames, &site)

	odh, err := getOdhParking()
	if err != nil {
		return ret, err
	}

	parkings := mapToNetex(odh)
	site.Parkings = Parkings{Parkings: parkings}

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
