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
	"slices"

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

	r.GET("/netex", netexEndpoint(netexAll))
	r.GET("/netex/parking", netexEndpoint(netexParking))
	r.GET("/netex/sharing", netexEndpoint(netexSharing))
	r.GET("/siri-lite/fm", siriLite(siriFM))
	r.GET("/siri-lite/fm/parking", siriLite(siriFMParking))
	r.GET("/siri-lite/fm/sharing", siriLite(siriFMSharing))

	r.GET("/health", health)
	r.Run()
}
func health(c *gin.Context) {
	c.Status(http.StatusOK)
}

type netexFn func() ([]netex.CompositeFrame, error)

func netexEndpoint(fn netexFn) func(*gin.Context) {
	return func(c *gin.Context) {
		comp, err := fn()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		n := netex.NewNetexFrame()
		n.DataObjects.Frames = append(n.DataObjects.Frames, comp...)
		prettyXML(c, http.StatusOK, n)
	}
}

func netexAll() ([]netex.CompositeFrame, error) {
	ret := []netex.CompositeFrame{}
	for _, s := range []netexFn{netexParking, netexSharing} {
		sr, err := s()
		if err != nil {
			return ret, err
		}
		ret = append(ret, sr...)
	}
	return ret, nil
}

func netexParking() ([]netex.CompositeFrame, error) {
	return netex.GetParking(provider.ParkingStatic)
}
func netexSharing() ([]netex.CompositeFrame, error) {
	return netex.GetSharing(provider.SharingBikesStatic, provider.SharingCarsStatic)
}

type siriFn func(siri.Query) (siri.Siri, error)

func siriLite(fn siriFn) func(*gin.Context) {
	return func(c *gin.Context) {
		query := siri.WrapQuery(c.Request.URL.Query())
		res, err := fn(query)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if wantsXml(c) {
			prettyXML(c, http.StatusOK, res)
		} else {
			c.JSONP(http.StatusOK, struct{ Siri siri.Siri }{Siri: res}) // wrap root level
		}
	}
}

func wantsXml(c *gin.Context) bool {
	accept := c.Request.Header.Values("Accept")
	return (slices.Contains(accept, "application/xml") &&
		!slices.Contains(accept, "application/json")) ||
		c.Query("format") == "xml"
}

func siriFM(query siri.Query) (siri.Siri, error) {
	return siri.FM(append(provider.ParkingRt, provider.SharingRt...), query)
}

func siriFMParking(query siri.Query) (siri.Siri, error) {
	return siri.FM(provider.ParkingRt, query)
}
func siriFMSharing(query siri.Query) (siri.Siri, error) {
	return siri.FM(provider.SharingRt, query)
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
