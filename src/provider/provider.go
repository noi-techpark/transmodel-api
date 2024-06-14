// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later
package provider

import "opendatahub/sta-nap-export/ninja"

func odhMob[T any](tp string, origin string) (T, error) {
	req := ninja.DefaultNinjaRequest()
	req.Limit = -1
	req.StationTypes = []string{tp}
	req.Where = `sorigin.eq.` + origin + `,sactive.eq.true`
	var res ninja.NinjaResponse[T]
	err := ninja.StationType(req, &res)
	return res.Data, err
}
