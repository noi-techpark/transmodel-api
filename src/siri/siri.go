// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package siri

type RtFMData struct {
	Conditions []FacilityCondition
}
type RtFMProvider interface {
	RtSharing() (RtFMData, error)
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

func FM(ps []RtFMProvider) (Siri, error) {
	siri := newSiri()

	for _, p := range ps {
		dt, err := p.RtSharing()
		if err != nil {
			return siri, err
		}
		siri.appencFcs(dt.Conditions)
	}

	return siri, nil
}
