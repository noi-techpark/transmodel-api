// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>

// SPDX-License-Identifier: AGPL-3.0-or-later

package ninja

type OdhStation[T any] struct {
	Scode   string
	Sname   string
	Sorigin string
	Scoord  struct {
		X    float32
		Y    float32
		Srid uint32
	} `json:"scoordinate"`
	Smeta T `json:"smetadata"`
}

type MetaAny map[string]any
