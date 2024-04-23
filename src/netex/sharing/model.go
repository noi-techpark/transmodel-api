// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package sharing

import (
	"encoding/xml"
	"opendatahub/sta-nap-export/netex"
)

type MobilityServiceFrame struct {
	XMLName       xml.Name `xml:"MobilityServiceFrame"`
	Id            string   `xml:"id,attr"`
	Version       string   `xml:"version,attr"`
	FrameDefaults struct {
		DefaultCurrency string
	}

	Fleets                         []Fleet                         `xml:"fleets"`
	ModesOfOperation               []VehicleSharing                `xml:"modesOfOperation"`
	MobilityServices               []VehicleSharingService         `xml:"mobilityServices"`
	MobilityServiceConstraintZones []MobilityServiceConstraintZone `xml:"mobilityServiceContraintZones"`
}

type ResourceFrame struct {
	XMLName        xml.Name `xml:"ResourceFrame"`
	Id             string   `xml:"id,attr"`
	Version        string   `xml:"version,attr"`
	TypeOfFrameRef netex.TypeOfFrameRef

	Vehicles      []Vehicle           `xml:"vehicles>Vehicle"`
	VehicleModels []VehicleModel      `xml:"models>VehicleModel"`
	CarModels     []CarModelProfile   `xml:"carModels>CarModelProfile"`
	CycleModels   []CycleModelProfile `xml:"cycleModels>CycleModelProfile"`
	Operators     []Operator          `xml:"operators>Operator"`
}

type Fleet struct {
	Id           string `xml:"id,attr"`
	Version      string `xml:"version,attr"`
	ValidBetween netex.ValidBetween
	Members      []netex.Ref `xml:"members"`
	OperatorRef  netex.Ref
}

type Vehicle struct {
	XMLName            xml.Name `xml:"Vehicle"`
	Id                 string   `xml:"id,attr"`
	Version            string   `xml:"version,attr"`
	ValidBetween       netex.ValidBetween
	Name               string
	ShortName          string
	RegistrationNumber string
	VehicleIdNumber    string
	PrivateCode        string
	OperatorRef        netex.Ref
	VehicleTypeRef     netex.Ref
}

type VehicleModel struct {
	Id             string `xml:"id,attr"`
	Version        string `xml:"version,attr"`
	ValidBetween   netex.ValidBetween
	Name           string
	Description    string
	Manufacturer   string
	VehicleTypeRef netex.Ref
}

type CycleModelProfile struct {
	Id        string `xml:"id,attr"`
	Version   string `xml:"version,attr"`
	ChildSeat string
	Battery   bool
	Lamps     bool
	Pump      bool
	Basket    bool
	Lock      bool
}

type CarModelProfile struct {
	Id              string `xml:"id,attr"`
	Version         string `xml:"version,attr"`
	ChildSeat       string
	Seats           uint16
	Doors           uint16
	Transmission    string
	CruiseControl   bool
	SatNav          bool
	AirConditioning bool
	Convertible     bool
	UsbPowerSocket  bool
	WinterTyres     bool
	Chains          bool
	TrailerHitch    bool
	RoofRack        bool
	CycleRack       bool
	SkiRack         bool
}

type Operator struct {
	Id             string `xml:"id,attr"`
	Version        string `xml:"version,attr"`
	PrivateCode    string
	Name           string
	ShortName      string
	LegalName      string
	TradingName    string
	ContactDetails struct {
		Email string
		Phone string
		URL   string
	}
	OrganizationType string
	Address          struct {
		Id          string `xml:"id,attr"`
		CountryName string
		Street      string
		Town        string
		PostCode    string
	}
	Departments any
}

type Submode struct {
	Id            string `xml:"id,attr"`
	Version       string `xml:"version,attr"`
	TransportMode string
	SelfDriveMode string
}

type VehicleSharing struct {
	Id       string    `xml:"id,attr"`
	Version  string    `xml:"version,attr"`
	Submodes []Submode `xml:"submodes"`
}

type VehicleSharingService struct {
	Id                string `xml:"id,attr"`
	Version           string `xml:"version,attr"`
	VehicleSharingRef netex.Ref
	FloatingVehicles  bool
	Fleets            []netex.Ref `xml:"fleets"`
}

type MobilityServiceConstraintZone struct {
	Id                string `xml:"id,attr"`
	Version           string `xml:"version,attr"`
	GmlPolygon        any    `xml:"http://www.opengis.net/gml Polygon"`
	VehicleSharingRef netex.Ref
}
