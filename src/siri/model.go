// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

import (
	"time"
)

type Siri struct {
	Siri struct {
		Version         string
		ServiceDelivery ServiceDelivery
	}
}

func (s *Siri) AppencFcs(fcs []FacilityCondition) {
	s.Siri.ServiceDelivery.FacilityMonitoringDelivery.FacilityCondition = append(s.Siri.ServiceDelivery.FacilityMonitoringDelivery.FacilityCondition, fcs...)
}

func NewSiri() Siri {
	siri := Siri{}
	siri.Siri.Version = "2"
	sd := &siri.Siri.ServiceDelivery
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
}
