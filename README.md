<!--
SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>

SPDX-License-Identifier: CC0-1.0
-->

# NeTEx and SIRI FM export for STA MaaS4Italy project

[![REUSE Compliance](https://github.com/noi-techpark/java-boilerplate/actions/workflows/reuse.yml/badge.svg)](https://github.com/noi-techpark/odh-docs/wiki/REUSE#badges)

This API provides data from the Open Data Hub in the standard formats NeTEx (Italian Profile) and SIRI FM.

It acts as a wrapper and has no internal data.

## Getting started

Clone the repository and cd into it

We suggest using `docker compose up` to build and test.  
The docker image provides live code reloading  

To run natively
```sh
cd src
go run .
```
Refer to the [docker file](infrastructure/docker/Dockerfile) and [docker-compose](docker-compose.yml) for more details

## Tech stack
The application is written in Go with
 - [gin](https://github.com/gin-gonic/gin) for REST
 - [air](https://github.com/cosmtrek/air) to provide live reloading during development (not mandatory)

## Domain documentation
For documentation on the procotols, see the [initial issue](https://github.com/noi-techpark/sta-nap-export/issues/1)

The documents linked there are also in the [documentation](./documentation/) directory

The subrepo [netex-italian-profile](netex-italian-profile) tracks the [Italian netex profile repo](https://github.com/5Tsrl/netex-italian-profile) and contains examples and validation xsd files.  
The Italian profile is an extension of the base NeTEx specification

## Open Data Hub API calls
Example calls such as the ones used by this API are provided in [calls.http](calls.http)

## Adding datasets to export
Filters on which datasets are included are configured via [datasets.yml](src/config/datasets.yml).  

## Running automated tests
Native: From `/src` folder run `go test ./...`  
Docker: `docker compose up test`

## XML Validation
The [validation](./validation) directory has a script to validate the produced XML files against the Italian Profile xsd.  
[xmllint](https://xmllint.com/) is needed to run it locally, but we have also set up a dockerized version.  

```bash
docker compose --profile validate up --attach validate --abort-on-container-exit
```
This starts a local API, requests the XML files, adds the missing header parts, and validates them

The XML file that was validated is saved in `./validation/validate.xml`

## General Information
### Guidelines

Find [here](https://opendatahub.readthedocs.io/en/latest/guidelines.html) guidelines for developers.

### Support

For support, please contact [info@opendatahub.com](mailto:info@opendatahub.com).

### Contributing

If you'd like to contribute, please follow our [Getting
Started](https://github.com/noi-techpark/odh-docs/wiki/Contributor-Guidelines:-Getting-started)
instructions.

### Documentation

More documentation can be found at [https://opendatahub.readthedocs.io/en/latest/index.html](https://opendatahub.readthedocs.io/en/latest/index.html).

### License

The code in this project is licensed under the GNU AFFERO GENERAL PUBLIC LICENSE Version 3 license. See the [LICENSE.md](LICENSE.md) file for more information.

### REUSE

This project is [REUSE](https://reuse.software) compliant, more information about the usage of REUSE in NOI Techpark repositories can be found [here](https://github.com/noi-techpark/odh-docs/wiki/Guidelines-for-developers-and-licenses#guidelines-for-contributors-and-new-developers).

Since the CI for this project checks for REUSE compliance you might find it useful to use a pre-commit hook checking for REUSE compliance locally. The [pre-commit-config](.pre-commit-config.yaml) file in the repository root is already configured to check for REUSE compliance with help of the [pre-commit](https://pre-commit.com) tool.

Install the tool by running:
```bash
pip install pre-commit
```
Then install the pre-commit hook via the config file by running:
```bash
pre-commit install
```
