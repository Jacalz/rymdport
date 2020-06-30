# Constants for cross compilation and packaging.
APPID = com.github.jacalz.wormhole-gui
ICON = assets/icon-512.png
NAME = wormhole-gui

# Default path to the go binary directory.
GOBIN ?= ~/go/bin/

# If PREFIX isn't provided, we check for /usr/local and use that if it exists. Otherwice we fall back to using /usr.
ifneq ("$(wildcard /usr/local)","")
PREFIX ?= /usr/local
else
PREFIX ?= /usr
endif

build:
	go build -o $(name)

install:
	install -Dm00755 $(name) $(DESTDIR)$(PREFIX)/bin/$(name)
	install -Dm00644 assets/icon-512.png $(DESTDIR)$(PREFIX)/share/pixmaps/$(name).png
	install -Dm00644 assets/$(name).desktop $(DESTDIR)$(PREFIX)/share/applications/$(name).desktop

bundle:
	# Bundle the correct logo into sparta/src/bundled/bundled.go
	$(GOBIN)fyne bundle -package assets -name AppIcon assets/icon-512.png > assets/bundled.go

check:
	# Check the whole codebase for misspellings.
	$(GOBIN)misspell -w .

	# Run full formating on the code.
	gofmt -s -w .

	# Check the whole program for security issues.
	$(GOBIN)gosec ./...

	# Run staticcheck on the codebase.
	$(GOBIN)staticcheck -f stylish ./...

darwin:
	$(GOBIN)fyne-cross darwin -arch amd64 -app-id $(APPID) -icon $(ICON) -output $(NAME)

linux:
	$(GOBIN)fyne-cross linux -arch amd64 -app-id $(APPID) -icon $(ICON)

windows:
	$(GOBIN)fyne-cross windows -arch amd64 -app-id $(APPID) -icon $(ICON)

cross-compile: darwin linux windows