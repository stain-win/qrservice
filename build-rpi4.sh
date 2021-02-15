#!/usr/bin/env bash
rm -rf dist
mkdir -p dist
CGO_ENABLED=0 env GOARCH=arm GOARM=7 GOOS=linux go build -o dist/qrservice github.com/stain-win/qrservice/cmd/qrservice
