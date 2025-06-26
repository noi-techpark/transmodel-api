// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package comp

import "github.com/noi-techpark/go-netex"

type StFlightData struct {
	Operators           []netex.Operator
	StopPlaces          []netex.StopPlace
	ServiceCalendars    []netex.ServiceCalendar
	Routes              []netex.Route
	Lines               []netex.Line
	ScheduledStopPoints []netex.ScheduledStopPoint
	ServiceLinks        []netex.ServiceLink
	StopAssignments     []netex.PassengerStopAssignment
	JourneyPatterns     []netex.ServiceJourneyPattern
	VehicleJourneys     []netex.ServiceJourney
}

type StFlights interface {
	StFlights() (StFlightData, error)
}

func GetFlights(ps []StFlights) ([]netex.CompositeFrame, error) {
	ret := []netex.CompositeFrame{}

	apd := StFlightData{}

	for _, p := range ps {
		pd, err := p.StFlights()
		if err != nil {
			return ret, err
		}
		apd.Operators = append(apd.Operators, pd.Operators...)
		apd.StopPlaces = append(apd.StopPlaces, pd.StopPlaces...)
		apd.ServiceCalendars = append(apd.ServiceCalendars, pd.ServiceCalendars...)
		apd.Routes = append(apd.Routes, pd.Routes...)
		apd.Lines = append(apd.Lines, pd.Lines...)
		apd.ScheduledStopPoints = append(apd.ScheduledStopPoints, pd.ScheduledStopPoints...)
		apd.ServiceLinks = append(apd.ServiceLinks, pd.ServiceLinks...)
		apd.StopAssignments = append(apd.StopAssignments, pd.StopAssignments...)
		apd.JourneyPatterns = append(apd.JourneyPatterns, pd.JourneyPatterns...)
		apd.VehicleJourneys = append(apd.VehicleJourneys, pd.VehicleJourneys...)
	}

	ret = append(ret, compFlights(apd))

	return ret, nil
}

func compFlights(pd StFlightData) netex.CompositeFrame {
	ret := DefaultCompositFrame()
	ret.Id = CreateFrameId(netex.TypeCompositeFrameLineOffer, "FLIGHTS", "ita")
	ret.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeLineOffer)

	site := siteFrame()
	site.StopPlaces = &pd.StopPlaces
	ret.Frames.Frames = append(ret.Frames.Frames, &site)

	res := netex.ResourceFrame{}
	res.Id = CreateFrameId(netex.TypeResourceFrameCommon, "ita")
	res.Version = "1"
	res.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeCommon)
	res.Operators = &pd.Operators
	ret.Frames.Frames = append(ret.Frames.Frames, &res)

	ser := netex.ServiceFrame{}
	ser.Id = CreateFrameId(netex.TypeServiceFrameNetwork, "ita")
	ser.Version = "1"
	ser.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeNetwork)
	ser.JourneyPatterns = netex.JustSlice(pd.JourneyPatterns)
	ser.Lines = netex.JustSlice(pd.Lines)
	ser.Routes = netex.JustSlice(pd.Routes)
	ser.ScheduledStopPoints = netex.JustSlice(pd.ScheduledStopPoints)
	ser.ServiceLinks = netex.JustSlice(pd.ServiceLinks)
	ser.StopAssignments = netex.JustSlice(pd.StopAssignments)
	ret.Frames.Frames = append(ret.Frames.Frames, ser)

	cal := netex.ServiceCalendarFrame{}
	cal.Id = CreateFrameId(netex.TypeServiceCalendarFrameCalendar, "ita")
	cal.Version = "1"
	cal.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeCalendar)
	cal.ServiceCalendar = append(cal.ServiceCalendar, pd.ServiceCalendars...)
	ret.Frames.Frames = append(ret.Frames.Frames, cal)

	tim := netex.TimetableFrame{}
	tim.Id = CreateFrameId(netex.TypeTimetableFrameTimetable, "ita")
	tim.Version = "1"
	tim.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeTimetable)
	tim.VehicleJourneys = netex.JustSlice(pd.VehicleJourneys)
	ret.Frames.Frames = append(ret.Frames.Frames, tim)

	return ret
}
