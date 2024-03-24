// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

type Siri struct {
	ResponseTimeStamp         string
	ProducerRef               string
	ResponseMessageIdentifier int64
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

func Fm(id string) {

}
