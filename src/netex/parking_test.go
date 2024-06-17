// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

// import (
// 	"encoding/xml"
// 	"opendatahub/sta-nap-export/config"
// 	"testing"

// 	"gotest.tools/v3/assert"
// )

// func TestGetOdhData(t *testing.T) {
// 	TestOdhGet(t, getOdhParking)
// 	TestOdhGet(t, getOdhEcharging)
// }

// func TestCallFull(t *testing.T) {
// 	config.InitConfig()
// 	NinjaTestSetup()

// 	f, err := GetParking()
// 	if err != nil {
// 		t.Log(err)
// 		t.Fail()
// 	}
// 	t.Log(f)
// }

// func bzCentro() OdhParking {
// 	var o OdhParking
// 	o.Scode = "bolzano-centro"
// 	o.Sname = "Bolzano Centro"
// 	o.Sorigin = "skidata"
// 	o.Scoordinate.X = 11.123123
// 	o.Scoordinate.Y = 46.432132
// 	o.Smetadata.Capacity = 10
// 	o.Smetadata.StandardName = "Bolzano Centro 1"
// 	o.Smetadata.Netex.Type = "parkingZone"
// 	o.Smetadata.Netex.VehicleTypes = "all"
// 	o.Smetadata.Netex.Layout = "covered"
// 	o.Smetadata.Netex.HazardProhibited = true
// 	o.Smetadata.Netex.Charging = true
// 	o.Smetadata.Netex.Surveillance = false
// 	o.Smetadata.Netex.Reservation = "noReservations"
// 	return o
// }
// func echargingPk() OdhEcharging {
// 	var o OdhEcharging
// 	o.Scode = "bolzano-centro"
// 	o.Sname = "Bolzano Centro"
// 	o.Sorigin = "skidata"
// 	o.Scoordinate.X = 11.123123
// 	o.Scoordinate.Y = 46.432132
// 	o.Smetadata.Capacity = 10
// 	return o
// }

// func TestMapParking(t *testing.T) {
// 	config.InitConfig()
// 	NinjaTestSetup()
// 	odh := bzCentro()

// 	ps, os := mapParking([]OdhParking{odh})

// 	marshall := func(a any) string {
// 		r, err := xml.MarshalIndent(a, "", " ")
// 		if err != nil {
// 			t.Log(err)
// 			t.Fail()
// 		}
// 		return string(r)
// 	}

// 	s := marshall(ps)
// 	t.Log(s)

// 	s = marshall(os)
// 	t.Log(s)

// }
// func TestMapEcharging(t *testing.T) {
// 	config.InitConfig()
// 	odh := echargingPk()

// 	ps, os := mapEcharging([]OdhEcharging{odh})

// 	marshall := func(a any) string {
// 		r, err := xml.MarshalIndent(a, "", " ")
// 		if err != nil {
// 			t.Log(err)
// 			t.Fail()
// 		}
// 		return string(r)
// 	}

// 	s := marshall(ps)
// 	t.Log(s)

// 	s = marshall(os)
// 	t.Log(s)

// }

// func TestParkingId(t *testing.T) {
// 	s1 := "merano0123"
// 	s2 := "me:test123"
// 	s3 := ":',|?=+"

// 	assert.Equal(t, CreateID("Parking", s1), "IT:ITH10:Parking:merano0123")
// 	assert.Equal(t, CreateID("Parking", s2), "IT:ITH10:Parking:me_test123")
// 	assert.Equal(t, CreateID("Parking", s3), "IT:ITH10:Parking:_______")
// }
