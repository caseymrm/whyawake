## v0.7 — "Until I close the lid"

New sleep-prevention option: **Until I close the lid**. Pick it from the menu, get on with your work, and the moment you close the lid the app stops caffeinating and lets the system sleep normally — no timer to guess at, no manual deactivation.

### How it works
A new sister library, [`caseymrm/go-clamshell`](https://github.com/caseymrm/go-clamshell), subscribes to IOKit's `kIOPMMessageClamshellStateChange` push notification (no polling — the kernel tells us the moment the lid moves). When the state flips to closed *and* you're in lid mode, whyawake calls `caffeinate.Stop()` and the deferred sleep proceeds.

### Other
- Picks up `go-clamshell` v1.0.1.
- Closes #1.

### Gatekeeper warning
This build is still **unsigned** (ad-hoc only). To open:

1. Right-click `WhyAwake.app` → **Open** → confirm in the dialog, **or**
2. Run `xattr -dr com.apple.quarantine /Applications/WhyAwake.app` after copying to `/Applications`.

### Install
- Download `WhyAwake.app.zip` below.
- Unzip, drag to `/Applications`.
- Launch (see Gatekeeper note above).
