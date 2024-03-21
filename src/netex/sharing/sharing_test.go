// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package sharing

import (
	"opendatahub/sta-nap-export/netex"
	"testing"

	"gotest.tools/v3/assert"
)

func TestOdhGetSharing(t *testing.T) {
	netex.TestOdhGet(t, bzSharing)
	netex.TestOdhGet(t, bzBike)
	netex.TestOdhGet(t, meBike)
	netex.TestOdhGet(t, papingSharing)
	netex.TestOdhGet(t, halSharing)
}

func meBike() ([]OdhMobility[metaAny], error) {
	return odhMob[[]OdhMobility[metaAny]]("Bicycle", "BIKE_SHARING_MERANO")
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
	_, err := frame([]SharingProvider{&Bz{}})
	assert.NilError(t, err)
}
