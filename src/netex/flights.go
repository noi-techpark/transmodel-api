// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

type StFlightData struct {
	Operators           []Operator
	StopPlaces          []StopPlace
	ServiceCalendars    []ServiceCalendar
	Routes              []Route
	Lines               []Line
	ScheduledStopPoints []ScheduledStopPoint
	ServiceLinks        []ServiceLink
	StopAssignments     []PassengerStopAssignment
	JourneyPatterns     []ServiceJourneyPattern
	VehicleJourneys     []ServiceJourney
}

type StFlights interface {
	StFlights() (StFlightData, error)
}

func GetFlights(ps []StFlights) ([]CompositeFrame, error) {
	ret := []CompositeFrame{}

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

func compFlights(pd StFlightData) CompositeFrame {
	var ret CompositeFrame
	ret.Defaults()
	ret.Id = CreateFrameId("CompositeFrame_EU_PI_STOP_OFFER", "FLIGHTS", "ita")
	ret.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_LINE_OFFER")

	site := siteFrame()
	site.StopPlaces = pd.StopPlaces
	ret.Frames.Frames = append(ret.Frames.Frames, &site)

	res := ResourceFrame{}
	res.Id = CreateFrameId("ResourceFrame_EU_PI_MOBILITY", "ita")
	res.Version = "1"
	res.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_COMMON")
	res.Operators = &pd.Operators
	ret.Frames.Frames = append(ret.Frames.Frames, &res)

	ser := ServiceFrame{}
	ser.Id = CreateFrameId("ServiceFrame_EU_PI_NETWORK", "ita")
	ser.Version = "1"
	ser.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_NETWORK")
	ser.JourneyPatterns = append(ser.JourneyPatterns, pd.JourneyPatterns...)
	ser.Lines = append(ser.Lines, pd.Lines...)
	ser.Routes = append(ser.Routes, pd.Routes...)
	ser.ScheduledStopPoints = append(ser.ScheduledStopPoints, pd.ScheduledStopPoints...)
	ser.ServiceLinks = append(ser.ServiceLinks, pd.ServiceLinks...)
	ser.StopAssignments = append(ser.StopAssignments, pd.StopAssignments...)
	ret.Frames.Frames = append(ret.Frames.Frames, ser)

	cal := ServiceCalendarFrame{}
	cal.Id = CreateFrameId("ServiceCalendarFrame_EU_PI_CALENDAR", "ita")
	cal.Version = "1"
	cal.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_CALENDAR")
	cal.ServiceCalendar = append(cal.ServiceCalendar, pd.ServiceCalendars...)
	ret.Frames.Frames = append(ret.Frames.Frames, cal)

	tim := TimetableFrame{}
	tim.Id = CreateFrameId("TimetableFrame_EU_PI_TIMETABLE", "ita")
	tim.Version = "1"
	tim.TypeOfFrameRef = MkTypeOfFrameRef("EU_PI_TIMETABLE")
	tim.VehicleJourneys = append(tim.VehicleJourneys, pd.VehicleJourneys...)
	ret.Frames.Frames = append(ret.Frames.Frames, tim)

	return ret
}
