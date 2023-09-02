APPID = io.github.jacalz.rymdport
NAME = rymdport

# If PREFIX isn't provided, default to /usr.
PREFIX ?= /usr

debug:
	go build -tags no_emoji -trimpath -o $(NAME)

release:
	go build -tags no_emoji -trimpath -ldflags="-s -w" -buildvcs=false -o $(NAME)

install:
	install -Dm00755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)
	install -Dm00644 internal/assets/icons/icon-512.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/512x512/apps/$(APPID).png
	install -Dm00644 internal/assets/icons/icon-256.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/256x256/apps/$(APPID).png
	install -Dm00644 internal/assets/icons/icon-128.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/128x128/apps/$(APPID).png
	install -Dm00644 internal/assets/icons/icon-64.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/64x64/apps/$(APPID).png
	install -Dm00644 internal/assets/icons/icon-48.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/48x48/apps/$(APPID).png
	install -Dm00644 internal/assets/icons/icon-32.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/32x32/apps/$(APPID).png
	install -Dm00644 internal/assets/icons/icon-24.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/24x24/apps/$(APPID).png
	install -Dm00644 internal/assets/icons/icon-16.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/16x16/apps/$(APPID).png
	install -Dm00644 internal/assets/svg/icon.svg $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/apps/$(APPID).svg
	install -Dm00644 internal/assets/unix/$(APPID).desktop $(DESTDIR)$(PREFIX)/share/applications/$(APPID).desktop
	install -Dm00644 internal/assets/unix/$(APPID).appdata.xml $(DESTDIR)$(PREFIX)/share/appdata/$(APPID).appdata.xml
	sudo gtk-update-icon-cache -f /usr/share/icons/hicolor/

uninstall:
	-rm $(DESTDIR)$(PREFIX)/bin/$(NAME)
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/512x512/apps/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/256x256/apps/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/128x128/apps/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/64x64/apps/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/48x48/apps/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/32x32/apps/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/24x24/apps/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/16x16/apps/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/apps/$(APPID).svg
	-rm $(DESTDIR)$(PREFIX)/share/applications/$(APPID).desktop
	-rm $(DESTDIR)$(PREFIX)/share/appdata/$(APPID).appdata.xml
