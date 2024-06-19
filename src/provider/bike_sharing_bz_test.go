// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package provider

import (
	"opendatahub/sta-nap-export/config"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"
	"slices"
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

	// Cardinalities
	assert.Equal(t, len(nt.Vehicles), len(cycles.Data))
	assert.Equal(t, len(nt.Fleets), 1)
	assert.Equal(t, len(nt.Fleets[0].Members), len(cycles.Data))
	assert.Equal(t, len(nt.Parkings), len(sharings.Data))
	assert.Equal(t, len(nt.Operators), 1)
	assert.Equal(t, len(nt.Constraints), 1)
	assert.Equal(t, len(nt.CycleModels), 2)
	assert.Equal(t, len(nt.Modes), 1)
	assert.Equal(t, len(nt.Modes[0].Submodes), 1)
	assert.Equal(t, len(nt.Services), 1)
	assert.Equal(t, len(nt.Services[0].Fleets), 1)

	var c52 *netex.Vehicle
	for _, v := range nt.Vehicles {
		if v.Id == "IT:ITH10:Vehicle:BIKE_SHARING_BOLZANO:City_52M" {
			c52 = &v
			break
		}
	}
	// Some basic field mapping
	assert.Assert(t, c52 != nil)
	assert.Equal(t, c52.Name, "Sunrise")
	assert.Equal(t, c52.ShortName, "Sunrise")
	assert.Equal(t, c52.RegistrationNumber, "")
	assert.Equal(t, c52.OperationalNumber, "")
	assert.Equal(t, c52.PrivateCode, "City 52M")

	// Operator correctly included
	// Referential integrity should already be checked by the validation
	assert.Equal(t, c52.OperatorRef.Ref, "IT:ITH10:Operator:Municipality_of_Bolzano_bikesharing")
	o := netex.GetOperator(&config.Cfg, b.origin)
	assert.Equal(t, o.Id, c52.OperatorRef.Ref)
	assert.DeepEqual(t, o, nt.Operators[0])

	assert.Equal(t, c52.VehicleTypeRef.Ref, "IT:ITH10:CycleModelProfile:BIKE_SHARING_BOLZANO:muscular")

	mus := nt.CycleModels[slices.IndexFunc(nt.CycleModels, func(m netex.CycleModelProfile) bool { return m.Id == c52.VehicleTypeRef.Ref })]
	assert.Equal(t, mus.Basket, true)
	assert.Equal(t, mus.Battery, false)
	assert.Equal(t, mus.ChildSeat, "none")
	assert.Equal(t, mus.Lamps, true)
	assert.Equal(t, mus.Lock, false)
	assert.Equal(t, mus.Pump, false)

	park := nt.Parkings[slices.IndexFunc(nt.Parkings, func(m netex.Parking) bool { return m.ShortName == "Viale Europa" })]
	assert.Equal(t, park.TotalCapacity, int32(12))
}
