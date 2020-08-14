# wormhole-gui

Wormhole-gui is a graphical interface for magic-wormhole. It uses the Go implementation [wormhole-william](https://github.com/psanford/wormhole-william) along with [fyne](https://github.com/fyne-io/fyne) for the graphical interface. The initial version was built in less than one day to show how quick and easy it is to use fyne for developing applications. The application has since developed into a more feature-full cross-platform application for sharing files and text between devices on the local network. 

<p align="center">
  <img src="assets/screenshot.png" />
</p>

It is built using the following Go modules:
- [fyne](https://github.com/fyne-io/fyne) (version 1.3.3 or later)
- [wormhole-william](https://github.com/psanford/wormhole-william) (version 1.0.4 or later)

## Future work
This list is part of a larger set of planned features. The list is in now way complete and doesn't mean that other features won't be incorporated.
- Folder sharing (blocked until fyne has folder select support in file-picker)

## Requirements

Wormhole-gui compiles into a statically linked binary with no runtime dependencies.

## Downloads

Please visit the [release page](https://github.com/Jacalz/wormhole-gui/releases) for downloading the latest releases.
Versions for Linux, MacOS and Windows on the x86-64 architecture are currently available. Support for ARM64 and BSD systems might come in the future.

## Building

Systems that have a recent [Go](https://golang.org) compiler and the [required prequsites for Fyne](https://fyne.io/develop/) should be able to just build the project using `go build` in the project root:
```bash
go build 
```

The project can also be built and installed using GNU Make (installing is only supported on Linux currently):
```bash
make
sudo make install
```

## Contributing

Contributions are strongly appreciated. Everything from creating bug reports to contributing code will help the project a lot, so please feel free to help in any way, shape or form that you feel comfortable doing.

## License
- Wormhole-gui is licensed under `GNU GENERAL PUBLIC LICENSE Version 3`.
