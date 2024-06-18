#!/bin/bash

# SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
#
# SPDX-License-Identifier: CC0-1.0

# Run tests and validations

exitcode=0

(cd src; go test ./...)
result=$?
if [ $result -ne 0 ]; then
    echo Go tests failed!
    exitcode=$result
fi

docker compose --profile validate up --attach validate --abort-on-container-exit
result=$?
if [ $result -ne 0 ]; then
    echo Live validation failed!
    exitcode=$result
fi

if [ $exitcode -ne 0 ]; then
    exit $exitcode
fi

echo "All tests OK!"
exit 0