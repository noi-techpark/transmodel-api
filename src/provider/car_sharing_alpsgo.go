// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package provider

import (
	"fmt"
	"log/slog"
	"opendatahub/transmodel-api/comp"
	"opendatahub/transmodel-api/config"
	"opendatahub/transmodel-api/ninja"
	"opendatahub/transmodel-api/siri"
	"regexp"

	"github.com/noi-techpark/go-netex"
	"golang.org/x/exp/maps"
)

type odhAlpsGoCar []struct {
	ninja.OdhStation[alpsGoCarMeta]
	Pmetadata struct {
		Company struct {
			Uid       string
			ShortName string
			FullName  string
		}
	}
}

type CarSharingAlpsGo struct {
	cars     odhAlpsGoCar
	origin   string
	provider string
}

//	type alpsGoCarMeta struct {
//		Brand        string
//		Model        string
//		LicensePlate string
//		Features     *struct {
//			Doors           uint8
//			Seats           uint8
//			Chains          bool
//			Satnav          bool
//			Skirack         bool
//			Roofrack        bool
//			Childseat       string
//			Cyclerack       bool
//			Wintertyres     bool
//			Transmission    string
//			Cruisecontrol   bool
//			Trailerhitch    bool
//			Usbpowersockets bool
//		}
//	}
type alpsGoCarMeta struct {
	FuelType     string `json:"fuel_type"`
	Transmission string `json:"transmission"`
	VehicleModel struct {
		ModelName string `json:"model_name"`
	} `json:"vehicle_model"`
}

type alpsGoSharingMeta struct {
	CapacityMax int `json:"capacity_max"`
}

const ORIGIN_CAR_SHARING_ALPSGO = "AlpsGo"

func NewCarSharingAlpsGo() *CarSharingAlpsGo {
	b := CarSharingAlpsGo{}
	b.origin = "AlpsGo"
	b.provider = b.origin
	return &b
}

func (b *CarSharingAlpsGo) GetOperator() netex.Operator {
	return comp.GetOperator(&config.Cfg, b.origin)
}
func (b *CarSharingAlpsGo) StSharing() (comp.StSharingData, error) {
	ret := comp.StSharingData{}
	if err := b.fetch(); err != nil {
		return ret, err
	}

	// Operators
	o := b.GetOperator()
	ret.Operators = append(ret.Operators, o)

	// Modes of Operation
	m := netex.VehicleSharing{}
	m.Id = comp.CreateID("VehicleSharing", b.provider)
	m.Version = "1"
	sub := netex.Submode{}
	sub.Id = comp.CreateID("Submode", b.provider)
	sub.Version = "1"
	sub.TransportMode = "car"
	sub.SelfDriveSubmode = "hireCar"
	m.Submodes = append(m.Submodes, sub)
	ret.Modes = append(ret.Modes, m)

	models := make(map[string]netex.CarModelProfile)

	for _, c := range b.cars {
		modelname := c.Smeta.VehicleModel.ModelName
		p, found := models[modelname]
		if !found {
			// Car model profile
			p = netex.CarModelProfile{}
			p.Id = comp.CreateID("CarModelProfile", b.provider, modelname)
			p.Version = "1"
			// p.ChildSeat = c.Smeta.Features.Childseat
			// p.Seats = c.Smeta.Features.Seats
			// p.Doors = c.Smeta.Features.Doors
			p.Transmission = &c.Smeta.Transmission
			// p.CruiseControl = c.Smeta.Features.Cruisecontrol
			// p.SatNav = c.Smeta.Features.Satnav
			// p.AirConditioning = true
			// p.Convertible = false
			// p.UsbPowerSockets = c.Smeta.Features.Usbpowersockets
			// p.WinterTyres = c.Smeta.Features.Wintertyres
			// p.Chains = c.Smeta.Features.Chains
			// p.TrailerHitch = c.Smeta.Features.Trailerhitch
			// p.RoofRack = c.Smeta.Features.Roofrack
			// p.CycleRack = c.Smeta.Features.Cyclerack
			// p.SkiRack = c.Smeta.Features.Skirack
			models[modelname] = p
		}

		// Vehicles
		v := netex.Vehicle{}
		v.Id = comp.CreateID("Vehicle", b.provider, c.Scode)
		v.Version = "1"
		v.ValidBetween = comp.ValidAYear()
		v.Name = c.Sname
		v.ShortName = c.Sname
		v.PrivateCode = c.Scode
		v.RegistrationNumber = regexp.MustCompile(`\b\w+$`).FindString(c.Sname) // format AlpsGo 702 - GT029GC
		v.OperatorRef = comp.MkRef("Operator", o.Id)
		v.VehicleTypeRef = comp.MkRef("CarModelProfile", p.Id)
		ret.Vehicles = append(ret.Vehicles, v)
	}
	ret.CarModels = maps.Values(models)

	// Fleets = all Vehicles + operator
	f := netex.Fleet{}
	f.Id = comp.CreateID("Fleet", b.provider)
	f.Version = "1"
	f.ValidBetween = comp.ValidAYear()
	members := []netex.Ref{}
	for _, v := range ret.Vehicles {
		members = append(members, comp.MkRef("Vehicle", v.Id))
	}
	if len(members) > 0 {
		f.Members = &members
	}
	f.OperatorRef = comp.MkRef("Operator", o.Id)
	ret.Fleets = append(ret.Fleets, f)

	// Mobility services = Fleet + mode
	s := netex.VehicleSharingService{}
	s.Id = comp.CreateID("VehicleSharingService", b.provider)
	s.Version = "1"
	s.VehicleSharingRef = comp.MkRef("VehicleSharing", m.Id)
	s.FloatingVehicles = false
	for _, fl := range ret.Fleets {
		s.Fleets = append(s.Fleets, comp.MkRef("Fleet", fl.Id))
	}
	ret.Services = append(ret.Services, s)

	// Constraint zone
	c := netex.MobilityServiceConstraintZone{}
	c.Id = comp.CreateID("MobilityServiceConstraintZone", b.provider)
	c.Version = "1"
	c.GmlPolygon.Id = b.provider
	c.GmlPolygon.Polygon = config.GML_PROVINCE_BZ
	c.VehicleSharingRef = comp.MkRef("VehicleSharingService", s.Id)
	ret.Constraints = append(ret.Constraints, c)

	// Sharing as Parking (for SIRI reference)
	ss, err := FetchOdhStations[[]ninja.OdhStation[alpsGoSharingMeta]]("CarsharingStation", b.origin)
	if err != nil {
		return ret, err
	}
	for _, s := range ss {
		p := netex.Parking{}
		p.Id = comp.CreateID("Parking", b.provider, s.Scode)
		p.Version = "1"
		p.ShortName = s.Sname
		p.Centroid.Location.Longitude = s.Scoord.X
		p.Centroid.Location.Latitude = s.Scoord.Y
		p.OperatorRef = comp.MkRef("Operator", o.Id)
		p.GmlPolygon = nil
		p.Entrances = nil
		p.ParkingType = "rentalCarParking"
		p.ParkingVehicleTypes = "car"
		p.ParkingLayout = "undefined"
		p.ProhibitedForHazardousMaterials = netex.Just(true)
		p.RechargingAvailable = nil
		p.Secure = nil
		p.ParkingReservation = "reservationRequired"
		p.ParkingProperties = nil

		p.Name = s.Sname
		p.PrincipalCapacity = int32(s.Smeta.CapacityMax)
		p.TotalCapacity = int32(s.Smeta.CapacityMax)
		ret.Parkings = append(ret.Parkings, p)
	}

	return ret, nil
}

func (b *CarSharingAlpsGo) fetch() error {
	cs, err := FetchOdhStations[odhAlpsGoCar]("CarsharingCar", b.origin)
	b.cars = cs
	return err
}

type OdhAlpsGoSharingLatest struct {
	ninja.OdhLatest
	Sname string
}

func (p CarSharingAlpsGo) odhLatest(q siri.Query) ([]OdhAlpsGoSharingLatest, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = q.MaxSize()
	req.Repr = ninja.FlatNode
	req.StationTypes = []string{"CarsharingStation"}
	req.DataTypes = []string{"number-available"}
	req.Select = "mperiod,mvalue,mvalidtime,scode,sname"
	req.Where = "sactive.eq.true"
	req.Where += fmt.Sprintf(",sorigin.eq.%s", p.origin)
	req.Where += apiBoundingBox(q)

	var res ninja.NinjaResponse[[]OdhAlpsGoSharingLatest]
	if err := ninja.Latest(req, &res); err != nil {
		slog.Error("Error retrieving parking state", "err", err)
		return res.Data, err
	}
	return res.Data, nil
}

func (p CarSharingAlpsGo) mapSiri(latest []OdhAlpsGoSharingLatest) []siri.FacilityCondition {
	ret := []siri.FacilityCondition{}

	for _, o := range latest {
		fc := siri.FacilityCondition{}
		fc.FacilityRef = comp.CreateID("Parking", p.provider, o.Scode)
		fc.FacilityStatus.Status = siri.MapFacilityStatus(o.MValue, 1)
		fc.MonitoredCounting = &siri.MonitoredCounting{}
		fc.MonitoredCounting.CountingType = "availabilityCount"
		fc.MonitoredCounting.CountedFeatureUnit = "bays"
		fc.MonitoredCounting.Count = o.MValue

		ret = append(ret, fc)
	}

	return ret
}
func (p CarSharingAlpsGo) SiriFM(query siri.Query) (siri.FMData, error) {
	ret := siri.FMData{}
	idFilter := maybeIdMatch(query.FacilityRef(), comp.CreateID("Parking"))
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

func (b *CarSharingAlpsGo) MatchOperator(id string) bool {
	return id == b.GetOperator().Id
}
