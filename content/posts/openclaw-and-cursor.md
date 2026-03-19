---
title: "Hooking Up OpenClaw with Cursor"
date: 2026-03-19T16:15:00-04:00
draft: false
---

[OpenClaw](https://docs.openclaw.ai/) is a **gateway** that sits on your machine (or a host you control), wires your AI assistant to **Slack**, **cron**, plugins, and whatever skills you add, and exposes those integrations to clients like **Cursor** over MCP. I spent a day getting the pieces aligned; this post is the condensed version so future-me (and you) can reproduce it without retracing every dead end.

## What you end up with

* A **user systemd service** (`openclaw-gateway`) that keeps the gateway running.
* A **workspace** on disk (markdown + scripts + memory) that the agent reads every session — personality, rules, and continuity live there, not in the model weights.
* **Cursor** talking to the gateway via **MCP**, so the IDE can trigger the same tools the gateway uses for Slack and automation.
* Optional: **skills** (Jira, Confluence, custom scripts) dropped into that workspace.

## Install the gateway and workspace

Follow the current docs for your OS; the important part is that you know where config lives. On my Linux laptop everything under **`~/.openclaw/`** is the source of truth: `openclaw.json`, cron definitions, lock files, etc.

The agent’s **workspace** path is whatever you set in config (for example `agents.defaults.workspace` pointing at a folder that holds `AGENTS.md`, `SOUL.md`, `USER.md`, `memory/`, and `scripts/`). That directory is *home* for the assistant — rules in `AGENTS.md` apply to every turn.

Quick sanity checks once the service is up:

```bash
systemctl --user status openclaw-gateway
journalctl --user -u openclaw-gateway -n 50 --no-pager
```

If the gateway is wedged after a config change:

```bash
systemctl --user restart openclaw-gateway
```

Same-day file logging (handy when you do not want to page through journald):

```bash
tail -f "/tmp/openclaw/openclaw-$(date +%F).log"
```

## Point Cursor at the gateway (MCP)

Cursor reads MCP servers from your user config (on my machine, `~/.cursor/mcp.json`). The OpenClaw entry should target your running gateway URL and use whatever **token** your `openclaw.json` expects for that connection.

After editing MCP config, **reload MCP** in Cursor (or restart Cursor). If `openclaw_discover` only shows a minimal tool set, the gateway may still be syncing plugins — a gateway restart and another reload usually fixes it.

In the chat, prefer the documented flow: **`openclaw_discover`** → **`openclaw_skill`** for full parameter docs → **`openclaw_invoke`** for anything not wrapped as a first-class tool.

## Workspace files that matter

These are not ceremonial; the agent is instructed to read them at session start:

* **`AGENTS.md`** — Operating rules (Slack behavior, memory, red lines, how to use OpenClaw on this machine).
* **`SOUL.md` / `USER.md`** — Who the assistant is and who you are; keeps tone and context consistent.
* **`memory/YYYY-MM-DD.md`** — Daily log; **`MEMORY.md`** — longer-lived curated notes (only for private / main sessions if you split that way).

That layout means you can version-control **behavior** the same way you version code.

## Cron jobs: mind the schema

OpenClaw’s cron store is JSON under `~/.openclaw/cron/jobs.json`. Each job needs a proper **`schedule`** object with **`schedule.kind`** (`cron`, `at`, or `every`). A bare `"cron": "0 9 * * *"` string at the top level will not fly with current gateway versions — you will see errors in logs about missing `kind`. The [cron jobs doc](https://docs.openclaw.ai/automation/cron-jobs) is the canonical reference.

## Jira from the CLI (`acli`) — auth is not your shell env

I use Atlassian’s [**acli**](https://developer.atlassian.com/cloud/acli/) (CLI for Jira Cloud) for Jira. Important detail: **`acli` does not authenticate via `JIRA_USERNAME` + `JIRA_API_TOKEN` in the environment** the way raw `curl` does. It stores credentials after you log in once.

Non-interactive login with an Atlassian **API token** (site + **account email**):

```bash
printf %s "$JIRA_TOKEN" | acli jira auth login \
  --site your-site.atlassian.net \
  --email 'you@company.com' \
  --token
```

Check status:

```bash
acli jira auth status
```

After that, normal commands work:

```bash
acli jira workitem search --jql 'project = MYPROJ AND assignee = currentUser()' --limit 10 --json
```

For one-off REST calls I still use **`curl`** with basic auth (`email:api_token`). If you hit **HTTP 410** on issue search, Atlassian has been retiring the old **`GET /rest/api/3/search`** endpoint — use **`POST /rest/api/3/search/jql`** with a JSON body instead. Migration notes are in [Atlassian’s changelog](https://developer.atlassian.com/changelog/).

Example:

```bash
curl -s -u "${JIRA_USER}:${JIRA_TOKEN}" \
  -H 'Accept: application/json' -H 'Content-Type: application/json' \
  -X POST 'https://your-site.atlassian.net/rest/api/3/search/jql' \
  -d '{"jql":"project = MYPROJ ORDER BY updated DESC","maxResults":5,"fields":["summary","status"]}'
```

## Slack and duplicate replies

If multiple Cursor windows can see the same inbound Slack message, you need a **single-flight** guard before posting. I keep a small script in the workspace that takes the channel, thread, text, and uses a lock under `~/.openclaw/locks/` so only one instance answers. Whether you use that script or your own mutex, the invariant is: **one reply per inbound message**.

## Skills

Skills are just documented capabilities in the workspace (see `.agents/skills/.../SKILL.md` in your tree). Cursor can load them when relevant; the gateway can expose **OpenClaw-native** tools for Slack, browser automation, etc., depending on what you install.

There is a growing **community catalog** of OpenClaw skills online; treat them like any third-party dependency — read what they do, where tokens go, and whether they match your employer’s policy before you plug them into a work machine.

## LinkedIn and “company page” automation

Separate from OpenClaw: **there is no blessed one-click skill** that fully manages a **LinkedIn company page** the way official APIs do for Slack. The sanctioned path is **LinkedIn’s developer products** (app registration, approvals, OAuth scopes for organization publishing). Cookie-based or browser-hammering tools exist in the wild; I am not using those for brand-facing accounts.

## Bottom line

OpenClaw + Cursor is **gateway + workspace + MCP**. Get systemd happy, wire MCP once, treat `AGENTS.md` and `memory/` as production config, log in **`acli`** explicitly, and use **`POST /rest/api/3/search/jql`** when the old Jira search endpoint 410s. After that, most of the work is deciding what you trust the agent to do — and writing that down in the repo.
