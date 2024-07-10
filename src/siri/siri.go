// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

import (
	"strconv"
)

type FMData struct {
	Conditions []FacilityCondition
}
type FMProvider interface {
	SiriFM(query Query) (FMData, error)
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

func (q Query) float(k string) float64 {
	i, _ := strconv.ParseFloat(q[k][0], 64)
	return i
}
func (q Query) int(k string) int {
	i, _ := strconv.Atoi(q[k][0])
	return i
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

func FM(ps []FMProvider, query Query) (Siri, error) {
	siri := NewSiri()
	maxSize := query.MaxSize()
	cnt := 0

	for _, p := range ps {
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
