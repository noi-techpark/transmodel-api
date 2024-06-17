// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"encoding/xml"
	"log/slog"
	"net/http"
	"opendatahub/sta-nap-export/config"
	"opendatahub/sta-nap-export/netex"
	"opendatahub/sta-nap-export/provider"
	"opendatahub/sta-nap-export/siri"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	sloggin "github.com/samber/slog-gin"
)

func main() {
	InitLogger()
	config.InitConfig()

	r := gin.New()

	if os.Getenv("GIN_LOG") == "PRETTY" {
		r.Use(gin.Logger())
	} else {
		// Enable slog logging for gin framework
		// https://github.com/samber/slog-gin
		r.Use(sloggin.New(slog.Default()))
	}

	r.Use(gin.Recovery())

	r.GET("/netex/parking", netexPark)
	r.GET("/netex/sharing", netexSharing)
	r.GET("/siri/fm/parking", siriParking)

	r.GET("/health", health)
	r.Run()
}
func health(c *gin.Context) {
	c.Status(http.StatusOK)
}

func netexPark(c *gin.Context) {
	res, err := netex.GetParking()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	prettyXML(c, http.StatusOK, res)
}

func netexSharing(c *gin.Context) {
	bikeProviders := []netex.StSharing{&provider.BikeBz{}, &provider.BikeMe{}, &provider.BikePapin{}}
	carProviders := []netex.StSharing{&provider.CarHAL{}}
	res, err := netex.GetSharing(bikeProviders, carProviders)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	prettyXML(c, http.StatusOK, res)
}

func siriParking(c *gin.Context) {
	res, err := siri.Parking()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSONP(http.StatusOK, res)
}

func prettyXML(c *gin.Context, code int, object any) {
	// Due to a request, this renders the xml in a pretty format.
	// Once this is production ready, should probably switch back to just c.XML(...)
	data, err := xml.MarshalIndent(object, "", "  ")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Render(code, render.Data{Data: data, ContentType: "application/xml; charset=utf-8"})
}
