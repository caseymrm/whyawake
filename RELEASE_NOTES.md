## v0.9 — menuet v2.2.0

Maintenance bump. No user-visible changes — the app behaves the same as v0.8.

### What changed
- Picks up [menuet v2.2.0](https://github.com/caseymrm/menuet/releases/tag/v2.2.0), where `MenuItem` becomes a Go interface with `Regular` and `Separator` as the first two concrete implementations. All internal call sites migrated:
  - `menuet.MenuItem{Text: …}` → `menuet.Regular{Text: …}`
  - `menuet.MenuItem{Type: menuet.Separator}` → `menuet.Separator{}`
- Slice element types in `[]menuet.MenuItem` literals are now qualified, as the new interface requires.

### Gatekeeper
Still **unsigned** (ad-hoc only).
1. Right-click `WhyAwake.app` → **Open** → confirm, **or**
2. `xattr -dr com.apple.quarantine /Applications/WhyAwake.app` after copying to `/Applications`.
