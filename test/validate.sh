#!/bin/bash

# SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
#
# SPDX-License-Identifier: CC0-1.0

endpoint="${ENDPOINT:-localhost:8000}"

tmpfile="${TMPFILE:-validate.xml}"

function vUrl () {
    xsddir=../netex-italian-profile/xsd
    content=`curl $1`

    xmllint --format - <<<"$content" > $tmpfile
    xmllint --noout  --schema $xsddir/NeTEx_publication_Lev4.xsd $tmpfile
}

vUrl $endpoint/netex/parking \
&& vUrl $endpoint/netex/sharing \
&& vUrl $endpoint/netex/flights \
&& vUrl $endpoint/netex
