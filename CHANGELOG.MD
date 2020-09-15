# Changelog

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