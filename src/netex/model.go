// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later
package netex

import (
	"encoding/xml"
	"time"
)

type CompositeFrame struct {
	XMLName        xml.Name `xml:"CompositeFrame"`
	Id             string   `xml:"id,attr"`
	Version        string   `xml:"version,attr"`
	ValidBetween   ValidBetween
	TypeOfFrameRef TypeOfFrameRef
	Codespaces     struct {
		Codespace struct {
			Id          string `xml:"id,attr"`
			Xmlns       string
			XmlnsUrl    string
			Description string
		}
	} `xml:"codespaces"`
	FrameDefaults struct {
		DefaultCodespaceRef Ref
	}
	Frames struct{ Frames []any } `xml:"frames"`
}

type SiteFrame struct {
	XMLName        xml.Name `xml:"SiteFrame"`
	Id             string   `xml:"id,attr"`
	Version        string   `xml:"version,attr"`
	TypeOfFrameRef TypeOfFrameRef
	Parkings       any `xml:"parkings,omitempty"`
}

func (c *CompositeFrame) Defaults() {
	c.Version = "1"
	c.ValidBetween.AYear()
	c.Codespaces.Codespace.Id = "ita"
	c.Codespaces.Codespace.Xmlns = "ita"
	c.Codespaces.Codespace.XmlnsUrl = "http://www.ita.it"
	c.Codespaces.Codespace.Description = "Italian Profile"
	c.FrameDefaults.DefaultCodespaceRef.Ref = "ita"
}

type Ref struct {
	XMLName xml.Name
	Ref     string `xml:"ref,attr"`
	Version string `xml:"version,attr,omitempty"`
}
type TypeOfFrameRef struct {
	XMLName xml.Name
	Ref     string `xml:"ref,attr"`
	Version string `xml:"versionRef,attr,omitempty"`
}

type ValidBetween struct {
	FromDate time.Time
	ToDate   time.Time
}
