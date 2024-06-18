// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

type FMData struct {
	Conditions []FacilityCondition
}
type FMProvider interface {
	SiriFM() (FMData, error)
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

func FM(ps []FMProvider) (Siri, error) {
	siri := NewSiri()

	for _, p := range ps {
		dt, err := p.SiriFM()
		if err != nil {
			return siri, err
		}
		siri.AppencFcs(dt.Conditions)
	}

	return siri, nil
}
