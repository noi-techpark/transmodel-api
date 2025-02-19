// SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: CC0-1.0

package provider

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"opendatahub/transmodel-api/netex"
	"sync"
	"time"
)

type FlightsGtfs struct {
	Url       string
	NUTS      string
	Vat       string
	Company   string
	maxAge    time.Duration
	cache     *netex.StFlightData
	cacheTime time.Time
}

func NewFlightsSkyalps() *FlightsGtfs {
	return &FlightsGtfs{
		Url:     "https://gtfs.api.opendatahub.com/v1/dataset/skyalps-flight-data/raw",
		NUTS:    "IT:ITH10",
		Company: "SKYALPS",
		Vat:     "03067170211",
		maxAge:  8 * time.Hour,
	}
}

func (FlightsGtfs) getRemoteGtfs(url string) (*[]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	gtfs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &gtfs, nil
}

func (f FlightsGtfs) gtfs2Netex(gtfs *[]byte) (*[]byte, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	form, err := writer.CreateFormFile("file", "flights.gtfs")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(form, bytes.NewReader(*gtfs))
	if err != nil {
		return nil, err
	}

	writer.WriteField("nuts", f.NUTS)
	writer.WriteField("az", f.Company)
	writer.WriteField("vat", f.Vat)
	writer.WriteField("version", time.Now().Format("060102"))

	if err := writer.Close(); err != nil {
		return nil, err
	}

	// call docker internal conversion API
	req, err := http.NewRequest("POST", `http://gtfs2netex:8080`, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
}

func (fs FlightsGtfs) fromNetex(netexXml *[]byte) (netex.StFlightData, error) {
	ret := netex.StFlightData{}
	var n netex.NetexFrame

	err := xml.Unmarshal(*netexXml, &n)
	if err != nil {
		return ret, err
	}

	for _, do := range n.DataObjects {
		for _, rf := range do.Frames.ResourceFrame {
			ret.Operators = append(ret.Operators, *rf.Operators...)
		}
		for _, sf := range do.Frames.SiteFrame {
			ret.StopPlaces = append(ret.StopPlaces, sf.StopPlaces...)
		}
		for _, sf := range do.Frames.ServiceFrame {
			ret.JourneyPatterns = append(ret.JourneyPatterns, sf.JourneyPatterns...)
			ret.Lines = append(ret.Lines, sf.Lines...)
			ret.Routes = append(ret.Routes, sf.Routes...)
			ret.ScheduledStopPoints = append(ret.ScheduledStopPoints, sf.ScheduledStopPoints...)
			ret.ServiceLinks = append(ret.ServiceLinks, sf.ServiceLinks...)
			ret.StopAssignments = append(ret.StopAssignments, sf.StopAssignments...)
		}
		for _, cf := range do.Frames.ServiceCalendarFrame {
			ret.ServiceCalendars = append(ret.ServiceCalendars, cf.ServiceCalendar...)
		}
		for _, tf := range do.Frames.TimetableFrame {
			ret.VehicleJourneys = append(ret.VehicleJourneys, tf.VehicleJourneys...)
		}
	}

	return ret, nil
}

var lock sync.Mutex

func (fs *FlightsGtfs) StFlights() (netex.StFlightData, error) {
	if err := func() error {
		lock.Lock()
		defer lock.Unlock()

		// only do this one at a time
		if time.Since(fs.cacheTime) > fs.maxAge {
			gtfs, err := fs.getRemoteGtfs(fs.Url)
			if err != nil {
				return err
			}
			netexIn, err := fs.gtfs2Netex(gtfs)
			if err != nil {
				return err
			}
			netexOut, err := fs.fromNetex(netexIn)
			if err != nil {
				return err
			}
			fs.cache = &netexOut
			fs.cacheTime = time.Now()
		}
		return nil
	}(); err != nil {
		return netex.StFlightData{}, err
	}
	return *fs.cache, nil
}
