---
title: "OpenClaw on a Machine That Never Saw It"
date: 2026-03-19T16:15:00-04:00
draft: false
---

[OpenClaw](https://docs.openclaw.ai/) is a **gateway** I run on hardware I control. It connects an assistant to **Slack**, **cron**, and whatever plugins I add, and exposes the same integrations to **Cursor** over MCP. I wanted one setup I could repeat without guessing which decade-old SSH key or half-forgotten API client was interfering.

Here’s what I did: I used a **fresh laptop** (you could use a VM or a new user account—the point is isolation), stood up GitHub and Slack identities the gateway would need, installed OpenClaw from the docs, pointed Cursor at the gateway, and got to a first session where I could ask the assistant to **edit its own** `AGENTS.md` and log what we’d done. The rest of this post is that path in order, plus the sharp edges I actually hit.

## Why I didn’t start on my daily driver

I *could* have installed OpenClaw on the machine I’ve carried for years. That machine already had a dozen API tokens in shell history, three GitHub identities worth of SSH config, Slack in the browser, corporate VPN, and tools I don’t fully remember installing. An assistant with MCP and a workspace full of scripts isn’t malicious by default, but it **amplifies** whatever access I give it.

On a clean machine I get a few things for free:

* **Blast radius.** If I mis-scope a Slack token or paste something dumb into `openclaw.json`, I’m not risking the rest of my life’s toolchain at the same time.
* **Debuggability.** When something fails, I’m not wondering whether it’s OpenClaw or the Node version I installed in 2019 for one contract.
* **A checklist I can trust.** “Works on the fresh install” is a reproducible story I can turn into a blog post or hand to future-me.

I’m not saying you must buy hardware. I’m saying I **separated** “prove the stack” from “merge it into my main environment,” and that made the first day calmer.

## Baseline: OS, Git, editor

I’m on **Linux** with **systemd user sessions** so the gateway can run as `openclaw-gateway` under my user without sudo drama. If you’re on macOS, the same ideas apply; the exact service commands differ—follow the current OpenClaw install guide rather than copying my `systemctl` lines blindly.

I installed **Git**, **Cursor**, and whatever runtime the OpenClaw docs asked for that week. I did **not** wire Cursor to anything yet. The gateway has to exist first; otherwise I’m debugging MCP and Slack at the same time, which is a bad hobby.

## GitHub: me and the bot

I use my normal GitHub account for *me*. For anything the assistant does as an automated identity—forking a repo, opening a PR—I use a **separate account** with its own SSH key. That way I never wonder which key `git@github.com` is going to offer, and I can revoke bot access without touching my personal keys.

I generated an ed25519 key for the bot, added the public key to the bot account via GitHub’s API, and gave **`~/.ssh/config`** a host alias so only bot remotes use that key:

```text
Host github.com-mybot
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_mybot
    IdentitiesOnly yes
```

Remotes for the bot look like `git@github.com-mybot:mybot/somerepo.git`. When I need the REST API (fork, PR), I use a **fine-grained PAT**—but for `git push` I use **SSH** so tokens never end up embedded in `git remote` URLs or shell transcripts. I learned that one the embarrassing way.

## Slack: the app and the token

I created a Slack app in the API portal, added the **bot token scopes** I actually need (posting, reading DMs or channels—your list will match how you want the bot to behave), installed the app to my workspace, and grabbed the **bot token** and any workspace metadata OpenClaw expects.

I started with **DMs** to the bot before I invited it to busy channels. Fewer moving parts: no “which channel ID,” no threading surprises, just “did a message get through.”

## Installing the gateway

I followed [OpenClaw’s documentation](https://docs.openclaw.ai/) for install and first run. On my machine everything authoritative lives under **`~/.openclaw/`**: `openclaw.json`, cron jobs, lock files, that sort of thing.

The **workspace**—where `AGENTS.md`, `SOUL.md`, `USER.md`, `memory/`, and `scripts/` live—is a **separate directory** whose path I set in config (e.g. `agents.defaults.workspace`). I treat that folder as the assistant’s home; the dotdir is the daemon’s home.

Once the **user service** was enabled, I checked it like any other service I own:

```bash
systemctl --user status openclaw-gateway
journalctl --user -u openclaw-gateway -n 50 --no-pager
```

After a bad config edit:

```bash
systemctl --user restart openclaw-gateway
```

I also tail the same-day file log when I don’t want to swim in journald:

```bash
tail -f "/tmp/openclaw/openclaw-$(date +%F).log"
```

I didn’t touch MCP until this stayed up across a reboot.

## The workspace is the memory

The model forgets everything between sessions. **Files don’t.** So I keep rules in `AGENTS.md`, persona in `SOUL.md`, facts about me in `USER.md`, and a running log under `memory/YYYY-MM-DD.md`. Longer-lived curated notes go in `MEMORY.md` when I want them loaded only in private sessions—that’s a policy choice I documented in `AGENTS.md`, not a framework mystery.

The first thing I asked in Cursor after the workspace existed was: read `AGENTS.md`, `SOUL.md`, `USER.md`, and today’s memory file, then summarize what constraints you’ll follow. If it can’t do that, the rest of the project is theater.

## Cursor and MCP

Cursor reads MCP servers from **`~/.cursor/mcp.json`**. I added an entry pointing at my gateway URL and the token `openclaw.json` expects for that connection. I did not paste secrets into the blog draft; I used placeholders in comments and real values only on disk.

After saving, I **reloaded MCP** in Cursor (sometimes I restart Cursor—same outcome, more coffee). I ran **`openclaw_discover`**. If the tool list looked empty or tiny, I **restarted the gateway** and reloaded again—sync lag is a real failure mode, not imagination.

For anything non-obvious I use **`openclaw_skill`** for docs and **`openclaw_invoke`** for tools that aren’t first-class in Cursor yet.

## Wiring Slack

I put the Slack credentials where the gateway expects them (again: follow the current doc layout; the shape is “config + secret,” not “magic env var I forgot”).

One invariant I care about: **one reply per inbound message.** I have multiple Cursor windows sometimes; the gateway doesn’t know that. I keep a small script in the workspace that acquires a lock under `~/.openclaw/locks/` before posting, so two sessions don’t double-text a human. The rule lives in `AGENTS.md` too, because culture and code both matter.

## The first loop that felt real

Once MCP and Slack were alive, I asked the assistant to **edit its own house**: tighten `AGENTS.md` for Slack etiquette, adjust `SOUL.md` for tone, and append what we’d done to today’s `memory/` file. That’s the moment the architecture stopped being theoretical—**behavior as code in git**, not vibes in a chat buffer.

## Sharp edges I actually hit

**Cron.** Jobs live in `~/.openclaw/cron/jobs.json`. Each schedule needs a **`schedule`** object with **`schedule.kind`** (`cron`, `at`, or `every`). A bare top-level `"cron": "0 9 * * *"` string blows up on current gateway versions with errors about missing `kind`. The [cron jobs doc](https://docs.openclaw.ai/automation/cron-jobs) is the source of truth.

**Jira and `acli`.** Atlassian’s [**acli**](https://developer.atlassian.com/cloud/acli/) doesn’t pick up `JIRA_USERNAME` + `JIRA_API_TOKEN` from the environment the way raw `curl` does. I logged in once with the token on stdin:

```bash
printf %s "$JIRA_TOKEN" | acli jira auth login \
  --site your-site.atlassian.net \
  --email 'you@company.com' \
  --token
```

Then `acli jira auth status` and normal commands work.

For ad-hoc REST I still use **`curl`** with basic auth. Atlassian has been retiring **`GET /rest/api/3/search`**; I use **`POST /rest/api/3/search/jql`** with a JSON body now. Details are in [their changelog](https://developer.atlassian.com/changelog/).

```bash
curl -s -u "${JIRA_USER}:${JIRA_TOKEN}" \
  -H 'Accept: application/json' -H 'Content-Type: application/json' \
  -X POST 'https://your-site.atlassian.net/rest/api/3/search/jql' \
  -d '{"jql":"project = MYPROJ ORDER BY updated DESC","maxResults":5,"fields":["summary","status"]}'
```

**Skills.** Anything under `.agents/skills/.../SKILL.md` is a playbook Cursor can load; the gateway can expose more tools depending on what I install. Community skills exist; I read them like any third-party dependency before I put tokens near them.

## When I called it “working”

I had: a gateway that survives reboot, Cursor talking to it over MCP with a non-empty tool list, Slack DMs round-tripping without duplicate spam, and a workspace under git where the assistant could change its own rules and log the day.

Everything else—more channels, Jira epics, Confluence—is just more tokens and more markdown. The skeleton was the hard part; the fresh machine is why I could **see** the skeleton.
