# Constants for cross compilation and packaging.
APPID = io.github.jacalz.wormhole_gui
NAME = wormhole-gui

# Variables for development.
GOBIN ?= ~/go/bin/
VERSION ?= dev

# If PREFIX isn't provided, we check for $(DESTDIR)/usr/local and use that if it exists.
# Otherwice we fall back to using /usr.
LOCAL != test -d $(DESTDIR)/usr/local && echo -n "/local" || echo -n ""
LOCAL ?= $(shell test -d $(DESTDIR)/usr/local && echo "/local" || echo "")
PREFIX ?= /usr$(LOCAL)

build:
	go build -ldflags="$(LDFLAGS)" -o $(NAME)

install:
	install -Dm00755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)
	install -Dm00644 internal/assets/icon/icon-512.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/512x512/$(APPID).png
	install -Dm00644 internal/assets/svg/icon.svg $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/$(APPID).svg
	install -Dm00644 internal/assets/unix/$(APPID).desktop $(DESTDIR)$(PREFIX)/share/applications/$(APPID).desktop
	install -Dm00644 internal/assets/unix/$(APPID).metainfo.xml $(DESTDIR)$(PREFIX)/share/metainfo/$(APPID).metainfo.xml

uninstall:
	-rm $(DESTDIR)$(PREFIX)/bin/$(NAME)
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/512x512/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/$(APPID).svg
	-rm $(DESTDIR)$(PREFIX)/share/applications/$(APPID).desktop
	-rm $(DESTDIR)$(PREFIX)/share/metainfo/$(APPID).metainfo.xml

freebsd:
	$(GOBIN)fyne-cross freebsd -arch amd64,arm64

darwin:
	$(GOBIN)fyne-cross darwin -arch amd64,arm64 -output $(NAME)

linux:
	$(GOBIN)fyne-cross linux -arch amd64,arm64

windows:
	$(GOBIN)fyne-cross windows -arch amd64

windows-debug:
	$(GOBIN)fyne-cross windows -console -arch amd64

bundle:
	# Move Linux package bundles to the root with correct naming.
	mv fyne-cross/dist/linux-amd64/wormhole-gui.tar.gz $(NAME)-$(VERSION)-linux-amd64.tar.gz
	mv fyne-cross/dist/linux-arm64/wormhole-gui.tar.gz $(NAME)-$(VERSION)-linux-arm64.tar.gz

	# Move FreeBSD package bundles to the root with correct naming.
	mv fyne-cross/dist/freebsd-amd64/wormhole-gui.tar.gz $(NAME)-$(VERSION)-freebsd-amd64.tar.gz
	mv fyne-cross/dist/freebsd-arm64/wormhole-gui.tar.gz $(NAME)-$(VERSION)-freebsd-arm64.tar.gz

	# Zip up the darwin packages with correct name and move to the root.
	(cd fyne-cross/dist/darwin-amd64/ && zip -r wormhole-gui-darwin-amd64.zip wormhole-gui.app/)
	mv fyne-cross/dist/darwin-amd64/wormhole-gui-darwin-amd64.zip $(NAME)-$(VERSION)-macOS-amd64.zip

	(cd fyne-cross/dist/darwin-arm64/ && zip -r wormhole-gui-darwin-arm64.zip wormhole-gui.app/)
	mv fyne-cross/dist/darwin-arm64/wormhole-gui-darwin-arm64.zip $(NAME)-$(VERSION)-macOS-arm64.zip

	# Move Windows package to the root with correct naming.
	mv fyne-cross/dist/windows-amd64/wormhole-gui.exe.zip $(NAME)-$(VERSION)-windows-amd64.zip

release: darwin freebsd linux windows bundle
	# Usage: make release VERSION=v3.x.x
	# If no version is specified, it will default to "dev" instead.
