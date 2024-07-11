// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

import (
	"slices"
	"strconv"
	"strings"
)

type FMData struct {
	Conditions []FacilityCondition
}
type FMProvider interface {
	SiriFM(query Query) (FMData, error)
	MatchOperator(id string) bool
}

func MapFacilityStatus(available int, partialThreshold int) string {
	switch {
	case available == 0:
		return "notAvailable"
	case available <= partialThreshold:
		return "partiallyAvailable"
	default:
		return "available"
	}
}

type Query map[string][]string

func WrapQuery(src map[string][]string) Query {
	q := Query{}
	for k, v := range src {
		q[k] = v
	}
	return q
}

func (q Query) int(k string) int {
	s, found := q[k]
	if found && len(s) > 0 {
		i, _ := strconv.Atoi(s[0])
		return i
	} else {
		return 0
	}
}
func (q Query) MaxSize() int {
	m := q.int("maxSize")
	if m == 0 {
		m = -1
	}
	return m
}
func (q Query) DatasetIds() []string {
	return q["datasetId"]
}
func (q Query) FacilityRef() []string {
	return q["facilityRef"]
}
func (q Query) Operators() []string {
	ret := []string{}
	for _, p := range q["operators"] {
		ss := strings.Split(p, ",")
		ret = append(ret, ss...)
	}
	return ret
}

func (q Query) float(k string) float64 {
	s, found := q[k]
	if found && len(s) > 0 {
		i, _ := strconv.ParseFloat(s[0], 64)
		return i
	} else {
		return 0
	}
}

func (q Query) Lat() float64 {
	return q.float("lat")
}
func (q Query) Lon() float64 {
	return q.float("lon")
}
func (q Query) MaxDistance() float64 {
	return q.float("maxDistance")
}

func matchAnyOp(prov FMProvider, ids []string) bool {
	return slices.ContainsFunc(ids, func(id string) bool {
		return prov.MatchOperator(id)
	})
}

func FM(ps []FMProvider, query Query) (Siri, error) {
	siri := NewSiri()
	maxSize := query.MaxSize()
	cnt := 0

	for _, p := range ps {
		os := query.Operators()
		if len(os) > 0 && !matchAnyOp(p, os) {
			continue
		}
		dt, err := p.SiriFM(query)
		if err != nil {
			return siri, err
		}
		if maxSize > 0 {
			query["maxSize"] = []string{strconv.Itoa(maxSize - cnt)}
			current := len(dt.Conditions)
			cnt += current
			if cnt >= maxSize {
				siri.AppencFcs(dt.Conditions[:current-(cnt-maxSize)])
				break
			}
		}
		siri.AppencFcs(dt.Conditions)
	}

	return siri, nil
}
