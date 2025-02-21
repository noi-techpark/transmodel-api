// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later
package netex

import (
	"encoding/xml"
	"time"
)

const NetexNamespace = "http://www.netex.org.uk/netex"

type NetexFrame struct {
	XMLName              xml.Name `xml:"http://www.netex.org.uk/netex PublicationDelivery"`
	Version              string   `xml:"version,attr"`
	NsXsi                string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation    string   `xml:"xsi:schemaLocation,attr"`
	PublicationTimestamp time.Time
	ParticipantRef       string
	Description          string
	DataObjects          []CompositeFrame `xml:"dataObjects>CompositeFrame"`
}

func NewNetexFrame() NetexFrame {
	n := NetexFrame{}
	n.Version = "1.0"
	n.NsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	n.XsiSchemaLocation = "http://www.netex.org.uk/netex https://raw.githubusercontent.com/5Tsrl/netex-italian-profile/main/xsd/NeTEx_publication_Lev4.xsd"
	n.PublicationTimestamp = time.Now()
	n.ParticipantRef = "RAP"
	n.Description = "Open Data Hub Netex export"
	return n
}

type Root struct {
	XMLName        xml.Name         `xml:"root"`
	CompositeFrame []CompositeFrame `xml:"CompositeFrame"`
}
type CompositeFrame struct {
	Id             string `xml:"id,attr"`
	Version        string `xml:"version,attr"`
	ValidBetween   ValidBetween
	TypeOfFrameRef Ref
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
	Frames struct {
		ResourceFrame        []ResourceFrame
		SiteFrame            []SiteFrame
		MobilityServiceFrame []MobilityServiceFrame
		ServiceFrame         []ServiceFrame
		ServiceCalendarFrame []ServiceCalendarFrame
		TimetableFrame       []TimetableFrame
		Frames               []any
	} `xml:"frames"`
}

type ResourceFrame struct {
	XMLName        xml.Name `xml:"ResourceFrame"`
	Id             string   `xml:"id,attr"`
	Version        string   `xml:"version,attr"`
	TypeOfFrameRef Ref

	Operators   *[]Operator          `xml:"organisations>Operator"`
	CarModels   *[]CarModelProfile   `xml:"vehicleModelProfiles>CarModelProfile"`
	CycleModels *[]CycleModelProfile `xml:"vehicleModelProfiles>CycleModelProfile"`
	Vehicles    *[]Vehicle           `xml:"vehicles>Vehicle"`
}

type SiteFrame struct {
	XMLName        xml.Name `xml:"SiteFrame"`
	Id             string   `xml:"id,attr"`
	Version        string   `xml:"version,attr"`
	TypeOfFrameRef Ref
	Parkings       []Parking   `xml:"parkings>Parking"`
	StopPlaces     []StopPlace `xml:"stopPlaces>StopPlace"`
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

// Avoid having the explicit namespace there for every Ref
func (r *Ref) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type ReRef Ref
	err := d.DecodeElement((*ReRef)(r), &start)
	if err == nil {
		r.XMLName.Space = ""
	}
	return err
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
	Members      []Ref `xml:"members>VehicleRef"`
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

type GmlPolygon struct {
	XMLName xml.Name `xml:"http://www.opengis.net/gml/3.2 Polygon"`
	Id      string   `xml:"id,attr"`
	Polygon string   `xml:",innerxml"`
}

type MobilityServiceConstraintZone struct {
	Id                string `xml:"id,attr"`
	Version           string `xml:"version,attr"`
	GmlPolygon        GmlPolygon
	VehicleSharingRef Ref
}

type Parking struct {
	XMLName                         xml.Name `xml:"Parking"`
	Id                              string   `xml:"id,attr"`
	Version                         string   `xml:"version,attr"`
	Name                            string
	ShortName                       string
	Centroid                        Centroid
	GmlPolygon                      any `xml:"http://www.opengis.net/gml/3.2 Polygon"`
	OperatorRef                     Ref
	Entrances                       any `xml:"entrances"`
	ParkingType                     string
	ParkingVehicleTypes             string
	ParkingLayout                   string
	PrincipalCapacity               int32
	TotalCapacity                   int32
	ProhibitedForHazardousMaterials N[bool]
	RechargingAvailable             N[bool]
	Secure                          N[bool]
	ParkingReservation              string
	ParkingProperties               any
}

type Centroid struct {
	Location struct {
		Longitude float32
		Latitude  float32
	}
}

type StopPlace struct {
	XMLName       xml.Name `xml:"StopPlace"`
	Id            string   `xml:"id,attr"`
	Version       string   `xml:"version,attr"`
	Name          string
	ShortName     string
	PrivateCode   string
	Centroid      Centroid
	AccessModes   string
	PublicCode    string
	TransportMode string
	StopPlaceType string
	Levels        []Level `xml:"levels>Level"`
	Quays         []Quay  `xml:"quays>Quay"`
}
type Quay struct {
	Id       string `xml:"id,attr"`
	Version  string `xml:"version,attr"`
	Name     string
	Centroid struct {
		Location struct {
			Longitude string `xml:"Longitude"`
			Latitude  string `xml:"Latitude"`
		}
	}
	LevelRef Ref
	QuayType string
}
type Level struct {
	Id         string `xml:"id,attr"`
	Version    string `xml:"version,attr"`
	Name       string
	PublicCode string
}
type ServiceCalendarFrame struct {
	XMLName         xml.Name `xml:"ServiceCalendarFrame"`
	Id              string   `xml:"id,attr"`
	Version         string   `xml:"version,attr"`
	TypeOfFrameRef  Ref
	ServiceCalendar []ServiceCalendar `xml:",omitempty"`
}
type ServiceCalendar struct {
	Id                 string `xml:"id,attr"`
	Version            string `xml:"version,attr"`
	Name               string
	FromDate           string
	ToDate             string
	DayTypes           []DayType            `xml:"dayTypes>DayType"`
	OperatingPeriods   []UicOperatingPeriod `xml:"operatingperiods>UicOperatingPeriod"`
	DayTypeAssignments []DayTypeAssignment  `xml:"dayTypeAssignments>DayTypeAssignment"`
}
type DayType struct {
	Id          string `xml:"id,attr"`
	Version     string `xml:"version,attr"`
	Name        string
	Description string
	Properties  struct {
		PropertyOfDay struct {
			DaysOfWeek   string
			HolidayTypes string
		}
	} `xml:"properties"`
}
type UicOperatingPeriod struct {
	Id           string `xml:"id,attr"`
	Version      string `xml:"version,attr"`
	FromDate     string
	ToDate       string
	ValidDayBits string
}
type DayTypeAssignment struct {
	Id                 string `xml:"id,attr"`
	Version            string `xml:"version,attr"`
	Order              string `xml:"order,attr"`
	OperatingPeriodRef Ref
	DayTypeRef         Ref
}
type Route struct {
	Id      string `xml:"id,attr"`
	Version string `xml:"version,attr"`
	Name    string
}
type Line struct {
	Id            string `xml:"id,attr"`
	Version       string `xml:"version,attr"`
	Name          string
	ShortName     string
	Description   string
	TransportMode string
	Url           string
	PublicCode    string
	PrivateCode   string
	OperatorRef   Ref
	Monitored     string
}
type ScheduledStopPoint struct {
	Id       string `xml:"id,attr"`
	Version  string `xml:"version,attr"`
	Name     string
	Location struct {
		Longitude string
		Latitude  string
		Altitude  string
		Precision string
	}
	ShortName   string
	Description string
	PublicCode  string
	PrivateCode string
}
type ServiceLink struct {
	Id           string `xml:"id,attr"`
	Version      string `xml:"version,attr"`
	Distance     string
	LineString   LineString `xml:"http://www.opengis.net/gml/3.2 LineString"`
	FromPointRef Ref
	ToPointRef   Ref
}

type LineString struct {
	Id      string `xml:"id,attr"`
	PosList string `xml:"posList"`
}

type PassengerStopAssignment struct {
	Order                 string `xml:"order,attr"`
	Id                    string `xml:"id,attr"`
	Version               string `xml:"version,attr"`
	ScheduledStopPointRef Ref
	StopPlaceRef          Ref
	QuayRef               Ref
}
type ServiceJourneyPattern struct {
	Id        string `xml:"id,attr"`
	Version   string `xml:"version,attr"`
	Name      string
	Distance  string
	RouteView struct {
		LineRef Ref
	}
	PointsInSequence []StopPointInJourneyPattern `xml:"pointsInSequence>StopPointInJourneyPattern"`
}

type StopPointInJourneyPattern struct {
	Id                    string `xml:"id,attr"`
	Version               string `xml:"version,attr"`
	Order                 string `xml:"order,attr"`
	ScheduledStopPointRef Ref
	OnwardServiceLinkRef  Ref
	ForAlighting          bool
	ForBoarding           bool
}

type ServiceFrame struct {
	XMLName             xml.Name `xml:"ServiceFrame"`
	Id                  string   `xml:"id,attr"`
	Version             string   `xml:"version,attr"`
	TypeOfFrameRef      Ref
	Routes              []Route                   `xml:"routes>Route"`
	Lines               []Line                    `xml:"lines>Line"`
	ScheduledStopPoints []ScheduledStopPoint      `xml:"scheduledStopPoints>ScheduledStopPoint"`
	ServiceLinks        []ServiceLink             `xml:"serviceLinks>ServiceLink"`
	StopAssignments     []PassengerStopAssignment `xml:"stopAssignments>PassengerStopAssignment"`
	JourneyPatterns     []ServiceJourneyPattern   `xml:"journeyPatterns>ServiceJourneyPattern"`
}

type TimetabledPassingTime struct {
	Id                           string `xml:"id,attr"`
	Version                      string `xml:"version,attr"`
	StopPointInJourneyPatternRef Ref
	ArrivalTime                  string
	ArrivalDayOffset             string
	DepartureTime                string
	DepartureDayOffset           string
}

type ServiceJourney struct {
	Id                       string `xml:"id,attr"`
	Version                  string `xml:"version,attr"`
	Name                     string
	Distance                 string
	TransportMode            string
	DepartureTime            string
	DepartureDayOffset       string
	JourneyDuration          string
	DayTypes                 []Ref `xml:"dayTypes>DayTypeRef"`
	ServiceJourneyPatternRef Ref
	OperatorRef              Ref
	PassingTimes             []TimetabledPassingTime `xml:"passingTimes>TimetabledPassingTime"`
}
type TimetableFrame struct {
	XMLName         xml.Name `xml:"TimetableFrame"`
	Id              string   `xml:"id,attr"`
	Version         string   `xml:"version,attr"`
	TypeOfFrameRef  Ref
	VehicleJourneys []ServiceJourney `xml:"vehicleJourneys>ServiceJourney"`
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

// Don't render value
func (n *N[t]) Ignore() {
	n.Maybe(nil)
}

func (n *N[any]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.set {
		return e.EncodeElement(n.v, start)
	}
	return nil
}
