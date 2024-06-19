// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package provider

import (
	"opendatahub/sta-nap-export/config"
	"opendatahub/sta-nap-export/ninja"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMapping(t *testing.T) {
	sharings, err := ninja.LoadJsonFile[bikeBzSharing]("test/bike_sharing_bz_sharing.json")
	assert.NilError(t, err, "Failed to load JSON")

	cycles, err := ninja.LoadJsonFile[bikeBzCycles]("test/bike_sharing_bz_cycles.json")
	assert.NilError(t, err, "Failed to load JSON")

	ninja.TestReqHook = func(nr *ninja.NinjaRequest) (any, error) {
		if nr.StationTypes[0] == "BikesharingStation" && len(nr.StationTypes) == 1 {
			return sharings, nil
		}
		if nr.StationTypes[0] == "Bicycle" && len(nr.StationTypes) == 1 {
			return cycles, nil
		}
		t.Error("Hmmmm. Should not get to here.", nr)
		return nil, nil
	}

	b := NewBikeBz()
	config.InitConfig()

	nt, err := b.StSharing()
	assert.NilError(t, err)
	assert.Equal(t, len(nt.Vehicles), len(cycles.Data))
	assert.Equal(t, len(nt.Parkings), len(sharings.Data))

}
