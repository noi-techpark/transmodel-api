// SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestReadOps(t *testing.T) {
	ops := readOps(getCsvPath())
	if len(ops) == 0 {
		t.Log("Operators loaded empty")
		t.Fail()
	}
}

func TestOpsContent(t *testing.T) {
	ops := readOps(getCsvPath())
	mapped := mapByOrigin(ops)

	bsb := mapped["BIKE_SHARING_BOLZANO"]
	assert.Equal(t, "urp@comune.bolzano.it", bsb.Email)
	assert.Equal(t, "0471997111", bsb.Phone)

	hal := mapped["HAL-API"]
	assert.Equal(t, "https://www.carsharing.bz.it/it/", hal.Url)
	assert.Equal(t, "Via Beda Weber 1", hal.Street)
}
