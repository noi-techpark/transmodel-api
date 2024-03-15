// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func main() {
	InitLogger()
	r := gin.Default()

	// Enable slog logging for gin framework
	// https://github.com/samber/slog-gin
	r.Use(sloggin.New(slog.Default()))
	r.Use(gin.Recovery())

	r.GET("/parking", parking)
	r.GET("/sharing/static", sharing)
	r.GET("/sharing/rt/:id", realtime)
	r.GET("/health", health)
	r.Run()
}
func health(c *gin.Context) {
	c.Status(http.StatusOK)
}
func parking(c *gin.Context) {
	c.XML(http.StatusOK, gin.H{"msg": "Hello parking world"})
}
func sharing(c *gin.Context) {
	c.XML(http.StatusOK, gin.H{"msg": "Hello sharing world"})
}
func realtime(c *gin.Context) {
	c.XML(http.StatusOK, gin.H{"msg": "Hello realtime world", "station": c.Param("id")})
}
