// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

import (
	"testing"

	"gotest.tools/v3/assert"
)

type prov struct{}

func (prov) SiriFM(q Query) (FMData, error) {
	maxSize := q.MaxSize()
	ret := FMData{}
	ret.Conditions = make([]FacilityCondition, maxSize)
	return ret, nil
}

func (prov) MatchOperator(id string) bool {
	return true
}

func TestMaxSize(t *testing.T) {
	q := Query{}
	q["maxSize"] = []string{"15"}
	si, err := FM([]FMProvider{prov{}}, q)
	assert.NilError(t, err)
	assert.Equal(t, len(si.ServiceDelivery.FacilityMonitoringDelivery.FacilityCondition), 15)
}
