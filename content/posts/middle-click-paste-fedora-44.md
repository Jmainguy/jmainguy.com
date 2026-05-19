---
title: "Middle-Click Paste Broke After Upgrading to Fedora 44"
date: 2026-05-18T10:00:00-04:00
draft: false
---

I upgraded to Fedora 44 a few days ago and immediately noticed something was wrong in GNOME Terminal. Highlighting text no longer felt like it was copying anything, and middle-click paste had stopped working entirely. I have been using Linux on the desktop for over a decade and this workflow is muscle memory at this point: select text in one terminal, middle-click in another, done. No Ctrl+C, no Ctrl+V.

It turns out nothing was broken in GNOME Terminal itself. GNOME changed the default.

## Two clipboards on Linux

On X11 (and still on Wayland, when enabled), Linux traditionally has two separate clipboards:

- **CLIPBOARD** — what you get with explicit Copy/Paste (Ctrl+C / Ctrl+V). This behaves like Windows and macOS.
- **PRIMARY** — updated when you *select* text. You paste it with the middle mouse button, without ever pressing Copy.

The Freedesktop clipboard specification describes PRIMARY as something expert users use; regular users can ignore it. That framing matters, because GNOME's recent change is essentially "we agree, so we turned it off for everyone by default."

- [Freedesktop: Clipboards wiki](https://www.freedesktop.org/wiki/Specifications/ClipboardsWiki/)

## What Fedora 44 / GNOME 50 changed

The setting lives in GSettings as `org.gnome.desktop.interface gtk-enable-primary-paste`. When it is `true`, GTK applications honor primary selection for middle-click paste. When it is `false`, middle-click does nothing useful in those apps.

The default flipped from `true` to `false` in [gsettings-desktop-schemas merge request !119](https://gitlab.gnome.org/GNOME/gsettings-desktop-schemas/-/merge_requests/119), opened by Jordan Petridis in January 2026. His merge request description calls primary paste an "X11ism" tied to the old Gtk/EnablePrimaryPaste XSetting ([GNOME bug 775844](https://bugzilla.gnome.org/show_bug.cgi?id=775844)), and argues that accidental middle-clicks dumping clipboard contents without warning is bad UX. The MR closes with "Goodbye X11."

The schema key itself is defined here (note `<default>false</default>`):

```xml
<key name="gtk-enable-primary-paste" type="b">
  <default>false</default>
  <summary>Enable the primary paste selection</summary>
  <description>
    If true, gtk+ uses the primary paste selection, usually triggered by a middle mouse button click.
  </description>
</key>
```

Source: [`org.gnome.desktop.interface.gschema.xml`](https://gitlab.gnome.org/GNOME/gsettings-desktop-schemas/-/blob/main/schemas/org.gnome.desktop.interface.gschema.xml) in the [gsettings-desktop-schemas](https://gitlab.gnome.org/GNOME/gsettings-desktop-schemas) repository.

Fedora 44 ships GNOME 50, which picks up that new default. So after upgrade, existing users who never touched the setting suddenly lose middle-click paste — and the "select to copy" workflow *feels* broken because the only paste path they used for it is gone.

## How I fixed it

Re-enable primary paste for your user:

```bash
gsettings set org.gnome.desktop.interface gtk-enable-primary-paste true
```

Verify:

```bash
gsettings get org.gnome.desktop.interface gtk-enable-primary-paste
```

Open a **new** GNOME Terminal tab or window after changing the setting; already-running windows may not pick it up immediately.

You can also toggle this in **GNOME Tweaks** if you prefer a GUI (the option is exposed there on recent GNOME versions).

## This is not a Wayland regression (exactly)

I am running a Wayland session (`echo $XDG_SESSION_TYPE` → `wayland`). Wayland has supported primary selection for years; GNOME documented the initiative [here](https://wiki.gnome.org/Initiatives/Wayland/PrimarySelection). The Fedora Project [Wayland features wiki](https://fedoraproject.org/wiki/Wayland_features) still lists primary selection under historical blockers for older releases.

What changed in Fedora 44 is policy and defaults, not missing Wayland plumbing.

## The discussion is worth reading

People have strong feelings about this, on both sides.

**For disabling it by default:** accidental middle-click paste can leak sensitive text (passwords, tokens) into the wrong window, especially during screen sharing. The Freedesktop wiki's "easter egg for experts" language is cited in coverage such as [It's FOSS: middle-click paste likely disabled in future GNOME](https://itsfoss.com/news/gnome-firefox-middle-click-paste-removal/). Jordan Petridis also posted about the coordinated GNOME/Firefox effort on [Mastodon](https://mastodon.social/@alatiera/115832011216789958).

**Against disabling it by default:** long-time Linux users rely on primary selection constantly. We liked this setting because it made a tiny, repeated workflow feel instant, which is exactly the kind of habit change [xkcd 1172: Workflow](https://xkcd.com/1172/) jokes about. Several LWN readers noted that Fedora 44 shipped without prominent release-note mention, leaving users to discover the breakage on their own — see the [middle-click paste thread](https://lwn.net/Articles/1070435/) under [Fedora Linux 44 has been released](https://lwn.net/Articles/1070198/). Ubuntu users hit the same surprise ([Bug 2145179 discussion](https://www.mail-archive.com/desktop-bugs@lists.ubuntu.com/msg830632.html)).

Mozilla opened a parallel change for Firefox: [Phabricator D277804](https://phabricator.services.mozilla.com/D277804).

## Quick reference

| What you want | What to use |
|---------------|-------------|
| Explicit copy/paste | Ctrl+C / Ctrl+V (CLIPBOARD) |
| Select text, paste with middle button | PRIMARY (requires `gtk-enable-primary-paste true`) |
| Re-enable on GNOME | `gsettings set org.gnome.desktop.interface gtk-enable-primary-paste true` |

## Links

- [GNOME MR !119 — Disable primary-paste by default](https://gitlab.gnome.org/GNOME/gsettings-desktop-schemas/-/merge_requests/119)
- [gsettings-desktop-schemas: `org.gnome.desktop.interface.gschema.xml`](https://gitlab.gnome.org/GNOME/gsettings-desktop-schemas/-/blob/main/schemas/org.gnome.desktop.interface.gschema.xml)
- [GNOME bug 775844 — Gtk/EnablePrimaryPaste gsetting](https://bugzilla.gnome.org/show_bug.cgi?id=775844)
- [Freedesktop: Clipboards wiki](https://www.freedesktop.org/wiki/Specifications/ClipboardsWiki/)
- [GNOME Wayland primary selection initiative](https://wiki.gnome.org/Initiatives/Wayland/PrimarySelection)
- [Fedora Project: Wayland features](https://fedoraproject.org/wiki/Wayland_features)
- [LWN: Fedora Linux 44 released](https://lwn.net/Articles/1070198/)
- [LWN: Middle-click paste discussion](https://lwn.net/Articles/1070435/)
- [It's FOSS on the GNOME/Firefox proposals](https://itsfoss.com/news/gnome-firefox-middle-click-paste-removal/)
- [Mozilla Phabricator D277804 (Firefox)](https://phabricator.services.mozilla.com/D277804)
- [Ubuntu desktop-bugs: middle click paste (Bug 2145179)](https://www.mail-archive.com/desktop-bugs@lists.ubuntu.com/msg830632.html)

If you upgraded to Fedora 44 and your terminal "stopped copying" when you highlight text, try the `gsettings` line above before you spend an hour debugging VTE or wl-clipboard.
