// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"fmt"
	"opendatahub/sta-nap-export/ninja"
)

type NetexParking struct {
	Id        string `xml:"id"`
	Version   string `xml:"version"`
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

type odhParking struct {
	Scode       string `json:"scode"`
	Sname       string `json:"sname"`
	Sorigin     string `json:"sorigin"`
	Scoordinate struct {
		X    float32 `json:"x"`
		Y    float32 `json:"y"`
		Srid uint32  `json:"srid"`
	} `json:"scoordinate"`
	Smetadata struct {
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

func getOdhParking() ([]odhParking, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{"ParkingStation"}
	var res ninja.NinjaResponse[[]odhParking]
	err := ninja.StationType(req, &res)
	return res.Data, err
}

func mapToNetex(os []odhParking) []NetexParking {
	var ps []NetexParking
	for _, o := range os {
		var p NetexParking

		p.Id = fmt.Sprintf("IT:OpenDataHub:%s:%s", o.Sorigin, o.Scode)
		p.Name = o.Smetadata.StandardName
		p.ShortName = o.Sname
		// p.Centroid.Location.Precision = 1  not sure what this actually does, according to specification not needed?
		p.Centroid.Location.Longitude = o.Scoordinate.X
		p.Centroid.Location.Latitude = o.Scoordinate.Y
		p.GmlPolygon = nil
		p.OperatorRef = o.Sorigin
		p.Entrances = nil
		p.ParkingType = o.Smetadata.ParkingType
		p.ParkingVehicleTypes = o.Smetadata.ParkingVehicleTypes
		p.ParkingLayout = o.Smetadata.ParkingLayout
		p.PrincipalCapacity = ""
		p.TotalCapacity = o.Smetadata.Capacity
		p.ProhibitedForHazardousMaterials = o.Smetadata.ParkingProhibitions
		p.RechargingAvailable = o.Smetadata.ParkingProhibitions
		p.Secure = o.Smetadata.ParkingSurveillance
		p.ParkingReservation = o.Smetadata.ParkingReservation
		p.ParkingProperties = nil

		ps = append(ps, p)
	}
	return ps
}

func validateXml(ps []NetexParking) error {
	// TODO: everything
	return nil
}

func GetNetexParking() ([]NetexParking, error) {
	odh, err := getOdhParking()
	if err != nil {
		return nil, err
	}
	netex := mapToNetex(odh)

	err = validateXml(netex)
	if err != nil {
		return nil, err
	}

	return netex, nil
}
