// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"encoding/xml"
	"opendatahub/sta-nap-export/ninja"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetOdhData(t *testing.T) {
	ninja.BaseUrl = "https://mobility.api.opendatahub.testingmachine.eu"
	ninja.Referer = "sta-nap-export-unit-test"

	res, err := getOdhParking()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(res)

	// Looks like we managed to parse the response structure
}

func bzCentro() odhParking {
	var o odhParking
	o.Scode = "bolzano-centro"
	o.Sname = "Bolzano Centro"
	o.Sorigin = "skidata"
	o.Scoord.X = 11.123123
	o.Scoord.Y = 46.432132
	o.Smeta.Capacity = 10
	o.Smeta.StandardName = "Bolzano Centro 1"
	o.Smeta.ParkingType = "parkingZone"
	o.Smeta.ParkingVehicleTypes = "all"
	o.Smeta.ParkingLayout = "covered"
	o.Smeta.ParkingProhibitions = true
	o.Smeta.ParkingCharging = true
	o.Smeta.ParkingSurveillance = false
	o.Smeta.ParkingReservation = "noReservations"
	return o
}

func TestMapNetex(t *testing.T) {
	o := bzCentro()

	p := mapToNetex([]odhParking{o})
	ps := NetexParkings{Parkings: p}

	x, err := xml.MarshalIndent(ps, "", " ")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log(string(x))
}

func TestParkingId(t *testing.T) {
	s1 := "merano0123"
	s2 := "me:test123"
	s3 := ":',|?=+"

	assert.Equal(t, ParkingId(s1), "IT:ITH10:Parking:merano0123")
	assert.Equal(t, ParkingId(s2), "IT:ITH10:Parking:me_test123")
	assert.Equal(t, ParkingId(s3), "IT:ITH10:Parking:_______")
}
