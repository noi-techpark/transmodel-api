// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package provider

import (
	"opendatahub/sta-nap-export/config"
	"opendatahub/sta-nap-export/netex"
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

	var c52 *netex.Vehicle
	for _, v := range nt.Vehicles {
		if v.Id == "IT:ITH10:Vehicle:BIKE_SHARING_BOLZANO:City_52M" {
			c52 = &v
			break
		}
	}
	assert.Assert(t, c52 != nil)
	assert.Equal(t, c52.Name, "Sunrise")
	assert.Equal(t, c52.ShortName, "Sunrise")
	assert.Equal(t, c52.RegistrationNumber, "")
	assert.Equal(t, c52.OperationalNumber, "")
	assert.Equal(t, c52.PrivateCode, "City 52M")

	assert.Equal(t, c52.OperatorRef.Ref, "IT:ITH10:Operator:Municipality_of_Bolzano_bikesharing")
	o := netex.GetOperator(&config.Cfg, b.origin)
	assert.Equal(t, o.Id, c52.OperatorRef.Ref)
	assert.Equal(t, c52.VehicleTypeRef.Ref, "IT:ITH10:CycleModelProfile:BIKE_SHARING_BOLZANO:muscular")
}
