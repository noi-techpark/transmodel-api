// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

import (
	"encoding/xml"
	"time"
)

type Siri struct {
	Version           string `xml:"version,attr"`
	ServiceDelivery   ServiceDelivery
	XMLName           xml.Name `json:"-" xml:"Siri"`
	NsNetex           string   `json:"-" xml:"xmlns,attr"`
	NsXsi             string   `json:"-" xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string   `json:"-" xml:"xsi:schemaLocation,attr"`
}

func (s *Siri) AppencFcs(fcs []FacilityCondition) {
	s.ServiceDelivery.FacilityMonitoringDelivery.FacilityCondition = append(s.ServiceDelivery.FacilityMonitoringDelivery.FacilityCondition, fcs...)
}

func NewSiri() Siri {
	siri := Siri{}
	siri.Version = "2.1"
	siri.NsNetex = "http://www.siri.org.uk/siri"
	siri.NsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	siri.XsiSchemaLocation = "http://www.siri.org.uk/siri"
	sd := &siri.ServiceDelivery
	sd.defaults()
	sd.FacilityMonitoringDelivery.defaults()

	return siri
}

type DeliveryThingy struct {
	ResponseTimestamp string
	ProducerRef       string
}

func (s *DeliveryThingy) defaults() {
	s.ResponseTimestamp = time.Now().Format(time.RFC3339)
	s.ProducerRef = "RAP Alto Adige - Open Data Hub"
}

type ServiceDelivery struct {
	DeliveryThingy
	FacilityMonitoringDelivery FacilityMonitoringDelivery
}

type FacilityMonitoringDelivery struct {
	DeliveryThingy
	FacilityCondition []FacilityCondition
}

type MonitoredCounting struct {
	CountingType       string
	CountedFeatureUnit string
	Count              int
}
type FacilityUpdatedPosition struct {
	Longitude float32
	Latitude  float32
}

type FacilityCondition struct {
	FacilityRef    string
	FacilityStatus struct {
		Status string
	}
	MonitoredCounting       *MonitoredCounting       `json:",omitempty"`
	FacilityUpdatedPosition *FacilityUpdatedPosition `json:",omitempty"`
	Facility                *Facility                `json:",omitempty"`
}

type Facility struct {
	FacilityClass    string
	FacilityLocation struct {
		VehicleRef  string
		OperatorRef string
	}
}
