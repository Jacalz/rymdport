# Constants for cross compilation and packaging.
appID = com.github.jacalz.wormhole-gui
icon = assets/icon-512.png
name = wormhole-gui

# Default path to the go binary directory.
GOBIN ?= ~/go/bin/

bundle:
	# Bundle the correct logo into sparta/src/bundled/bundled.go
	${GOBIN}fyne bundle -package assets -name AppIcon assets/icon-512.png > assets/bundled.go

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
	${GOBIN}fyne-cross darwin -arch amd64 -app-id ${appID} -icon ${icon}  -output ${name}

linux:
	${GOBIN}fyne-cross linux -arch amd64 -app-id ${appID} -icon ${icon}

windows:
	${GOBIN}fyne-cross windows -arch amd64 -app-id ${appID} -icon ${icon}

cross-compile: darwin linux windows