# Constants for cross compilation and packaging.
APPID = com.github.jacalz.wormhole-gui
ICON = internal/assets/icon/icon-512.png
NAME = wormhole-gui

# Default path to the go binary directory.
GOBIN ?= ~/go/bin/

# If PREFIX isn't provided, we check for $(DESTDIR)/usr/local and use that if it exists.
# Otherwice we fall back to using /usr.
LOCAL != test -d $(DESTDIR)/usr/local && echo -n "/local" || echo -n ""
LOCAL ?= $(shell test -d $(DESTDIR)/usr/local && echo "/local" || echo "")
PREFIX ?= /usr$(LOCAL)

build:
	go build -o $(NAME)

install:
	install -Dm00755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)
	install -Dm00644 $(ICON) $(DESTDIR)$(PREFIX)/share/pixmaps/$(NAME).png
	install -Dm00644 internal/assets/$(NAME).desktop $(DESTDIR)$(PREFIX)/share/applications/$(NAME).desktop

uninstall:
	-rm $(DESTDIR)$(PREFIX)/share/applications/$(NAME).desktop
	-rm $(DESTDIR)$(PREFIX)/bin/$(NAME)
	-rm $(DESTDIR)$(PREFIX)/share/pixmaps/$(NAME).png

check:
	# Check the whole codebase for misspellings.
	$(GOBIN)misspell -w .

	# Run full formating on the code.
	gofmt -s -w .

	# Check the whole program for security issues.
	$(GOBIN)gosec ./...

	# Run staticcheck on the codebase.
	$(GOBIN)staticcheck -f stylish ./...

freebsd:
	$(GOBIN)fyne-cross freebsd -arch amd64 -app-id $(APPID) -icon $(ICON)

darwin:
	$(GOBIN)fyne-cross darwin -arch amd64 -app-id $(APPID) -icon $(ICON) -output $(NAME)

linux:
	$(GOBIN)fyne-cross linux -arch amd64,arm64 -app-id $(APPID) -icon $(ICON)

windows:
	$(GOBIN)fyne-cross windows -arch amd64 -app-id $(APPID) -icon $(ICON)

windows-debug:
	$(GOBIN)fyne-cross windows -console -arch amd64 -app-id $(APPID) -icon $(ICON)

cross-compile: darwin freebsd linux windows