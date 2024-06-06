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
	o.Scoordinate.X = 11.123123
	o.Scoordinate.Y = 46.432132
	o.Smetadata.Capacity = 10
	o.Smetadata.StandardName = "Bolzano Centro 1"
	o.Smetadata.Netex.Type = "parkingZone"
	o.Smetadata.Netex.VehicleTypes = "all"
	o.Smetadata.Netex.Layout = "covered"
	o.Smetadata.Netex.HazardProhibited = true
	o.Smetadata.Netex.Charging = true
	o.Smetadata.Netex.Surveillance = false
	o.Smetadata.Netex.Reservation = "noReservations"
	return o
}

func TestMapNetex(t *testing.T) {
	netex.InitConfig()
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
