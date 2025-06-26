// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package comp

import "github.com/noi-techpark/go-netex"

type OdhEcharging struct {
	Scode       string
	Sname       string
	Sorigin     string
	Scoordinate struct {
		X    float32
		Y    float32
		Srid uint32
	}
	Smetadata struct {
		State    string
		Capacity int32
	}
}

func compFrame(serviceName string, pd StParkingData) netex.CompositeFrame {
	ret := DefaultCompositFrame()
	ret.Id = CreateFrameId(netex.TypeCompositeFrameStopOffer, serviceName)
	ret.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeStopOffer)

	site := siteFrame(serviceName)
	ret.Frames.Frames = append(ret.Frames.Frames, &site)

	res := netex.ResourceFrame{}
	res.Id = CreateFrameId(netex.TypeResourceFrameCommon, serviceName)
	res.Version = "1"
	res.TypeOfFrameRef = MkTypeOfFrameRef(netex.EpipTypeCommon)
	ret.Frames.Frames = append(ret.Frames.Frames, &res)

	site.Parkings = &pd.Parkings
	res.Operators = &pd.Operators

	return ret
}

type StParkingData struct {
	Parkings  []netex.Parking
	Operators []netex.Operator
}

type StParking interface {
	StParking() (StParkingData, error)
}

func GetParking(ps []StParking) ([]netex.CompositeFrame, error) {
	ret := []netex.CompositeFrame{}

	apd := StParkingData{}

	for _, p := range ps {
		pd, err := p.StParking()
		if err != nil {
			return ret, err
		}
		apd.Parkings = append(apd.Parkings, pd.Parkings...)
		apd.Operators = append(apd.Operators, pd.Operators...)
	}

	ret = append(ret, compFrame("PARKING", apd))

	return ret, nil
}
