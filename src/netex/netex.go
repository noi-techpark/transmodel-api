// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"opendatahub/sta-nap-export/ninja"
	"regexp"
	"testing"
	"time"
)

func NinjaTestSetup() {
	ninja.BaseUrl = "https://mobility.api.opendatahub.com"
	ninja.Referer = "sta-nap-export-unit-test"
}

func TestOdhGet[T any](t *testing.T, f func() (T, error)) {
	NinjaTestSetup()

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

func CreateFrameId(segments ...string) string {
	return "edp:" + CreateID(segments...)
}

func MkRef(tp string, segments ...string) Ref {
	r := Ref{}
	r.Ref = CreateID(append([]string{tp}, segments...)...)
	r.Version = "1"
	r.XMLName.Local = tp + "Ref"
	return r
}
func MkTypeOfFrameRef(tp string) TypeOfFrameRef {
	r := TypeOfFrameRef{}
	r.Ref = "epip:" + tp
	r.Version = "1"
	r.XMLName.Local = "TypeOfFrameRef"
	return r
}

func (v *ValidBetween) AYear() {
	v.FromDate = time.Now().Truncate(time.Hour * 24)
	v.ToDate = time.Now().AddDate(1, 0, 0).Truncate(time.Hour * 24)
}

func AppendSafe[T any](h *[]T, t ...T) *[]T {
	if len(t) > 0 {
		if h == nil {
			h = &t
		}
		*h = append(*h, t...)
	}
	return h
}
