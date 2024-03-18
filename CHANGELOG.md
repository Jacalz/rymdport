# Changelog

## 3.5.3 - Dependency updates and mindor documentation improvements
- Updated `fyne` to [v2.4.4](https://github.com/fyne-io/fyne/releases/tag/v2.4.4).
  - This brings various bug fixes and some sizeable performance improvements.
- Updated `compress` to [v1.17.7](https://github.com/klauspost/compress/releases/tag/v1.17.7).
  - Decompression performance for directory transfers is improved slightly.
- Update `fyneselfupdate` to [v0.1.1](https://github.com/fynelabs/fyneselfupdate/releases/tag/v0.1.1).
  - This removes some usage of deprecated functions and cleans up the code. Nothing major.
- Some improvements to documentation.
  - Some incorrect wording and textual errors were fixed in [CONTRIBUTING.md](CONTRIBUTING.md). 
  - Information on how to build distribution packages (for Linux repositories etc.) was added in a new [PACKAGING.md](PACKAGING.md) file.

## 3.5.2 - Updated Fyne UI toolkit dependency
- Updated `fyne` to [v2.4.3](https://github.com/fyne-io/fyne/releases/tag/v2.4.3).

## 3.5.1 - Bug fixes and improved performance
- The `Makefile` now includes `.PHONY` targets to avoid problems with targets having the same name as local files. 
- The `Makefile` no longer updates the icon cache (per request from Linux package maintainers).
  - A user installing locally should manually run `make update-icon-cache` afterwards instead. 
- Updated `fyne` to [v2.4.1](https://github.com/fyne-io/fyne/releases/tag/v2.4.1).
  - File names inside the file dialog grid view are now shown on two lines and with better truncation.
  - Various other improvements and bug fixes.
- Updated `compress` to [v1.17.0](https://github.com/klauspost/compress/releases/tag/v1.17.0) for more efficient decompression of directory transfers.

## 3.5.0 - Major theme redesign plus drag and drop support
- Major theme rework with rounded corners.
- Added support for dragging and dropping files and folders onto the application (#73).
- Users on Linux and BSD can now open files or folders in Rymdport from the file manager to start a transfer (#91).
- A transfer can now be started by passing files or folders as command-line arguments (#75).
- Pressing enter after writing a custom code will now start the transfer (#55).
- Instead of giving an error, a number is now added to the end of a file if a duplicate exists (#87).
- Fixed long codes flowing into the progress bar (#96).
- Fixed the QR code showing for failed and completed transfers (#90).
- Refocus the main window after sending text (#85).
- Enabling the option to overwrite files now shows the submit button using a warning colour (#48).
- Fixed a memory leak when starting a second text transfer with the window already open.
- The text send window now opens faster for transfers after the initial start.
- Directory transfers now check that the size and file number are as expected before extraction.
- Long filenames and codes are now truncated using ellipsis instead of clipping.
- The `Makefile` for building and installing on Linux and BSD has been improved.
  - Binaries are now more reproducible and cleaned up more aggressively for release builds.
  - The icon path has been updated to install the files in a more suitable location.
  - Icons are now installed in more sizes to look better in more places.
  - The icon cache is now automatically updated after installing files.
- Release bianries now install Rymdport with a capital R.
- Go 1.18 or later is now required for compiling the application. 
- Updated `fyne` to [v2.4.0](https://github.com/fyne-io/fyne/releases/tag/v2.4.0).
  - Added an option to disable animations (which can be found within our Settings tab).
  - Rectangles now use OpenGL hardware-acceleration for improved rendering performance.
  - Various improvements to keyboard navigation in widgets.
  - Fix text in the progress bar flickering and having the wrong colour sometimes.
  - Various bug fixes, performance improvements, and code cleanups.
- Updated `wormhole-william` to [v1.0.7](https://github.com/psanford/wormhole-william/releases/tag/v1.0.7) for minor bug fixes and performance improvements.
- Updated `compress` to [v1.16.7](https://github.com/klauspost/compress/releases/tag/v1.16.7) for slightly more efficient compression and decompression of directory transfers.
- Various cleanups and performance improvements throughout the codebase.

## 3.4.0 - QR code support and backend rewrite
- QR codes can now be generated for easily sending items to [supported apps](https://github.com/Jacalz/rymdport/wiki/Supported-clients#clients-with-qr-code-scanning) (fixes #10).
- Rewrote a large part of the backend for displaying sends and receives.
  - Progress bars no longer display the wrong percentage sometimes.
  - Fixes a potential crash that could happen sometimes when sending.
- Added support for removing items that have completed the transfer (fixes #32).
- Added keyboard shortcuts for switching between tabs. See the [Keyboard Shortcuts](https://github.com/Jacalz/rymdport/wiki/Keyboard-shortcuts) wiki for more information.
- Added an option to disable update checking on startup for release binaries (fixes #66).
- The receive page now shows which code a received item came from.
- Various improvements to performance and memory usage.
- Updated `fyne` to [v2.3.5](https://github.com/fyne-io/fyne/releases/tag/v2.3.5) for many bug fixes.

## 3.3.6 - Crash fix for copying received text
- Fixed a crash when pressing the copy button in the text receive window (fixes #83).
- Fixed the AppStream metadata not containing the v3.3.5 release information.

## 3.3.5 - Performance and memory usage improvements
- Updated `fyne` to [v2.3.4](https://github.com/fyne-io/fyne/releases/tag/v2.3.4).
  - Binary size improvements (roughly 20% smaller binaries).
  - Various performance and memory usage improvements.

## 3.3.4 - Add missing AppStream metadata for v3.3.3
- Fixed the AppStream metadata not containing the v3.3.3 release information.

## 3.3.3 - Better Apple M2 support and improved rendering 
- Updated `selfupdate` to [v0.2.0](https://github.com/fynelabs/selfupdate/releases/tag/v0.2.0).
  - This fixes update notifications showing when there were no new releases (issue #76). 
- Updated `fyne` to [v2.3.3](https://github.com/fyne-io/fyne/releases/tag/v2.3.3). 
  - Many improvements when running on Apple M2 devices.
  - Performance and memory improvements for text rendering.
  - Faster performance when resizing the window on Linux/BSD.
  - Fixed letters being cropped in some cases.
  - Other minor fixes and improvements.
- Updated indirect dependencies for some minor security fixes. 

## 3.3.2 - Performance improvements and bug fixes
- Minor performance improvements to the list of sent/received items.
- Updated `fyne` to [v2.3.1](https://github.com/fyne-io/fyne/releases/tag/v2.3.1).
  - Fixes for various potential crashes (most importantly when minimizing on Windows sometimes).
  - Support being displayed over VNC on RaspbianOS (for Raspberry Pi computers).
  - Fixes an issue with list items sometimes appearing hovered when scrolling.

## 3.3.1 - Fix Windows builds not starting
- Fixes a cross-compilation issue that caused the Windows binary to not work correctly.

## 3.3.0 - New theme and improved usability
- The application can now automatically update itself.
- Notifications for received files now mentions where the file was saved to.
- Present more useful information about how to write custom codes.
- The dialog for custom codes now focuses the text entry automatically.
- The selection of the download path has been reworked to be more useful.
- Minor redesign of the send and receive tabs.
- Minor wording improvements inside the settings tab.
- Updated `fyne` to [v2.3.0](https://github.com/fyne-io/fyne/releases/tag/v2.3.0).
  - An entirely new theme with improved visuals.
  - The application will now follow the FreeDesktop Dark Style Preference on Linux/BSD.
  - Added an option to create a new folder within the folder selection dialog.
  - Lots of other improvements and fixes.
- Update `compress` to [v1.15.14](https://github.com/klauspost/compress/releases/tag/v1.15.14).
  - Includes performance improvements that benefit sending and receiving of directories.
- Minimal supported Go compiler version is now Go 1.17.
- Release binaries are now built using Go 1.19.
  - Various performance improvements, security fixes and other improvements.

## 3.2.0 - Improvements and bug fixes
- Work around send and receive windows not being focused correctly sometimes.
- Go backwards in tab completion using Shift + Tab.
- Update `compress` to [v1.15.8](https://github.com/klauspost/compress/releases/tag/v1.15.8) for CVE-2022-30631 fix.
- Clicking the icon on the about page now opens the repository in a browser.
- The `zip` and `completion` packages are now exported and considered stable.
- Updated `fyne` to [v2.2.3](https://github.com/fyne-io/fyne/releases/tag/v2.2.3) for minor bug fixes.
- Minor performance improvements.
- The `Makefile` for Linux and BSD release binaries now includes an option to install into the home directory.
- Release binaries are now built using Go 1.18.
  - Performance of arm64 binaries should improve by around 10-20%.
  - Various security fixes and other improvements.

## 3.1.0 - Tab completion and more BSD support
- Added support for tab completion of receive-codes (fixes #35).
- Initial support for OpenBSD and NetBSD.
  - Feel free to report any feedback. This support is new and experimental.
- Notifications on Linux and BSD now show the application icon.
- The list view in the file chooser now display extensions correctly (fixes #39).
- The entry for receive-codes now becomes unfocused when pressing escape.
- Updated `fyne` to [v2.2.1](https://github.com/fyne-io/fyne/releases/tag/v2.2.1).
  - Better error reporting on Windows when OpenGL is not available.
  - Many optimizations and widget performance enhancements.
- Updated `compress` to [v1.15.5](https://github.com/klauspost/compress/releases/tag/v1.15.5).
  - Includes minor performance improvements that benefit sending and receiving of directories.
- Go 1.16 is now the oldest supported compiler. Support for older versions has been removed.
- Various improvements and fixes for the AppStream metadata.
- Notifications are now enabled by default.

## 3.0.2 - Improvements to AppStream Metadata
- Added release summaries and removed markdown leftovers from AppStream metadata.
  - This should mostly benfit the Flatpak package and Linux packages.

## 3.0.1 - Flatpak support and performance improvement
- Optimized text receives to be faster and use significantly less memory.
  - This change also means that saving the text to a file will be faster.
- The Windows release binary no longer opens a terminal on start-up (#34).
- Various improvements and fixes for the AppStream metadata.
  - These changes should make it possible to create a Flatpak (for #23)

## 3.0.0 - Rymdport is the new wormhole-gui
- Added support for sending using custom codes.
- UI scaling and primary color can now be changed in the settings tab.
- Dialogs are not scaled to window size (fixes #16).
- Progress is now shown for receives as well (fixes #20).
- Added support for verifying sends and receives before accepting them (fixes #18).
- Improved application startup time by optimizing how settings are handled on startup.
- Long codes are now truncated to avoid moving other UI elements.
- The user now has to confirm before enabling overwriting of files.
- Fixed an issue where sending with received text open would remove the text.
- Removed support for removing completed sends and receives (see #32).
  - This has been broken broken a long time. Will be introduced again in a later version.
- The filename when saving received text now also contains the current time. 
- Many improvements to the contents of the appstream metadata.
- Various minor performance improvements and race condition fixes.
- Updated `wormhole-william` to [v1.0.6](https://github.com/psanford/wormhole-william/releases/tag/v1.0.6).
  - Fixes compatability with the [magic-wormhole.rs](https://github.com/magic-wormhole/magic-wormhole.rs) client.
  - Switched to a faster websocket library.
- Updated `fyne` to [v2.1.3](https://github.com/fyne-io/fyne/releases/tag/v2.1.3).
  - Improves performance, fixes a few memory leaks and minor visual refresh among many other improvements.
- Updated `compress` to [v1.15.0](https://github.com/klauspost/compress/releases/tag/v1.15.0).
  - Includes various performance improvements that benefit sending and receiving of directories.
- Release binaries are now built using Go 1.17.
  - Performance of amd64 binaries should improve by around 5-10%.
  - Lowest supported macOS release is now 10.13 High Sierra.
  - Includes various other fixes and improvements.
- Release binaires for FreeBSD and Linux are now `xz` compressed to decrease sizes.
- Release binaries on macOS now contain the correct version and build number metadata.

## 2.3.1 - Rebuilt release binaries
- Updated `compress` to [v1.13.6](https://github.com/klauspost/compress/releases/tag/v1.13.6).

## 2.3.0 - FreeBSD and macOS arm64 binaries
- Added support for receiving from custom codes (sending will be in the next big release)
- Fixed received data not showing until after download completes, #17.
- Fixed a possible incorrect error that could happen when a text receive failed.
- Fixed an issue with the project module structure that made it impossible to download using `go get` or `go install`.
- Slightly faster application startup time.
- Fix issue with send items sometimes not being unselected correctly.
- Very minor performance improvement for receives.
- Avoid hardcoded defaults for advanced settings. We now use the defaults from `ẁormhole-william` directly instead.
- Release binaries are now available for FreeBSD and macOS (M1) on the arm64 architecture.
- Release binaries for macOS are now called that instead of `darwin` for clarity.
- Release binaries are now built with `Go 1.16.7`.
  - Fixes a couple security issues and contains a few bug fixes. 
- Updated `compress` to [v1.13.3](https://github.com/klauspost/compress/releases/tag/v1.13.3).
  - Better and faster zip compression and decompression (brings faster directory sends and receives).
- Updated `fyne` to [v2.0.4](https://github.com/fyne-io/fyne/releases/tag/v2.0.4).
  - The title bar on Windows 10 now matches the system theme (light or dark theme).
  - Fixed the Windows 10 notifications view showing the text "app-id" as application name.
  - Fixed a couple issues when running in fullscreen.
  - Improved performance when drawing transparent rectangles or whitespace strings.  

## 2.2.2 - A small hotfix release
- Fixed the receive code validation being too strict in some cases.
- Fixed incorrect version information on the about tab.

## 2.2.1 - A few minor bug fixes
- Fixed text wrapping being disabled for the text send/receive windows.
- Updated `fyne` to [v2.0.3](https://github.com/fyne-io/fyne/releases/tag/v2.0.3).
  - Fixed compilation on FreeBSD 13.
  - Fixed an issue when clicking links on macOS.
  - Improvements and fixes for text selection.
  - A few minor performance improvements.

## 2.2.0 - Much faster directory transfers
- Replace [mholt/archiver](https://github.com/mholt/archiver) with a custom zip extractor using [klauspost/compress/zip](https://github.com/klauspost/compress).
  - Binaries are about 0.5 MB smaller due to not including unused compression standards.
  - Improved performance when receiving directories.
- Added advanced settings for the wormhole client.
  - Support for changing the default AppID, Rendezvous Server URL and Transit Relay Address used for transfers.
- Improved error handling for receives.
- Improved memory usage when when receiving text.
- Fixed files not being closed if send failed to start.
- Fixed main window being unresponsive when sending text.
- Fixed sent/received text staying in memory until the next send/receive.
- Updated `wormhole-william` to [v1.0.5](https://github.com/psanford/wormhole-william/releases/tag/refs%2Ftags%2Fv1.0.5).
  - Switched to [klauspost/compress/zip](https://github.com/klauspost/compress) for up to 2.5x faster zip compression when sending directories.
- Updated `fyne` to [v2.0.1](https://github.com/fyne-io/fyne/releases/tag/v2.0.1).
  - Improved refresh and resizing of dialogs.
  - Initial support for building on Apple M1 computers (arm64).
  - Fixed some buttons not showing hover effects.
  - Fixed progress bars not having correct background.
  - Fixed pointer and cursor misalignment when typing text.
  - Fixed possible panic when selecting text.
  - Fixed cursor animation sometimes distorting the text.
- Release binaries are now built with `Go 1.16.2`.
  - Binaries for macOS no longer support 10.11 El Capitan and instead require 10.12 Sierra or later.
  - Windows binaries are built with [ASLR](https://en.wikipedia.org/wiki/Address_space_layout_randomization) support for improved security.
  - All binaries are now smaller due to using an improved linker.
  - Small performance improvements and other minor changes.

## 2.1.0 - Major ui changes and a lot of improvements
- Major rework of the receive tab to use progress bars.
  - Dialogs are no longer used to indicate finished receives.
- Redesigned settings tab to use a more modern layout.
  - The component slider setting now displays the currently selected number.
  - Fixed a bug where the component length slider did not have distinct steps.
  - Added a setting to allow existing files to overwritten (disabled by default).
- Multiple improvements to the text send and receive windows.
  - Opens faster by only being created once instead of on each send/receive.
  - Clicking `CTRL + SHIFT` in the send window now starts the send.
  - Tab characters are now displayed correctly (bug fix in fyne).
  - Buttons now use better wording and better icons.
- Some small performance improvements for send and receive of files and directories.
- Receives are now properly rejected instead of just not being downloaded.
- Better notification handling by indicating success and fail for both sends and receives.
- Copying the code of a sent item is now slightly faster.
- Fixed an issue where existing files could be overwritten.
- Fixed a couple possible race conditions on sending data.
- Fixed an issue that prevented enter on the numpad from starting the receive.
- Fixed an issue where file extensions would be displayed as `.` when waiting for data.
- Fixed a bug that caused the window to not be able to shrink to the correct smallest size.
- Multiple other code cleanups, restructurings and minor fixes.
- Added an [appstream metadata](https://www.freedesktop.org/software/appstream/docs/) file for Linux and BSD systems (installed via `make install`).
- Updated `fyne` to [v2.0.0](https://github.com/fyne-io/fyne/releases/tag/v2.0.0).
  - The tabs are now animated to be more responsible on change.
  - Buttons now show an animation on tapped.
  - Theme changes and other improvements to styling.
  - Multiple smaller performance optimizations for widgets.
  - Improved scaling on HIDPI displays.
- Release binaries are now built with `Go 1.14.15`.
  - A security fix for `crypto/elliptic` and a few smaller bug fixes.

## 2.0.2 - A few fixes while waiting for 2.1.0
- Fix a bug that prevented folder send to work on Windows.
- Fix an issue where the ui would become unresponsive on dismissing text send.
- Fixed a bug where the application could crash on typing an incorrect code.
- Sending large files, folders or text will no longer slow down the ui.
- Updated `fyne` to [v1.4.3](https://github.com/fyne-io/fyne/releases/tag/v1.4.2).
  - Fix an issue with notifications sometimes not showing on MacOS.
- Release binaries are now built with `Go 1.14.13`.
  - Improved performance thanks to multiple runtime improvements.
  - Windows binaries now have [DEP (Data Execution Prevention)](https://docs.microsoft.com/en-us/windows/win32/memory/data-execution-prevention) enabled.

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
