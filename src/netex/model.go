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
		//} `xml:"codespaces"`
	} `xml:"-"`
	FrameDefaults struct {
		DefaultCodespaceRef Ref
	} `xml:"-"`
	Frames struct{ Frames []any } `xml:"frames"`
}

type ResourceFrame struct {
	XMLName        xml.Name `xml:"ResourceFrame"`
	Id             string   `xml:"id,attr"`
	Version        string   `xml:"version,attr"`
	TypeOfFrameRef TypeOfFrameRef

	Operators   *[]Operator          `xml:"organisations>Operator"`
	CarModels   *[]CarModelProfile   `xml:"vehicleModelProfiles>CarModelProfile"`
	CycleModels *[]CycleModelProfile `xml:"vehicleModelProfiles>CycleModelProfile"`
	Vehicles    *[]Vehicle           `xml:"vehicles>Vehicle"`
}

type SiteFrame struct {
	XMLName        xml.Name `xml:"SiteFrame"`
	Id             string   `xml:"id,attr"`
	Version        string   `xml:"version,attr"`
	TypeOfFrameRef TypeOfFrameRef
	Parkings       Parkings `xml:"parkings,omitempty"`
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
	OrganisationType string
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

	Fleets                         []Fleet                         `xml:"fleets>Fleet"`
	ModesOfOperation               []VehicleSharing                `xml:"modesOfOperation>VehicleSharing"`
	MobilityServices               []VehicleSharingService         `xml:"mobilityServices>VehicleSharingService"`
	MobilityServiceConstraintZones []MobilityServiceConstraintZone `xml:"mobilityServiceConstraintZones>MobilityServiceConstraintZone"`
}

type Fleet struct {
	Id           string `xml:"id,attr"`
	Version      string `xml:"version,attr"`
	ValidBetween ValidBetween
	Members      Members
	OperatorRef  Ref
}

// structure is members>VehicleRef[], we need this hack to not render the members tag if vehicleref list is empty
type Members []Ref

func (m *Members) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(*m) == 0 {
		return nil
	}
	members := xml.StartElement{Name: xml.Name{Local: "members"}}
	e.EncodeToken(members)
	for _, r := range *m {
		if err := e.EncodeElement(r, xml.StartElement{Name: xml.Name{Local: "VehicleRef"}}); err != nil {
			return err
		}
	}
	e.EncodeToken(members.End())
	return nil
}

type Vehicle struct {
	XMLName            xml.Name `xml:"Vehicle"`
	Id                 string   `xml:"id,attr"`
	Version            string   `xml:"version,attr"`
	ValidBetween       ValidBetween
	Name               string
	ShortName          string
	RegistrationNumber string
	OperationalNumber  string
	PrivateCode        string
	OperatorRef        Ref
	VehicleTypeRef     Ref
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
	Seats           uint8
	Doors           uint8
	Transmission    string
	CruiseControl   bool
	SatNav          bool
	AirConditioning bool
	Convertible     bool
	UsbPowerSockets bool
	WinterTyres     bool
	Chains          bool
	TrailerHitch    bool
	RoofRack        bool
	CycleRack       bool
	SkiRack         bool
}

type Submode struct {
	Id               string `xml:"id,attr"`
	Version          string `xml:"version,attr"`
	TransportMode    string
	SelfDriveSubmode string
}

type VehicleSharing struct {
	Id       string    `xml:"id,attr"`
	Version  string    `xml:"version,attr"`
	Submodes []Submode `xml:"submodes>Submode"`
}

type VehicleSharingService struct {
	Id                string `xml:"id,attr"`
	Version           string `xml:"version,attr"`
	VehicleSharingRef Ref
	FloatingVehicles  bool
	Fleets            []Ref `xml:"fleets>FleetRef"`
}

// Some hackery to always get the "gml" namespace prefix for our xmlns
// The hardcoded polygons use the gml prefix, so we have to bind the namespace to it
type GmlPolygon struct {
	XMLName xml.Name `xml:"gml:Polygon"`
	Id      string   `xml:"gml:id,attr"`
	Polygon string   `xml:",innerxml"`
	Xmlns   string   `xml:"xmlns:gml,attr"`
}

func (g *GmlPolygon) SetPoly(p string) {
	g.Polygon = p
	g.Xmlns = "http://www.opengis.net/gml/3.2"
}

type MobilityServiceConstraintZone struct {
	Id                string `xml:"id,attr"`
	Version           string `xml:"version,attr"`
	GmlPolygon        GmlPolygon
	VehicleSharingRef Ref
}

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
		}
	}
	GmlPolygon                      any `xml:"gml:Polygon"`
	OperatorRef                     Ref
	Entrances                       any `xml:"entrances"`
	ParkingType                     string
	ParkingVehicleTypes             string
	ParkingLayout                   string
	PrincipalCapacity               int32
	TotalCapacity                   int32
	ProhibitedForHazardousMaterials N[bool]
	RechargingAvailable             bool
	Secure                          N[bool]
	ParkingReservation              string
	ParkingProperties               any
}

// Nullable wrapper. If not explicitly set, it doesn't render in xml
type N[T any] struct {
	set bool
	v   *T
}

// Sets value and makes it render
func (n *N[t]) Set(v t) {
	n.set = true
	n.v = &v
}

// like set, but if value is nil, don't render
func (n *N[t]) Maybe(v *t) {
	n.set = v == nil
	n.v = v
}

// Sets value and makes it render
func (n *N[t]) Ignore() {
	n.Maybe(nil)
}

func (n *N[any]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.set {
		return e.EncodeElement(n.v, start)
	}
	return nil
}
