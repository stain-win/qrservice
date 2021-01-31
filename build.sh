#!/usr/bin/env bash
rm -rf dist
mkdir -p dist
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dist/qrservice gitlab.com/stain-win/qrservice/cmd/qrservice
