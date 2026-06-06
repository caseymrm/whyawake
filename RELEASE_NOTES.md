## v0.6 — Apple Silicon universal build

First release supporting Apple Silicon natively. Universal binary (`x86_64` + `arm64`).

### What changed
- Universal `darwin/amd64` + `darwin/arm64` binary via `lipo` (was: `amd64` only).
- Modernized build to Go modules; new `Makefile` handles both arches and packaging.
- Picked up [`menuet` v1.1.0](https://github.com/caseymrm/menuet/releases/tag/v1.1.0), which migrated notifications off the deprecated `NSUserNotificationCenter` to `UNUserNotificationCenter` (closes the long-standing menuet issue #18).
- `LSMinimumSystemVersion` bumped to **macOS 10.14** (the floor required by the `UserNotifications` framework).

### Gatekeeper warning
This build is **unsigned** (ad-hoc only) and **not notarized**. On first launch macOS will refuse to open it. To open:

1. Right-click `WhyAwake.app` → **Open** → confirm in the dialog, **or**
2. Run `xattr -dr com.apple.quarantine /Applications/WhyAwake.app` after copying it to `/Applications`.

A signed + notarized build will come back once the Developer ID identity is restored.

### Install
- Download `WhyAwake.app.zip` below.
- Unzip, drag `WhyAwake.app` to `/Applications`.
- Launch (see Gatekeeper note above).
