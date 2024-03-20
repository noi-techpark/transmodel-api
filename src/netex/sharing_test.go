// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"opendatahub/sta-nap-export/ninja"
	"testing"
)

func testOdhGet[T any](t *testing.T, f func() (T, error)) {
	ninja.BaseUrl = "https://mobility.api.opendatahub.com"
	ninja.Referer = "sta-nap-export-unit-test"

	res, err := f()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(res)
}

func TestOdhGetSharing(t *testing.T) {
	testOdhGet(t, bikeSharingBz)
	testOdhGet(t, bikeMe)
	testOdhGet(t, bikeSharingPapin)
	testOdhGet(t, carSharingHal)
}
