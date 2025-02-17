// SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: CC0-1.0

package provider

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestUnmarshalFlight(t *testing.T) {
	fs := NewFlightsSkyalps()
	f, err := os.ReadFile("./testdata/skyalps.xml")
	assert.NilError(t, err)
	dt, err := fs.fromNetex(&f)
	assert.NilError(t, err, "error getting the flights")

	assert.Equal(t, len(dt.Operators), 1, "wrong number of operators")
	assert.Equal(t, len(dt.JourneyPatterns), 52, "wrong number of journey patterns")
	assert.Equal(t, len(dt.StopPlaces), 28, "wrong number of stop places")
	assert.Equal(t, len(dt.Lines), 52, "wrong number of lines")
	assert.Equal(t, len(dt.Routes), 52, "wrong number of routes")
	assert.Equal(t, len(dt.ScheduledStopPoints), 28, "wrong number of stop points")
	assert.Equal(t, len(dt.ServiceCalendars), 1, "wrong number of service calendars")
	assert.Equal(t, len(dt.ServiceLinks), 52, "wrong number of service links")
	assert.Equal(t, len(dt.StopAssignments), 28, "wrong number of stop assignments")
	assert.Equal(t, len(dt.VehicleJourneys), 1814, "wrong number of vehicle journeys")

	// unmarshalling of gml with namespaces and such
	assert.Assert(t, dt.ServiceLinks[0].LineString.Id != "")
	assert.Assert(t, dt.ServiceLinks[0].LineString.PosList != "")
}
