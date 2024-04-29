// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later
package netex

import (
	"encoding/xml"
	"time"
)

type CompositeFrame struct {
	XMLName        xml.Name `xml:"CompositeFrame"`
	Id             string   `xml:"id,attr"`
	Version        string   `xml:"version,attr"`
	ValidBetween   ValidBetween
	TypeOfFrameRef TypeOfFrameRef
	Codespaces     struct {
		Codespace struct {
			Id          string `xml:"id,attr"`
			Xmlns       string
			XmlnsUrl    string
			Description string
		}
	} `xml:"codespaces"`
	FrameDefaults struct {
		DefaultCodespaceRef Ref
	}
	Frames struct{ Frames []any } `xml:"frames"`
}

type ResourceFrame struct {
	XMLName        xml.Name `xml:"ResourceFrame"`
	Id             string   `xml:"id,attr"`
	Version        string   `xml:"version,attr"`
	TypeOfFrameRef TypeOfFrameRef

	Vehicles      []Vehicle           `xml:"vehicles>Vehicle"`
	VehicleModels []VehicleModel      `xml:"models>VehicleModel"`
	CarModels     []CarModelProfile   `xml:"carModels>CarModelProfile"`
	CycleModels   []CycleModelProfile `xml:"cycleModels>CycleModelProfile"`
	Operators     []Operator          `xml:"operators>Operator"`
}
type SiteFrame struct {
	XMLName        xml.Name `xml:"SiteFrame"`
	Id             string   `xml:"id,attr"`
	Version        string   `xml:"version,attr"`
	TypeOfFrameRef TypeOfFrameRef
	Parkings       any `xml:"parkings,omitempty"`
}

func (c *CompositeFrame) Defaults() {
	c.Version = "1"
	c.ValidBetween.AYear()
	c.Codespaces.Codespace.Id = "ita"
	c.Codespaces.Codespace.Xmlns = "ita"
	c.Codespaces.Codespace.XmlnsUrl = "http://www.ita.it"
	c.Codespaces.Codespace.Description = "Italian Profile"
	c.FrameDefaults.DefaultCodespaceRef.Ref = "ita"
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
		Url   string
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

type Ref struct {
	XMLName xml.Name
	Ref     string `xml:"ref,attr"`
	Version string `xml:"version,attr,omitempty"`
}
type TypeOfFrameRef struct {
	XMLName xml.Name
	Ref     string `xml:"ref,attr"`
	Version string `xml:"versionRef,attr,omitempty"`
}

type ValidBetween struct {
	FromDate time.Time
	ToDate   time.Time
}

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
	GmlPolygon        any    `xml:"http://www.opengis.net/gml Polygon"`
	VehicleSharingRef Ref
}
