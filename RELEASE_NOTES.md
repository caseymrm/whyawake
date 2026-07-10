## v0.10 — signed & notarized

First release signed with a Developer ID certificate and notarized by
Apple. No code changes beyond the menuet upgrade — the headline is
distribution quality.

### What changed
- **Signed + notarized**: hardened runtime, secure timestamp, stapled
  notarization ticket. Downloads open without the "unidentified
  developer" warning — no right-click-Open, no `xattr` workaround.
- **Notifications now work**: `UNUserNotificationCenter` requires a
  Developer ID signature; previous ad-hoc builds silently dropped
  notifications.
- Picks up [menuet v2.10.4](https://github.com/caseymrm/menuet/releases/tag/v2.10.4)
  (snapshot mode for the [menuet.app](https://menuet.app/apps/whyawake/)
  showcase, hardened-runtime signing in menuet.mk).

### Gatekeeper
Just download, unzip, and open. The quarantine instructions from
earlier releases are no longer needed.
