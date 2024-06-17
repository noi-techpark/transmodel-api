// SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Operators []OperatorCfg
	Dataset   DatasetCfg
}

var Cfg Config

type DatasetCfg struct {
	Parking struct {
		Origins []string
	}
}

type OperatorCfg struct {
	Origin   []string
	Id       string
	Name     string
	Email    string
	Phone    string
	Url      string
	Street   string
	Town     string
	Postcode string
	Country  string
}

func InitConfig() {
	Cfg = *ReadConfig()
}
func ReadConfig() *Config {
	cfg := Config{}
	readYaml(fixRelPath("config", "operators.yml"), &cfg.Operators)
	readYaml(fixRelPath("config", "datasets.yml"), &cfg.Dataset)
	return &cfg
}

func readYaml(path string, o any) {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Panicln("Cannot open config", path, err)
	}

	if err := yaml.Unmarshal(f, o); err != nil {
		log.Panicln("Cannot unmarshal Operators config", path, err)
	}
}

func fixRelPath(path ...string) string {
	cwd, _ := os.Getwd()
	// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
	// Relative paths are a pain in the butt with unit tests because they always execute from the module they are in
	// This is a hack to always start from root folder and compose the full "absolute" path
	if testing.Testing() {
		_, b, _, _ := runtime.Caller(0)
		cwd = filepath.Join(filepath.Dir(b), "..")
	}

	return filepath.Join(append([]string{cwd}, path...)...)
}

func (c *Config) ParkingOrigins() []string {
	return c.Dataset.Parking.Origins
}
