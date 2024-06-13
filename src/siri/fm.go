// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

import (
	"fmt"
	"log/slog"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/netex/parking"
	"opendatahub/sta-nap-export/ninja"
	"time"
)

type Siri struct {
	Siri struct {
		Version         string
		ServiceDelivery ServiceDelivery
	}
}

type DeliveryThingy struct {
	ResponseTimestamp string
	ProducerRef       string
}
type ServiceDelivery struct {
	DeliveryThingy
	FacilityMonitoringDelivery FacilityMonitoringDelivery
}

type FacilityMonitoringDelivery struct {
	DeliveryThingy
	FacilityCondition []FacilityCondition
}

type FacilityCondition struct {
	FacilityRef    string
	FacilityStatus struct {
		Status string
	}
	MonitoredCounting struct {
		CountingType       string
		CountedFeatureUnit string `xml:"countedFeatureUnit"`
		Count              int
	}
}

type OdhLatest struct {
	MPeriod     int             `json:"mperiod"`
	MValidTime  ninja.NinjaTime `json:"mvalidtime"`
	MValue      int             `json:"mvalue"`
	Scode       string          `json:"scode"`
	Stype       string          `json:"stype"`
	Capacity    int             `json:"smetadata.capacity"`
	TotalPlaces int             `json:"smetadata.totalPlaces"`
}

func odhParkingState() ([]OdhLatest, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.Repr = ninja.FlatNode
	req.StationTypes = []string{"ParkingStation", "BikeParking"}
	req.DataTypes = []string{"free"}
	req.Select = "mperiod,mvalue,mvalidtime,scode,stype,smetadata.capacity,smetadata.totalPlaces"
	req.Where = "sactive.eq.true"
	req.Where += fmt.Sprintf(",sorigin.in.(%s)", parking.ParkingOrigins())
	var res ninja.NinjaResponse[[]OdhLatest]
	err := ninja.Latest(req, &res)
	if err != nil {
		slog.Error("Error retrieving parking state", "err", err)
	}
	return res.Data, err
}

func (s *DeliveryThingy) defaults() {
	s.ResponseTimestamp = time.Now().Format(time.RFC3339)
	s.ProducerRef = "RAP Alto Adige - Open Data Hub"
}

func mapParkingStatus(free int) string {
	switch {
	case free == 0:
		return "notAvailable"
	case free <= 10:
		return "partiallyAvailable"
	default:
		return "available"
	}
}

func Parking() (Siri, error) {
	siri := Siri{}
	siri.Siri.Version = "2"
	sd := &siri.Siri.ServiceDelivery
	sd.defaults()
	sd.FacilityMonitoringDelivery.defaults()

	os, err := odhParkingState()
	if err != nil {
		return siri, err
	}

	for _, o := range os {
		fc := FacilityCondition{}
		fc.FacilityRef = netex.CreateID("Parking", o.Scode)
		fc.FacilityStatus.Status = mapParkingStatus(o.MValue)
		fc.MonitoredCounting.CountingType = "presentCount"

		if o.Stype == "BikeParking" {
			fc.MonitoredCounting.CountedFeatureUnit = "otherSpaces"
			fc.MonitoredCounting.Count = o.TotalPlaces - o.MValue
		} else {
			fc.MonitoredCounting.CountedFeatureUnit = "bays"
			fc.MonitoredCounting.Count = o.Capacity - o.MValue
		}

		sd.FacilityMonitoringDelivery.FacilityCondition = append(sd.FacilityMonitoringDelivery.FacilityCondition, fc)
	}

	return siri, nil
}
