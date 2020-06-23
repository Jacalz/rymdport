# Constants for cross compilation and packaging.
appID = com.github.jacalz.wormhole-gui
name = wormhole-gui

# Default path to the go binary directory.
GOBIN ?= ~/go/bin/

check:
	# Check the whole codebase for misspellings.
	${GOBIN}misspell -w .

	# Run full formating on the code.
	gofmt -s -w .

	# Check the whole program for security issues.
	${GOBIN}gosec ./...

	# Run staticcheck on the codebase.
	${GOBIN}staticcheck -f stylish ./...

darwin:
	${GOBIN}fyne-cross darwin -arch amd64 -app-id ${appID} -output ${name}

linux:
	${GOBIN}fyne-cross linux -arch amd64 -app-id ${appID}

windows:
	${GOBIN}fyne-cross windows -arch amd64 -app-id ${appID}

cross-compile: darwin linux windows