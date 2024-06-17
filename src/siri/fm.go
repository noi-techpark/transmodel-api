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

func (s *Siri) appencFcs(fcs []FacilityCondition) {
	s.Siri.ServiceDelivery.FacilityMonitoringDelivery.FacilityCondition = append(s.Siri.ServiceDelivery.FacilityMonitoringDelivery.FacilityCondition, fcs...)
}

func newSiri() Siri {
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
	req.StationTypes = []string{"ParkingStation", "BikeParking", "EChargingStation"}
	req.DataTypes = []string{"free", "number-available"}
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

func mapStatus(free int, partialThreshold int) string {
	switch {
	case free == 0:
		return "notAvailable"
	case free <= partialThreshold:
		return "partiallyAvailable"
	default:
		return "available"
	}
}

func Parking() (Siri, error) {
	siri := newSiri()

	os, err := odhParkingState()
	if err != nil {
		return siri, err
	}
	siri.appencFcs(mapParking(os))

	return siri, nil
}

func mapParking(latest []OdhLatest) []FacilityCondition {
	ret := []FacilityCondition{}

	for _, o := range latest {
		fc := FacilityCondition{}
		fc.FacilityRef = netex.CreateID("Parking", o.Scode)
		fc.MonitoredCounting.CountingType = "presentCount"

		switch o.Stype {
		case "BikeParking":
			fc.FacilityStatus.Status = mapStatus(o.MValue, 10)
			fc.MonitoredCounting.CountedFeatureUnit = "otherSpaces"
			fc.MonitoredCounting.Count = o.TotalPlaces - o.MValue
		case "EChargingStation":
			fc.FacilityStatus.Status = mapStatus(o.MValue, 1)
			fc.MonitoredCounting.CountedFeatureUnit = "devices"
			fc.MonitoredCounting.Count = o.Capacity - o.MValue
		case "ParkingStation":
			fc.FacilityStatus.Status = mapStatus(o.MValue, 10)
			fc.MonitoredCounting.CountedFeatureUnit = "bays"
			fc.MonitoredCounting.Count = o.Capacity - o.MValue
		}

		ret = append(ret, fc)
	}

	return ret
}

func Sharing() (Siri, error) {
	// station based bike
	// free floating bike
	// car sharing
	siri := newSiri()

	return siri, nil
}
