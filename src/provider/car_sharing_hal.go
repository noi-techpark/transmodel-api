// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package provider

import (
	"fmt"
	"log/slog"
	"opendatahub/sta-nap-export/config"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/ninja"
	"opendatahub/sta-nap-export/siri"

	"golang.org/x/exp/maps"
)

type odhHALCar []struct {
	ninja.OdhStation[halCarMeta]
	Pmetadata struct {
		Company struct {
			Uid       string
			ShortName string
			FullName  string
		}
	}
}

type CarHAL struct {
	cars     odhHALCar
	origin   string
	provider string
}

type halCarMeta struct {
	Brand        string
	Model        string
	LicensePlate string
	Features     *struct {
		Doors           uint8
		Seats           uint8
		Chains          bool
		Satnav          bool
		Skirack         bool
		Roofrack        bool
		Childseat       string
		Cyclerack       bool
		Wintertyres     bool
		Transmission    string
		Cruisecontrol   bool
		Trailerhitch    bool
		Usbpowersockets bool
	}
}

type halSharingMeta struct {
	Company struct {
		Uid       string
		ShortName string
		FullName  string
	}
	Bookahead         bool
	FixedParking      bool
	Spontaneously     bool
	AvailableVehicles int
}

const ORIGIN_CAR_SHARING_HAL_API = "HAL-API"

func NewCarSharingHal() *CarHAL {
	b := CarHAL{}
	b.origin = ORIGIN_CAR_SHARING_HAL_API
	return &b
}

func (b *CarHAL) GetOperator() netex.Operator {
	return netex.GetOperator(&config.Cfg, b.origin)
}
func (b *CarHAL) StSharing() (netex.StSharingData, error) {
	ret := netex.StSharingData{}
	if err := b.fetch(); err != nil {
		return ret, err
	}
	b.provider = b.cars[0].Pmetadata.Company.ShortName // As soon as there is more than one, we have to change some more stuff anyways

	// Operators
	o := b.GetOperator()
	ret.Operators = append(ret.Operators, o)

	// Modes of Operation
	m := netex.VehicleSharing{}
	m.Id = netex.CreateID("VehicleSharing", b.provider)
	m.Version = "1"
	sub := netex.Submode{}
	sub.Id = netex.CreateID("Submode", b.provider)
	sub.Version = "1"
	sub.TransportMode = "car"
	sub.SelfDriveSubmode = "hireCar"
	m.Submodes = append(m.Submodes, sub)
	ret.Modes = append(ret.Modes, m)

	models := make(map[string]netex.CarModelProfile)

	for _, c := range b.cars {
		modelname := c.Smeta.Brand
		p, found := models[modelname]
		if !found {
			// Car model profile
			p = netex.CarModelProfile{}
			p.Id = netex.CreateID("CarModelProfile", b.provider, modelname)
			p.Version = "1"
			p.ChildSeat = c.Smeta.Features.Childseat
			p.Seats = c.Smeta.Features.Seats
			p.Doors = c.Smeta.Features.Doors
			p.Transmission = c.Smeta.Features.Transmission
			p.CruiseControl = c.Smeta.Features.Cruisecontrol
			p.SatNav = c.Smeta.Features.Satnav
			p.AirConditioning = true
			p.Convertible = false
			p.UsbPowerSockets = c.Smeta.Features.Usbpowersockets
			p.WinterTyres = c.Smeta.Features.Wintertyres
			p.Chains = c.Smeta.Features.Chains
			p.TrailerHitch = c.Smeta.Features.Trailerhitch
			p.RoofRack = c.Smeta.Features.Roofrack
			p.CycleRack = c.Smeta.Features.Cyclerack
			p.SkiRack = c.Smeta.Features.Skirack
			models[modelname] = p
		}

		// Vehicles
		v := netex.Vehicle{}
		v.Id = netex.CreateID("Vehicle", b.provider, c.Scode)
		v.Version = "1"
		v.ValidBetween.AYear()
		v.Name = c.Sname
		v.ShortName = c.Sname
		v.PrivateCode = c.Scode
		v.RegistrationNumber = c.Smeta.LicensePlate
		v.OperatorRef = netex.MkRef("Operator", o.Id)
		v.VehicleTypeRef = netex.MkRef("CarModelProfile", p.Id)
		ret.Vehicles = append(ret.Vehicles, v)
	}
	ret.CarModels = maps.Values(models)

	// Fleets = all Vehicles + operator
	f := netex.Fleet{}
	f.Id = netex.CreateID("Fleet", b.provider)
	f.Version = "1"
	f.ValidBetween.AYear()
	for _, v := range ret.Vehicles {
		f.Members = append(f.Members, netex.MkRef("Vehicle", v.Id))
	}
	f.OperatorRef = netex.MkRef("Operator", o.Id)
	ret.Fleets = append(ret.Fleets, f)

	// Mobility services = Fleet + mode
	s := netex.VehicleSharingService{}
	s.Id = netex.CreateID("VehicleSharingService", b.provider)
	s.Version = "1"
	s.VehicleSharingRef = netex.MkRef("VehicleSharing", m.Id)
	s.FloatingVehicles = false
	for _, fl := range ret.Fleets {
		s.Fleets = append(s.Fleets, netex.MkRef("Fleet", fl.Id))
	}
	ret.Services = append(ret.Services, s)

	// Constraint zone
	c := netex.MobilityServiceConstraintZone{}
	c.Id = netex.CreateID("MobilityServiceConstraintZone", b.provider)
	c.Version = "1"
	c.GmlPolygon.Id = b.provider
	c.GmlPolygon.SetPoly(config.GML_PROVINCE_BZ)
	c.VehicleSharingRef = netex.MkRef("VehicleSharingService", s.Id)
	ret.Constraints = append(ret.Constraints, c)

	// Sharing as Parking (for SIRI reference)
	ss, err := FetchOdhStations[[]ninja.OdhStation[halSharingMeta]]("CarsharingStation", b.origin)
	if err != nil {
		return ret, err
	}
	for _, s := range ss {
		p := netex.Parking{}
		p.Id = netex.CreateID("Parking", s.Smeta.Company.ShortName, s.Scode)
		p.Version = "1"
		p.ShortName = s.Sname
		p.Centroid.Location.Longitude = s.Scoord.X
		p.Centroid.Location.Latitude = s.Scoord.Y
		p.OperatorRef = netex.MkRef("Operator", o.Id)
		p.GmlPolygon = nil
		p.Entrances = nil
		p.ParkingType = "rentalCarParking"
		p.ParkingVehicleTypes = "car"
		p.ParkingLayout = "undefined"
		p.ProhibitedForHazardousMaterials.Set(true)
		p.RechargingAvailable.Ignore()
		p.Secure.Ignore()
		p.ParkingReservation = "reservationRequired"
		p.ParkingProperties = nil

		p.Name = s.Sname
		p.PrincipalCapacity = int32(s.Smeta.AvailableVehicles)
		p.TotalCapacity = int32(s.Smeta.AvailableVehicles)
		ret.Parkings = append(ret.Parkings, p)
	}

	return ret, nil
}

func (b *CarHAL) fetch() error {
	cs, err := FetchOdhStations[odhHALCar]("CarsharingCar", b.origin)
	b.cars = cs
	return err
}

type OdhHalSharingLatest struct {
	ninja.OdhLatest
	Sname    string
	Provider string `json:"smetadata.company.shortName"`
}

func (p CarHAL) odhLatest(q siri.Query) ([]OdhHalSharingLatest, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = q.MaxSize()
	req.Repr = ninja.FlatNode
	req.StationTypes = []string{"CarsharingStation"}
	req.DataTypes = []string{"number-available"}
	req.Select = "mperiod,mvalue,mvalidtime,scode,sname,smetadata.company.shortName"
	req.Where = "sactive.eq.true"
	req.Where += fmt.Sprintf(",sorigin.eq.%s", p.origin)
	req.Where += apiBoundingBox(q)

	var res ninja.NinjaResponse[[]OdhHalSharingLatest]
	if err := ninja.Latest(req, &res); err != nil {
		slog.Error("Error retrieving parking state", "err", err)
		return res.Data, err
	}
	return res.Data, nil
}

func (p CarHAL) mapSiri(latest []OdhHalSharingLatest) []siri.FacilityCondition {
	ret := []siri.FacilityCondition{}

	for _, o := range latest {
		fc := siri.FacilityCondition{}
		fc.FacilityRef = netex.CreateID("Parking", o.Provider, o.Scode)
		fc.FacilityStatus.Status = siri.MapFacilityStatus(o.MValue, 1)
		fc.MonitoredCounting = &siri.MonitoredCounting{}
		fc.MonitoredCounting.CountingType = "availabilityCount"
		fc.MonitoredCounting.CountedFeatureUnit = "bays"
		fc.MonitoredCounting.Count = o.MValue

		ret = append(ret, fc)
	}

	return ret
}
func (p CarHAL) SiriFM(query siri.Query) (siri.FMData, error) {
	ret := siri.FMData{}
	idFilter := maybeIdMatch(query.FacilityRef(), netex.CreateID("Parking"))
	if len(query.FacilityRef()) > 0 && len(idFilter) == 0 {
		return ret, nil
	}

	l, err := p.odhLatest(query)
	if err != nil {
		return ret, err
	}
	ret.Conditions = filterFacilityConditions(p.mapSiri(l), idFilter)
	return ret, nil
}

func (b *CarHAL) MatchOperator(id string) bool {
	return id == b.GetOperator().Id
}
