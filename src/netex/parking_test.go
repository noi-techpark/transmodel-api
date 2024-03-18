// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"opendatahub/sta-nap-export/ninja"
	"testing"
)

func TestGetOdhData(t *testing.T) {
	ninja.BaseUrl = "https://mobility.api.opendatahub.testingmachine.eu"
	ninja.Referer = "sta-nap-export-unit-test"

	res, err := getOdhParking()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(res)

	// Looks like we managed to parse the response structure
}
