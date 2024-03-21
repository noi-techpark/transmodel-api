// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"opendatahub/sta-nap-export/ninja"
	"regexp"
	"testing"
)

func TestOdhGet[T any](t *testing.T, f func() (T, error)) {
	ninja.BaseUrl = "https://mobility.api.opendatahub.com"
	ninja.Referer = "sta-nap-export-unit-test"

	res, err := f()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(res)
}

// As per NeTEx spec, IDs must only contain non-accented charaters, numbers, hyphens and underscores
var idInvalid = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func CreateID(segments ...string) string {
	id := "IT:ITH10"
	for _, s := range segments {
		id += (":" + idInvalid.ReplaceAllString(s, "_"))
	}
	return id
}
