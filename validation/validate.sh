#!/bin/bash

# SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
#
# SPDX-License-Identifier: CC0-1.0

template=$(<delivery_container.xml)
endpoint="${ENDPOINT:-localhost:8000}"
tmpfile="${TMPFILE:-validate.xml}"

function vUrl () {
    parking=`curl $1`
    xmllint --format - <<<"${template//PLACEHOLDER/$parking}" > $tmpfile
    xmllint --noout --schema ../netex-italian-profile/xsd/NeTEx_publication_Lev5.xsd $tmpfile
}

vUrl $endpoint/netex/parking \
&& vUrl $endpoint/netex/sharing

