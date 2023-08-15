APPID = io.github.jacalz.rymdport
NAME = rymdport

# If PREFIX isn't provided, we check for $(DESTDIR)/usr/local and use that if it exists.
# Otherwice we fall back to using /usr.
LOCAL != test -d $(DESTDIR)/usr/local && echo -n "/local" || echo -n ""
LOCAL ?= $(shell test -d $(DESTDIR)/usr/local && echo "/local" || echo "")
PREFIX ?= /usr$(LOCAL)

debug:
	go build -tags no_emoji -trimpath -o $(NAME)

release:
	go build -tags no_emoji -trimpath -ldflags="-s -w" -buildvcs=false -o $(NAME)

install:
	install -Dm00755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)
	install -Dm00644 internal/assets/icon/icon-512.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/512x512/$(APPID).png
	install -Dm00644 internal/assets/svg/icon.svg $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/$(APPID).svg
	install -Dm00644 internal/assets/unix/$(APPID).desktop $(DESTDIR)$(PREFIX)/share/applications/$(APPID).desktop
	install -Dm00644 internal/assets/unix/$(APPID).appdata.xml $(DESTDIR)$(PREFIX)/share/appdata/$(APPID).appdata.xml

uninstall:
	-rm $(DESTDIR)$(PREFIX)/bin/$(NAME)
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/512x512/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/$(APPID).svg
	-rm $(DESTDIR)$(PREFIX)/share/applications/$(APPID).desktop
	-rm $(DESTDIR)$(PREFIX)/share/appdata/$(APPID).appdata.xml
