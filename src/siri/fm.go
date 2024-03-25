// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

import (
	"fmt"
	"log/slog"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"
	"sync"
	"time"
)

type Siri struct {
	ResponseTimeStamp         string
	ProducerRef               string
	ResponseMessageIdentifier uint64
	SubscriberRef             string
	SubscriptionRef           string
}

type SiriFM struct {
	Siri
	FacilityRef    string
	FacilityStatus struct {
		Status string
	}
	MonitoredCounting struct {
		CountingType          string
		CountedFeatureUnit    string `xml:"countedFeatureUnit"`
		TypeOfCountedFeatures struct {
			TypeOfValueCode string
			NameOfClass     string
		}
		Count uint16
	}
	FacilityUpdatedPosition struct {
		Longitude float32
		Latitude  float32
	}
}

type OdhLatest struct {
	TName      string          `json:"tname"`
	MPeriod    uint32          `json:"mperiod"`
	MValidTime ninja.NinjaTime `json:"mvalidtime"`
	MValue     uint16          `json:"mvalue"`
}

func Fm(id string) {
	// parse ID format, determine type
	// get parking or get sharing, depending on type
}

func odhParkingState(scode string) ([]OdhLatest, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = 10
	req.Repr = ninja.FlatNode
	req.StationTypes = []string{"ParkingStation"}
	req.DataTypes = []string{"free", "occupied"}
	req.Select = "tname,mperiod,mvalue,mvalidtime,scode"
	req.Where = fmt.Sprintf("scode.eq.\"%s\"", scode)
	var res ninja.NinjaResponse[[]OdhLatest]
	err := ninja.Latest(req, &res)
	if err != nil {
		slog.Error("Error retrieving parking state", "scode", scode, "err", err)
	}
	return res.Data, err
}

func odhSharingState(scode string) ([]OdhLatest, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = 10
	req.Repr = ninja.FlatNode
	req.StationTypes = []string{"BicycleSharing"}
	req.DataTypes = []string{"free", "occupied"}
	req.Select = "tname,mperiod,mvalue,mvalidtime,scode"
	req.Where = fmt.Sprintf("scode.eq.\"%s\"", scode)
	var res ninja.NinjaResponse[[]OdhLatest]
	err := ninja.Latest(req, &res)
	if err != nil {
		slog.Error("Error retrieving parking state", "scode", scode, "err", err)
	}
	return res.Data, err
}

type id struct {
	id uint64
	m  sync.Mutex
}

var responseId id

func (r *id) next() uint64 {
	// TODO: Must persist between sessions?
	r.m.Lock()
	defer r.m.Unlock()
	r.id++
	return r.id
}

func (s *SiriFM) defaults() {
	s.ResponseTimeStamp = time.Now().Format(time.RFC3339)
	s.ProducerRef = "RAP Alto Adige - Open Data Hub"
	s.ResponseMessageIdentifier = responseId.next()
}

func mapParkingStatus(free uint16) string {
	switch {
	case free == 0:
		return "notAvailable"
	case free <= 10:
		return "partiallyAvailable"
	default:
		return "available"
	}
}

func Parking(scode string) (SiriFM, error) {
	os, err := odhParkingState(scode)
	if err != nil {
		return SiriFM{}, err
	}
	var free, occupied OdhLatest
	for _, o := range os {
		switch o.TName {
		case "free":
			free = o
		case "occupied":
			occupied = o
		}
	}

	s := SiriFM{}
	s.defaults()
	s.FacilityRef = netex.CreateID("Parking", scode)
	s.FacilityStatus.Status = mapParkingStatus(free.MValue)
	s.MonitoredCounting.CountingType = "presentCount"
	s.MonitoredCounting.CountedFeatureUnit = "vehicles"
	s.MonitoredCounting.TypeOfCountedFeatures.TypeOfValueCode = "car"
	s.MonitoredCounting.TypeOfCountedFeatures.NameOfClass = "car"
	s.MonitoredCounting.Count = occupied.MValue

	return s, nil
}

func sharing(scode string) (SiriFM, error) {
	os, err := odhParkingState(scode)
	if err != nil {
		return SiriFM{}, err
	}
	var free, occupied OdhLatest
	for _, o := range os {
		switch o.TName {
		case "free":
			free = o
		case "occupied":
			occupied = o
		}
	}

	s := SiriFM{}
	s.defaults()
	s.FacilityRef = netex.CreateID("Parking", scode)
	s.FacilityStatus.Status = mapParkingStatus(free.MValue)
	s.MonitoredCounting.CountingType = "presentCount"
	s.MonitoredCounting.CountedFeatureUnit = "vehicles"
	s.MonitoredCounting.TypeOfCountedFeatures.TypeOfValueCode = "car"
	s.MonitoredCounting.TypeOfCountedFeatures.NameOfClass = "car"
	s.MonitoredCounting.Count = occupied.MValue

	return s, nil
}
