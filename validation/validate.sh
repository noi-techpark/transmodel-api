#!/bin/bash

# SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
#
# SPDX-License-Identifier: CC0-1.0

template=$(<delivery_container.xml)
endpoint="${ENDPOINT:-localhost:8000}"

tmpfile="${TMPFILE:-validate.xml}"

function vUrl () {
    content=""
    for arg in "$@"
    do
        content+=`curl $arg | xmllint --xpath "//*[local-name()='CompositeFrame']" -`
    done

    xmllint --format - <<<"${template//PLACEHOLDER/$content}" > $tmpfile
    xmllint --noout --schema ../netex-italian-profile/xsd/NeTEx_publication_Lev4.xsd $tmpfile
}

vUrl $endpoint/netex/parking \
&& vUrl $endpoint/netex/sharing \
&& vUrl $endpoint/netex/parking $endpoint/netex/sharing
