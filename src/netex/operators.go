// SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"log"
	"opendatahub/transmodel-api/config"

	"golang.org/x/exp/maps"
)

func mapByOrigin(p []config.OperatorCfg) map[string]config.OperatorCfg {
	ret := make(map[string]config.OperatorCfg)
	for _, o := range p {
		for _, origin := range o.Origin {
			ret[origin] = o
		}
	}
	return ret
}

func GetOperatorOrigins(c *config.Config, id string) []string {
	origins := []string{}
	for _, op := range c.Operators {
		if CreateID("Operator", op.Id) == id {
			origins = append(origins, op.Origin...)
		}
	}
	return origins
}

func GetOperator(c *config.Config, id string) Operator {
	mapped := mapByOrigin(c.Operators)
	cfg, found := mapped[id]
	if !found {
		log.Panicln("Unable to map operator. Probably got some origin that we shouldn't have?", id, maps.Keys(mapped))
	}

	o := Operator{}
	o.Id = CreateID("Operator", cfg.Id)
	o.Version = "1"
	o.PrivateCode = cfg.Id
	o.Name = cfg.Name
	o.ShortName = cfg.Name
	o.LegalName = cfg.Name
	o.TradingName = cfg.Name
	o.ContactDetails.Email = cfg.Email
	o.ContactDetails.Phone = cfg.Phone
	o.ContactDetails.Url = cfg.Url
	o.OrganisationType = "operator"
	o.Address.Id = CreateID("Address", cfg.Name)
	o.Address.CountryName = cfg.Country
	o.Address.Street = cfg.Street
	o.Address.Town = cfg.Town
	o.Address.PostCode = cfg.Postcode
	return o
}
