// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>

// SPDX-License-Identifier: AGPL-3.0-or-later

package ninja

import (
	"encoding/json"
	"os"
	"reflect"
)

// For testing purposes, set this function and it will be used to retrieve a request instead of http
var TestReqHook func(*NinjaRequest) (any, error)

func runReqHook(req *NinjaRequest, result any) error {
	r, err := TestReqHook(req)
	if err != nil {
		return err
	}
	// unholy memcpy: sets the memory at p to the value of pv. Obviously they have to be the same type
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(r).Elem())
	return nil
}

func LoadJsonFile[T any](f string) (*NinjaResponse[T], error) {
	r := &NinjaResponse[T]{}
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(b, r)
	return r, nil
}

// func fixRelPath(path string) string {
// 	_, b, _, _ := runtime.Caller(0)
// 	cwd := filepath.Join(filepath.Dir(b), "..")
// 	return filepath.Join(cwd, path)
// }
