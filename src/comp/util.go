// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package comp

import (
	"time"

	"github.com/noi-techpark/go-netex"
)

func CreateID(segments ...string) string {
	return netex.NewId(append([]string{"IT", "ITH10"}, segments...)...)
}

func CreateFrameId(segments ...string) string {
	return netex.NewFrameId(segments...)
}

func MkRef(tp string, id string) netex.Ref {
	return netex.NewRef(tp, id, "1")
}

func MkTypeOfFrameRef(tp string) netex.TypeOfFrameRef {
	return netex.NewTypeOfFrameRef(tp, "1")
}

func ValidAYear() netex.ValidBetween {
	v := netex.ValidBetween{}
	v.FromDate = time.Now().Truncate(time.Hour * 24)
	v.ToDate = time.Now().AddDate(1, 0, 0).Truncate(time.Hour * 24)
	return v
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

func siteFrame(serviceName string) netex.SiteFrame {
	var site netex.SiteFrame
	site.Id = CreateFrameId(netex.TypeSiteFrameStop, serviceName)
	site.Version = "1"
	site.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeStop)
	return site
}
