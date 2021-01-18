# Changelog

## 2.0.2 - A few fixes while waiting for 2.1.0
- Fix a bug that prevented folder send to work on Windows.
- Fix an issue where the ui would become unresponsive on dismissing text send.
- Fixed a bug where the application could crash on typing an incorrect code.
- Sending large files, folders or text will no longer slow down the ui.
- Updated `fyne` to [v1.4.3](https://github.com/fyne-io/fyne/releases/tag/v1.4.2).
  - Fix an issue where with notifications sometimes not showing on MacOS.
- Release binaries are built with `Go 1.14.15` for faster performance.

## 2.0.1 - Minor fixes and FreeBSD release binaries
- Binaries for `freeBSD/amd64` are now available on the release page.
- Corrected the icon for the receive tab.
- Updated `fyne` to [v1.4.2](https://github.com/fyne-io/fyne/releases/tag/v1.4.2).
  - Dialog shadow does not resize correctly sometimes.
  - Possible crash when minimising app on Windows.
  - File chooser ignores drive Z on Windows.

## 2.0.0 - Code rework and many new features
- Massive rework and rewrite of code to simplify and make it more maintainable.
- Use new list widget with custom layout for showing sends and receives.
- Show a button for copying the send code. See #3 for more information.
- Show an icon for each file, folder, or text snippet that is sent.
  - Includes MIME type and extension information.
- Big UI refresh thanks to new theme rework in `fyne v1.4.0`.
- Added support for sending and receiving folders.
- Sending files now properly closes them afterwards.
- Fixed a bug that caused component length to not be saved between application restarts.
- Switch to adaptive theme by default.
  - Will changes depending on dark/light mode on `windows 10` and `macOS`.
- Added folder picker for selecting a downloads directory.
- Multiple performance and memory improvements.
  - Sends and receives are reusing the same `wormhole` client instead of creating a new one each time.
  - Dialogs are now created once and then shown when appropriate (not when showing errors).
  - Themes are no longer checked too many times on startup.
  - Using less goroutines and channels internally.
- The `Makefile` now supports uninstalling too (for Linux and BSD).
- Release binaries are now built for `linux/arm64` as well.
- Add initial build and package support for BSD.
  - The next release will have binaries for `freeBSD/amd64`.
- Updated `fyne` to [v1.4.0](https://github.com/fyne-io/fyne/releases/tag/v1.4.0).

## 1.3.0 - Code refractoring, new features and fixes
- Refactored code to simplify and be more maintainable.
- Added support for sending notifications on send and receive.
  - Can be turned on in settings.
- Added an about page with logo and version number.
- Added build scripts and `.desktop` file for Linux packaging.
- Added a new fancy way of displaying text files.
  - Support for saving text to a file on receive.
  - Support for copying all text to clipboard on receive.
- Make sure that sending text updates progress too.
- Changed arrow down icon to download icon.
- Make file saves more reliable.
- Release binaries are built using `fyne-cross` v2.2.0.
  - Now built with `Go 1.13.15`.
- Updated `wormhole-william` to v1.0.4.
- Updated `fyne` from v1.3.1 to v1.3.3.
  - Brings a bunch of bugfixes and favourite icons in file picker.

## 1.2.0 - Application icon, fixes and new features
- Fixed text transfer between devices.
- Only set the max value for progressbars once.
- Add support for receiving on pressing enter/return.
- Added an application icon based on an actual wormhole.
- Moved out custom widget code to it's own package.
  - Added code copy popup menu on right click to send codes.

## 1.1.0 - New features and fixes
- Added progression bars when sending files.
- Added filename and status information to the receive page.
- Made headers bold for information on the send and receive tabs.
- Fixed an issue with the EventQueue filling up due to blocking calls.

## 1.0.0 - Initial version
The first release of `wormhole-gui`.
