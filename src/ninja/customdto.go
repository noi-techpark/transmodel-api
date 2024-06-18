// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>

// SPDX-License-Identifier: AGPL-3.0-or-later

package ninja

type OdhStation[T any] struct {
	Scode   string
	Sname   string
	Sorigin string
	Scoord  OdhCoord `json:"scoordinate"`
	Smeta   T        `json:"smetadata"`
}

type OdhCoord struct {
	X    float32
	Y    float32
	Srid uint32
}

type MetaAny map[string]any

type OdhLatest struct {
	MPeriod    int       `json:"mperiod"`
	MValidTime NinjaTime `json:"mvalidtime"`
	MValue     int       `json:"mvalue"`
	Scode      string    `json:"scode"`
	Stype      string    `json:"stype"`
	Tname      string    `json:"tname"`
}
