// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>

// SPDX-License-Identifier: AGPL-3.0-or-later

package ninja

import (
	"testing"

	"gotest.tools/v3/assert"
)

type BikeShareMeta struct {
	Address   string
	TotalBays int `json:"total-bays"`
}

func TestLoadJson(t *testing.T) {
	j, err := LoadJsonFile[[]OdhStation[BikeShareMeta]]("test/ninja_loadjson.json")
	assert.NilError(t, err, "Failed to load JSON")
	assert.Equal(t, j.Data[0].Sname, "Viale della Stazione - Bahnhofsallee", "Unexpected mapping from JSON")
	assert.Equal(t, j.Data[0].Smeta.TotalBays, 12, "Unexpected mapping from JSON")
}

func TestReqHook1(t *testing.T) {
	j, err := LoadJsonFile[[]OdhStation[BikeShareMeta]]("test/ninja_loadjson.json")
	assert.NilError(t, err, "Failed to load JSON")

	req := NinjaRequest{}
	req.Origin = "test"
	TestReqHook = func(nr *NinjaRequest) (any, error) {
		assert.Equal(t, nr.Origin, req.Origin, "Passed request not matching the one in hook")
		return j, nil
	}

	res := NinjaResponse[[]OdhStation[BikeShareMeta]]{}
	err = StationType(&req, &res)
	assert.NilError(t, err, "Error calling ninja with req hook")

	assert.Assert(t, res.Data[0].Scode != "", "zero value in returned data")
	assert.Equal(t, j.Data[0].Scode, res.Data[0].Scode, "Mismatch between returned value and hook")
}
