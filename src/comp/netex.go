// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package comp

import (
	"regexp"
	"time"

	"github.com/noi-techpark/go-netex"
)

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

func MkRef(tp string, id string) netex.Ref {
	r := netex.Ref{}
	r.Ref = id
	r.Version = "1"
	r.XMLName.Local = tp + "Ref"
	return r
}

func MkTypeOfFrameRef(tp string) netex.Ref {
	r := netex.Ref{}
	r.Ref = "epip:" + tp
	r.Version = "1"
	r.XMLName.Local = "TypeOfFrameRef"
	return r
}

func ValidAYear() netex.ValidBetween {
	v := netex.ValidBetween{}
	v.FromDate = time.Now().Truncate(time.Hour * 24)
	v.ToDate = time.Now().AddDate(1, 0, 0).Truncate(time.Hour * 24)
	return v
}

func AppendSafe[T any](h *[]T, t ...T) *[]T {
	if len(t) > 0 {
		if h == nil {
			h = &t
		} else {
			*h = append(*h, t...)
		}
	}
	return h
}

func DefaultCompositFrame() netex.CompositeFrame {
	c := netex.CompositeFrame{}
	c.Version = "1"
	c.ValidBetween = ValidAYear()
	c.Codespaces.Codespace.Id = "ita"
	c.Codespaces.Codespace.Xmlns = "ita"
	c.Codespaces.Codespace.XmlnsUrl = "http://www.ita.it"
	c.Codespaces.Codespace.Description = "Italian Profile"
	c.FrameDefaults.DefaultCodespaceRef.Ref = "ita"
	return c
}
