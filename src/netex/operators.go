// SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"log"

	"golang.org/x/exp/maps"
)

func mapByOrigin(p []operatorCfg) map[string]operatorCfg {
	ret := make(map[string]operatorCfg)
	for _, o := range p {
		for _, origin := range o.Origin {
			ret[origin] = o
		}
	}
	return ret
}

func (c *Config) GetOperator(id string) Operator {
	mapped := mapByOrigin(c.operators)
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
