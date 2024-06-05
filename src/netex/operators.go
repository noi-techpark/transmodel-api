// SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package netex

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gocarina/gocsv"
)

type operatorCfg struct {
	Origin   string `csv:"origin"`
	Email    string `csv:"email"`
	Phone    string `csv:"phone"`
	Url      string `csv:"url"`
	Street   string `csv:"street"`
	Town     string `csv:"town"`
	Postcode string `csv:"postcode"`
	Country  string `csv:"country"`
}

func readOps(path string) []operatorCfg {
	f, err := os.Open(path)
	if err != nil {
		wd, _ := os.Getwd()
		log.Panicln("Cannot open Operators csv.", wd, err)
	}
	defer f.Close()

	ops := []operatorCfg{}
	if err := gocsv.UnmarshalFile(f, &ops); err != nil {
		log.Panic("Cannot unmarshal Operators csv", err)
	}
	return ops
}

var ops []operatorCfg

func getCsvPath() string {
	// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
	// Relative paths are a pain in the butt with unit tests because they always execute from the module they are in
	// This is a hack to always start from root folder and compose the full "absolute" path
	if testing.Testing() {
		_, b, _, _ := runtime.Caller(0)
		root := filepath.Join(filepath.Dir(b), "../..")
		return filepath.Join(root, "src", "resources", "operators.csv")
	}

	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "resources", "operators.csv")
}

func mapByOrigin(p []operatorCfg) map[string]operatorCfg {
	ret := make(map[string]operatorCfg)
	for _, o := range p {
		ret[o.Origin] = o
	}
	return ret
}
func opsByOrigin() map[string]operatorCfg {
	if ops == nil {
		ops = readOps(getCsvPath())
	}
	return mapByOrigin(ops)
}

func GetOperator(id string) Operator {
	cfg, found := opsByOrigin()[id]
	if !found {
		log.Panic("Unable to map operator. Probably got some origin that we shouldn't have?")
	}

	o := Operator{}
	o.Id = CreateID("Operator", id)
	o.Version = "1"
	o.PrivateCode = id
	o.Name = id
	o.ShortName = id
	o.LegalName = id
	o.TradingName = id
	o.ContactDetails.Email = cfg.Email
	o.ContactDetails.Phone = cfg.Phone
	o.ContactDetails.Url = cfg.Url
	o.OrganisationType = "operator"
	o.Address.Id = CreateID("Address", id)
	o.Address.CountryName = cfg.Country
	o.Address.Street = cfg.Street
	o.Address.Town = cfg.Town
	o.Address.PostCode = cfg.Postcode
	return o
}
