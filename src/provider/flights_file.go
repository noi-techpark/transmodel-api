// SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: CC0-1.0

package provider

import (
	"encoding/xml"
	"opendatahub/transmodel-api/netex"
	"os"
)

type FlightsSkyalps struct {
	path string
}

func NewFlightsSkyalps(path string) FlightsSkyalps {
	return FlightsSkyalps{path: path}
}

func (fs FlightsSkyalps) StFlights() (netex.StFlightData, error) {
	ret := netex.StFlightData{}
	f, err := os.ReadFile(fs.path)
	if err != nil {
		return ret, err
	}
	var n netex.NetexFrame
	err = xml.Unmarshal(f, &n)
	if err != nil {
		return ret, err
	}

	for _, do := range n.DataObjects {
		for _, rf := range do.Frames.ResourceFrame {
			ret.Operators = append(ret.Operators, *rf.Operators...)
		}
		for _, sf := range do.Frames.SiteFrame {
			ret.StopPlaces = append(ret.StopPlaces, sf.StopPlaces...)
		}
		for _, sf := range do.Frames.ServiceFrame {
			ret.JourneyPatterns = append(ret.JourneyPatterns, sf.JourneyPatterns...)
			ret.Lines = append(ret.Lines, sf.Lines...)
			ret.Routes = append(ret.Routes, sf.Routes...)
			ret.ScheduledStopPoints = append(ret.ScheduledStopPoints, sf.ScheduledStopPoints...)
			ret.ServiceLinks = append(ret.ServiceLinks, sf.ServiceLinks...)
			ret.StopAssignments = append(ret.StopAssignments, sf.StopAssignments...)
		}
		for _, cf := range do.Frames.ServiceCalendarFrame {
			ret.ServiceCalendars = append(ret.ServiceCalendars, cf.ServiceCalendar...)
		}
		for _, tf := range do.Frames.TimetableFrame {
			ret.VehicleJourneys = append(ret.VehicleJourneys, tf.VehicleJourneys...)
		}
	}

	return ret, nil
}
