// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package sharing

import (
	"encoding/xml"
	"time"
)

type MobilityServiceFrame struct {
	Id                             string                          `xml:"id,attr"`
	Version                        string                          `xml:"version,attr"`
	Fleets                         []Fleet                         `xml:"fleets"`
	ModesOfOperation               []VehicleSharing                `xml:"modesOfOperation"`
	MobilityServices               []VehicleSharingService         `xml:"mobilityServices"`
	MobilityServiceConstraintZones []MobilityServiceConstraintZone `xml:"mobilityServiceContraintZones"`
}

type Ref struct {
	XMLName xml.Name
	Ref     string `xml:"ref,attr"`
	Version string `xml:"version,attr"`
}

type ValidBetween struct {
	FromDate *time.Time `xml:"FromDate,omitempty"`
	ToDate   *time.Time `xml:"ToDate,omitempty"`
}

type Fleet struct {
	Id           string `xml:"id,attr"`
	Version      string `xml:"version,attr"`
	ValidBetween ValidBetween
	Members      []Ref `xml:"members"`
	OperatorRef  Ref
}

type Vehicle struct {
	XMLName            xml.Name `xml:"Vehicle"`
	Id                 string   `xml:"id,attr"`
	Version            string   `xml:"version,attr"`
	ValidBetween       ValidBetween
	Name               string
	ShortName          string
	RegistrationNumber string
	VehicleIdNumber    string
	PrivateCode        string
	OperatorRef        Ref
	VehicleTypeRef     Ref
}

type VehicleModel struct {
	Id             string `xml:"id,attr"`
	Version        string `xml:"version,attr"`
	ValidBetween   ValidBetween
	Name           string
	Description    string
	Manufacturer   string
	VehicleTypeRef Ref
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
	VehicleSharingRef Ref
	FloatingVehicles  bool
	Fleets            []Ref `xml:"fleets"`
}

type MobilityServiceConstraintZone struct {
	Id                string `xml:"id,attr"`
	Version           string `xml:"version,attr"`
	GmlPolygon        any    `xml:"gml:Polygon"`
	VehicleSharingRef Ref
}
