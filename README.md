# Rymdport

[![Static Analysis](https://github.com/Jacalz/rymdport/actions/workflows/static_analysis.yml/badge.svg?branch=main)](https://github.com/Jacalz/rymdport/actions/workflows/static_analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Jacalz/rymdport/v3)](https://goreportcard.com/report/github.com/Jacalz/rymdport/v3)
[![version](https://img.shields.io/github/v/tag/Jacalz/rymdport?label=version)](https://github.com/Jacalz/rymdport/releases/latest)

Rymdport (formerly wormhole-gui) is a cross-platform application that lets you easily and safely share files, folders, and text between devices.
The data is sent peer-to-peer securely using end-to-end encryption using the same protocol as [magic-wormhole](https://github.com/magic-wormhole/magic-wormhole). This means that Rymdport can talk not only to itself, but also to other wormhole clients.

The transfers are implemented using [wormhole-william](https://github.com/psanford/wormhole-william) as a native [Go](https://go.dev/) implementation of magic-wormhole. As a result, Rymdport compiles into a native binary with no runtime dependencies while also outperforming the reference implementation of magic-wormhole.

<p align="center">
  <img src="internal/assets/screenshot1.png" />
</p>

#### Rymdport is built using the following Go modules:
- [fyne](https://github.com/fyne-io/fyne) (version 2.1.3)
- [wormhole-william](https://github.com/psanford/wormhole-william) (version 1.0.6)
- [compress](https://github.com/klauspost/compress) (version 1.15.0)

The initial version was built in less than one day to show how quick and easy it is to use [Fyne](https://github.com/fyne-io/fyne) for developing applications.

## Sponsoring

Rymdport is an open source project that is provided free of charge, and that will continue to be the case forever. If you use this project and appreciate the work being put into it, please consider supporting its development through [GitHub Sponsors](https://github.com/sponsors/Jacalz). This is in no way a requirement, but would be greatly appreciated and would allow for even more improvements to come further down the road.

## Requirements

Rymdport compiles into a statically linked binary with no explicit runtime dependencies.
Compiling requires a [Go](https://go.dev) compiler (1.15 or later officially supported, but older versions might work) and the [prerequisites for Fyne](https://developer.fyne.io/started/).

## Downloads

Please visit the [release page](https://github.com/Jacalz/rymdport/releases) to download the latest release.
Pre-built binaries are available for FreeBSD, Linux, macOS (`x86-64` and `arm64`) and Windows (`x86-64`).

For Linux users, Rymdport is also avaliable as a Flatpak on [Flathub](https://flathub.org/apps/details/io.github.jacalz.rymdport):

<a href='https://flathub.org/apps/details/io.github.jacalz.rymdport'><img width='200' alt='Download on Flathub' src='https://flathub.org/assets/badges/flathub-badge-en.png'/></a>

The following distributions also have binary packages available through their respective package managers:

[![Packaging status](https://repology.org/badge/vertical-allrepos/rymdport.svg)](https://repology.org/project/rymdport/versions)

## Building

Systems with the compile-time requirements satisfied can build the project using `go build` in the project root:
```bash
go build
```

The project is available in the [Fyne Apps Listing](https://apps.fyne.io/apps/rymdport.html) and can be installed either using the `fyne get` command or using the [Fyne Apps Installer](https://apps.fyne.io/apps/io.fyne.apps.html).

Installation can also be performed using GNU Make (installing this way is currently only supported on Linux and BSD):
```bash
make
sudo make install
```

## Contributing

Contributions are strongly appreciated. Everything from creating bug reports to contributing code will help the project a lot, so please feel free to help in any way, shape, or form that you feel comfortable doing.

## Name

The word "rymdport" comes from the Swedish language and means "space gate".
As a wormhole is a kind of gateway through space, it became the new name after "wormhole-gui".

## License
- Rymdport is licensed under `GNU GENERAL PUBLIC LICENSE Version 3`.
