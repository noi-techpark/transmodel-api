// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log/slog"
	"net/http"
	"opendatahub/sta-nap-export/netex"
	"os"

	"github.com/gin-gonic/gin"
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
	res, err := netex.GetNetexParking()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.XML(http.StatusOK, res)
}
func sharing(c *gin.Context) {
	c.XML(http.StatusOK, gin.H{"msg": "Hello sharing world"})
}
func realtime(c *gin.Context) {
	c.XML(http.StatusOK, gin.H{"msg": "Hello realtime world", "station": c.Param("id")})
}
