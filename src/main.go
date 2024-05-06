// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"encoding/xml"
	"log/slog"
	"net/http"
	nParking "opendatahub/sta-nap-export/netex/parking"
	nSharing "opendatahub/sta-nap-export/netex/sharing"
	"opendatahub/sta-nap-export/siri"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	sloggin "github.com/samber/slog-gin"
)

func main() {
	InitLogger()
	r := gin.New()

	if os.Getenv("GIN_LOG") == "PRETTY" {
		r.Use(gin.Logger())
	} else {
		// Enable slog logging for gin framework
		// https://github.com/samber/slog-gin
		r.Use(sloggin.New(slog.Default()))
	}

	r.Use(gin.Recovery())

	r.GET("/netex/parking", parking)
	r.GET("/netex/sharing", sharing)
	r.GET("/siri/fm/:id", realtime)
	r.GET("/health", health)
	r.Run()
}
func health(c *gin.Context) {
	c.Status(http.StatusOK)
}
func parking(c *gin.Context) {
	res, err := nParking.GetParking()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	prettyXML(c, http.StatusOK, res)
}
func sharing(c *gin.Context) {
	res, err := nSharing.GetSharing()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	prettyXML(c, http.StatusOK, res)
}
func realtime(c *gin.Context) {
	scode := c.Param("id")
	res, err := siri.Parking(scode)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	prettyXML(c, http.StatusOK, res)
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
