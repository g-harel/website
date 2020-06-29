#!/bin/sh

# Update build function package's deps to include the most
# recent internal changes outside of "../functions".
cd "./functions"
go get -u "github.com/g-harel/website"
go mod tidy
