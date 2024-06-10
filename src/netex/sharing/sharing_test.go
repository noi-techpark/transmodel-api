// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package sharing

import (
	"encoding/xml"
	"opendatahub/sta-nap-export/netex"
	"testing"

	"gotest.tools/v3/assert"
)

func TestOdhGetSharing(t *testing.T) {
	netex.TestOdhGet(t, bzBike)
	// netex.TestOdhGet(t, meBike)
	netex.TestOdhGet(t, papingSharing)
	netex.TestOdhGet(t, halSharing)
}

func papingSharing() ([]OdhMobility[metaAny], error) {
	return odhMob[[]OdhMobility[metaAny]]("BikesharingStation", "BIKE_SHARING_PAPIN")
}
func halSharing() ([]OdhMobility[metaAny], error) {
	return odhMob[[]OdhMobility[metaAny]]("CarsharingStation", "HAL-API")
}

func TestEmptyProvider(t *testing.T) {
	_, err := frame([]SharingProvider{})
	assert.NilError(t, err)
}

func TestBzProvider(t *testing.T) {
	netex.NinjaTestSetup()
	netex.InitConfig()
	_, err := frame([]SharingProvider{&BikeBz{}})
	assert.NilError(t, err)
}

func TestGetSharingData(t *testing.T) {
	netex.NinjaTestSetup()

	d, err := GetSharing()
	assert.NilError(t, err)

	x, err := xml.MarshalIndent(d, "", " ")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log(string(x))
}
