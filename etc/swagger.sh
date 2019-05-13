#!/usr/bin/env bash

# make sure to be in project root
go mod vendor

# now run swag
swag init -o src/docs --parseVendor