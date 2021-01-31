package qrservice

import (
	"os"
	"path/filepath"
)

var (
	// Name is name of this service as it is known to the outside world.
	Name = "qrservice"
	// Description is a short service description.
	Description = "QR code server"
	// Author is list of authors of a service that should be notified in event of any trouble.
	Author = "Miroslav Mitrovic <stainnn@gmai.com>"
	// Version is a SemVer number
	// Set the value in build time using: -ldflags "-X go.sbgenomics.com/tcpserve.Version=$(VERSION)"
	Version = "0"
	// BuildInfo is the version control current build revision.
	// Set the value in build time using: -ldflags "-X go.sbgenomics.com/tcpserve.BuildInfo=$(shell git describe --long --dirty --always)"
	BuildInfo = ""
	// BuildTime is the date of the build.
	// Set the value in build time using: -ldflags "-X go.sbgenomics.com/tcpserve.BuildTime=$(date +%FT%T%z)"
	BuildTime = ""
	// BaseDir is the root directory of the tcpserve executable.
	BaseDir = func() string {
		baseDir := filepath.Dir(os.Args[0])
		baseDir, err := filepath.Abs(baseDir)
		if err != nil {
			panic(err)
		}
		return baseDir
	}()
)
