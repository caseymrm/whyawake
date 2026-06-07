## v0.8 — Default left-click action

Pick a default duration from the new **Left-click** submenu, and from then on a single left-click on the menubar icon toggles caffeinate at that duration — no menu, no submenu. Right-click (or Ctrl-click) still opens the full menu.

### What changed
- New **"Left-click: …"** menu section. Choices: Off · Until I close the lid · Indefinitely · 10 / 30 / 60 / 180 min. Stored in NSUserDefaults; the current choice is shown in the parent label and persists across launches.
- When a default is set, left-click *toggles*: if the Mac is being kept awake, the click cancels it; otherwise it starts a new session at the chosen duration. (Accidentally caffeinated → one click to undo.)
- Picks up `menuet` v2.1.1 (new `Application.Clicked` field, plus the v2 Go module path fix).
- Closes #2.

### Gatekeeper
Still **unsigned** (ad-hoc only).
1. Right-click `WhyAwake.app` → **Open** → confirm, **or**
2. `xattr -dr com.apple.quarantine /Applications/WhyAwake.app` after copying to `/Applications`.

### Install
Download `WhyAwake.app.zip` → unzip → drag to `/Applications` → launch (see Gatekeeper note).
