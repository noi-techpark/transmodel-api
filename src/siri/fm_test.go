// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

import (
	"opendatahub/sta-nap-export/netex"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetOdhParkingState(t *testing.T) {
	netex.NinjaTestSetup()
	res, err := odhParkingState("115")
	assert.NilError(t, err)
	t.Log(res)
}

func TestResponseId(t *testing.T) {
	i := responseId.next()
	i2 := responseId.next()
	assert.Equal(t, i+1, i2)
}
