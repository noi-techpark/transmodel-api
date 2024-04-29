// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package parking

import (
	"encoding/xml"
	"opendatahub/sta-nap-export/netex"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetOdhData(t *testing.T) {
	netex.TestOdhGet(t, getOdhParking)
}

func bzCentro() OdhParking {
	var o OdhParking
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
	odh := bzCentro()

	ps, os := mapToNetex([]OdhParking{odh})

	marshall := func(a any) string {
		r, err := xml.MarshalIndent(a, "", " ")
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		return string(r)
	}

	s := marshall(ps)
	t.Log(s)

	s = marshall(os)
	t.Log(s)

}

func TestParkingId(t *testing.T) {
	s1 := "merano0123"
	s2 := "me:test123"
	s3 := ":',|?=+"

	assert.Equal(t, netex.CreateID("Parking", s1), "IT:ITH10:Parking:merano0123")
	assert.Equal(t, netex.CreateID("Parking", s2), "IT:ITH10:Parking:me_test123")
	assert.Equal(t, netex.CreateID("Parking", s3), "IT:ITH10:Parking:_______")
}
