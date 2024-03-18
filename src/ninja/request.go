// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>

// SPDX-License-Identifier: AGPL-3.0-or-later

package ninja

import "time"

type Repr string

const (
	FlatNode  Repr = "flat,node"
	TreeNode  Repr = "tree,node"
	FlatEdge  Repr = "flat,edge"
	TreeEdge  Repr = "tree,edge"
	FlatEvent Repr = "flat,event"
	TreeEvent Repr = "tree,event"
)

type NinjaRequest struct {
	Repr     Repr
	Origin   string
	Limit    int64
	Offset   uint64
	Select   string
	Where    string
	Shownull bool
	Distinct bool
	Timezone string

	EventOrigins []string
	EdgeTypes    []string

	Timepoint time.Time

	StationTypes []string
	DataTypes    []string

	From time.Time
	To   time.Time
}

// Defaults according to Ninja Swagger documentation
func DefaultNinjaRequest() *NinjaRequest {
	def := new(NinjaRequest)
	def.Repr = FlatNode
	def.Limit = 200
	def.Offset = 0
	def.Shownull = false
	def.Distinct = true
	return def
}

func (nr *NinjaRequest) AddStationType(stationType string) {
	nr.StationTypes = append(nr.StationTypes, stationType)
}
func (nr *NinjaRequest) AddDataType(dataType string) {
	nr.DataTypes = append(nr.DataTypes, dataType)
}
func (nr *NinjaRequest) AddEdgeType(edgeType string) {
	nr.EdgeTypes = append(nr.EdgeTypes, edgeType)
}
func (nr *NinjaRequest) AddEventOrigin(eventOrigin string) {
	nr.EventOrigins = append(nr.EventOrigins, eventOrigin)
}
