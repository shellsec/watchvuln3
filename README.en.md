**语言 / Language:** [中文](README.md) | English

[aiv123.com](https://aiv123.com/) · AI tools directory, 600+ tools in one place

## 🚀 Recommended: [ofox.ai](https://ofox.io/x/aiv123)

> **In short**: One account for the latest GPT / Claude / Gemini and **100+** top models. First top-up gets an extra **$3** credit.

Text, image, video, and embeddings in one place. Caching supported — repeat calls stay cheaper and faster.

[👉 Sign up](https://ofox.io/x/aiv123) · Global dedicated lines · Enterprise SLA · No conversation retention

| ⚡️ Faster & Leaner | 🧠 Models & Modalities | 🛡️ Privacy |
|:---:|:---:|:---:|
| Global lines, enterprise SLA, plus caching | 100+ models · text / image / video / embeddings | No conversation retention |

## ☕ Buy Me a Coke

Open source takes effort — sponsorship is welcome:  
👉 [爱发电 / Afdian](https://ifdian.net/a/shellsec)

---

# WatchVuln — High-Value Vulnerability Collection & Alerts

[![GitHub Release](https://img.shields.io/github/v/release/shellsec/watchvuln3?label=release)](https://github.com/shellsec/watchvuln3/releases)
[![License](https://img.shields.io/github/license/shellsec/watchvuln3)](https://github.com/shellsec/watchvuln3)

**Repository**: [github.com/shellsec/watchvuln3](https://github.com/shellsec/watchvuln3) · **Version**: v3.1.0

> WatchVuln collects high-severity vulnerabilities from AVD, Chaitin, QiAnXin, OSCS, ThreatBook, Seebug, KEV, and more; filters them by policy; and pushes alerts to DingTalk, WeCom, and other channels. It supports multiple same-type pushers, disabling the startup notification, and a local Web board to browse the intelligence DB (including RSS, REST API, and one-click MCP). DingTalk pushes can append a board link; you can one-click copy an analysis prompt and jump to ChatGPT / Gemini / DeepSeek.

## Vulnerability Intelligence Board

Browse the local intelligence DB with search, filters, and sorting. Next to each title you can one-click copy an analysis prompt and jump to **ChatGPT / Gemini / DeepSeek**. The board also exposes a **REST API** and **MCP** (URLs follow the host you used to open the page). Start with `--web-addr 127.0.0.1:8765` or `watchvuln board`. More details are in the **Vulnerability Intelligence Board (Web, no login)** foldout below.

![Vulnerability Intelligence Board: list, filters/sorting, and AI analysis buttons](https://github.com/user-attachments/assets/6c10562e-4238-42cc-befd-052a72b0f7f7)

*Board list page: sorted by disclosure date; ChatGPT / Gemini / DS next to titles for one-click copy and AI analysis*

![Vulnerability Intelligence Board: vulnerability detail modal](https://github.com/user-attachments/assets/f52964bd-b69d-4925-91f2-a8a2947b0e9c)

*Click a table row to view description, tags, remediation advice, and reference links*

## About This Repository

This project is maintained and extended from the original author **[zema1/watchvuln](https://github.com/zema1/watchvuln)**. The core ideas (multi-source high-value vulnerability collection, filtering, and multi-channel push) come from upstream. Thanks to the original author **zema1**.

In v3.0, `shellsec/watchvuln3` mainly delivered the following:

**Bug fixes**

- Fixed incorrect `Referer` / `Origin` request headers for the QiAnXin TI source
- Fixed `GO_SKIP_TLS_CHECK` not taking effect under Docker / CLI
- Fixed the same CVE being reprocessed every cycle after CVE deduplication
- Fixed a concurrent data race when initializing multiple sources
- Unified default `sources` with the docs; boolean env vars now parse standard `true`/`false`

**Feature enhancements**

- Vulnerability Intelligence Board (`--web-addr` / `watchvuln board`: browse the local DB with sort/filter; one-click copy and jump to ChatGPT / Gemini / DeepSeek for analysis; no login)
- Board **RSS feed** (`/feed.xml`, latest 50 pushed vulns; RSS entry next to the page title)
- Board **REST API / MCP**: `/api`, `/mcp`, one-click connect from the page; URLs follow the current request Host so a new IP is picked up automatically
- DingTalk and other push messages can append a **local board link** (`web_public_host` / `web_public_url`; port follows `web_addr` automatically)
- Config file / `--pusher-file` supports **multiple same-type pushers** (e.g. multiple DingTalk groups)
- `-nm` / `--quiet` disables the “init finished” push on startup
- CLI: `list-sources`, `init-config`, and startup config self-check

See [CHANGELOG.md](./CHANGELOG.md) for the full changelog. If upstream merges equivalent fixes, follow [zema1/watchvuln](https://github.com/zema1/watchvuln); this repo will keep syncing where needed.

---

As everyone knows, over 99% of entries in CVE databases are just numbers with little practical value. I wanted to focus on high-value vulnerabilities that actually need attention right now — instead of drowning in RSS feeds and WeChat public-account ~~threat intel~~. So I wrote this small project to pull from a few high-quality vulnerability sources and push alerts. `WatchVuln` means **watching** for vulnerability updates, and also that these vulns deserve **attention**.

Currently the following sites are scraped:

| Name | URL | Push policy |
|----------------------------|-------------------------------------------------------------------------------------------------|--------------------------------------------------|
| Alibaba Cloud Vuln DB (AVD) | https://avd.aliyun.com/high-risk/list | Severity High or Critical |
| Chaitin Vuln DB | https://stack.chaitin.com/vuldb/index | Severity High or Critical **and** title contains Chinese |
| OSCS Open Source Security Intelligence | https://www.oscs1024.com/cm | Severity High or Critical **and** has the `预警` (alert) tag |
| QiAnXin Threat Intelligence Center | https://ti.qianxin.com/ | Severity High/Critical **and** has one of: `奇安信CERT验证`, `POC公开`, `技术细节公布` |
| ThreatBook Research Response Center (WeChat) | https://x.threatbook.com/v5/vulIntelligence | Severity High or Critical |
| Knownsec Seebug Vuln DB | https://www.seebug.org/ | Severity High or Critical |
| Venustech Vulnerability Advisories | https://www.venustech.com.cn/new_type/aqtg/ | Severity High or Critical |
| CISA KEV | https://www.cisa.gov/known-exploited-vulnerabilities-catalog | Push all |
| Struts2 Security Bulletins | [Struts2 Security Bulletins](https://cwiki.apache.org/confluence/display/WW/Security+Bulletins) | Severity High or Critical |

> All information comes from publicly available pages. If anything infringes rights, please open an issue and I will remove the related source.
>
> Feedback on better sources is also welcome — they need to be timely and allow filtering for valuable vulnerabilities.

Specifically, there are two push cases (both have built-in deduplication, so the same item is not pushed twice):

- A newly created vulnerability matches the push policy → push immediately
- A newly created vulnerability did not match the policy, but later updates make it match → also push

![app](./.github/assets/app.jpg)

## What to Do After a Push (Suggested Workflow)

WatchVuln handles **monitoring and alerting**. After you receive a push, follow this closed loop (trim to your org’s process as needed):

1. **Receive**: Get the alert in DingTalk / WeCom (or similar), or open the local [Vulnerability Intelligence Board](#vulnerability-intelligence-board) for details
2. **Triage**: Confirm whether the product/version is in your asset scope; use tags (e.g. public POC, KEV) and the AI analysis prompt to judge exploit conditions and impact
3. **Severity & routing**:
   - Not affected in this environment → record the conclusion, close or archive
   - Needs fix but not urgent → file a regular vuln ticket (patch window / change process)
   - **Critical, KEV, actively exploited in the wild, and your environment is affected** → escalate to your org’s incident response (this tool does not replace IR standards)
4. **Remediate**: Patch, upgrade, or apply temporary mitigations (WAF / isolation / hardening, etc.)
5. **Retest & close**: Close the ticket after verifying the vuln is no longer exploitable; update assets and baselines if needed

> Full incident command, notification, recovery, and related requirements should follow your organization’s emergency management standards; this section only defines the handoff points from WatchVuln alerts.

## Quick Start

Supported push channels:

- [DingTalk group bot](https://open.dingtalk.com/document/robots/custom-robot-access)
- [WeCom (WeChat Work) group bot](https://open.work.weixin.qq.com/help2/pc/14931)
- [Lark / Feishu group bot](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN)
- [Lanxin group bot](https://developer.lanxin.cn/official/article?id=646ecae03d4e4adb7039c0e4&module=development-help&article_id=646f193b3d4e4adb7039c21c)
- [ServerChan](https://sct.ftqq.com/)
- [PushPlus](https://pushplus.plus/)
- [Slack Webhook](https://docs.slack.dev/messaging/sending-messages-using-incoming-webhooks/)
- [Telegram Bot](https://core.telegram.org/bots/tutorial)
- [Custom Bark service](https://github.com/Finb/Bark)
- [Custom Webhook service](./examples/webhook)

### Using Docker

With Docker, environment variables are the recommended way to configure the service:

| Env var | Description | Default |
|-------------------------|-----------------------------------------------------------------------------------|---------------------------------------------------|
| `DB_CONN` | Database connection string; see [Database Connection](#database-connection) | `sqlite3://vuln_v3.sqlite3` |
| `DINGDING_ACCESS_TOKEN` | The `access_token` part of the DingTalk bot URL | |
| `DINGDING_SECRET` | DingTalk bot sign secret (sign mode only) | |
| `LARK_ACCESS_TOKEN` | The part after `/open-apis/bot/v2/hook/` in the Lark bot URL; a full URL for self-hosted Lark is also supported | |
| `LARK_SECRET` | Lark bot sign secret (sign mode only) | |
| `WECHATWORK_KEY ` | The `key` part of the WeCom bot URL | |
| `SERVERCHAN_KEY ` | ServerChan `SCKEY` | |
| `WEBHOOK_URL` | Full URL of a custom webhook service | |
| `BARK_URL` | Full Bark service URL; path must include DeviceKey | |
| `PUSHPLUS_KEY` | PushPlus token | |
| `LANXIN_DOMAIN` | Domain of the Lanxin webhook bot | |
| `LANXIN_TOKEN` | Hook token of the Lanxin webhook bot | |
| `LANXIN_SECRET` | Sign secret of the Lanxin webhook bot | |
| `TELEGRAM_BOT_TOKEN` | Telegram Bot Token | |
| `TELEGRAM_CHAT_IDS` | Telegram chat ID list to send to, comma-separated | |
| `SLACK_WEBHOOK_URL` | Full Slack webhook URL | |
| `SLACK_CHANNEL` | Slack channel to push to | |
| `SOURCES` | Enabled vulnerability sources, comma-separated. Options: `avd`, `chaitin`, `nox`/`ti`, `oscs`, `threatbook`, `seebug`, `struts2`, `kev`, `venustech` | `avd,chaitin,nox,oscs,threatbook,seebug,struts2,kev,venustech` |
| `INTERVAL` | Check interval; supports seconds `60s`, minutes `10m`, hours `1h`; minimum `1m` | `30m` |
| `ENABLE_CVE_FILTER` | Enable CVE filter; when on, the same CVE from multiple sources is pushed only once | `true` |
| `NO_FILTER` | Disable the push filter policies above; all newly discovered vulns will be pushed | `false` |
| `NO_START_MESSAGE` | Disable the “init finished” push after startup (`true`/`1` enables, `false` disables) | `false` |
| `WHITELIST_FILE` | Whitelist file for push filtering; see [Push Content Filtering](#push-content-filtering) | |
| `BLACKLIST_FILE` | Blacklist file for push filtering; see [Push Content Filtering](#push-content-filtering) | |
| `DIFF` | Skip the init phase; check for updates and push immediately | |
| `HTTPS_PROXY` | Proxy for all requests; see [Proxy Configuration](#proxy-configuration) | |
| `GO_SKIP_TLS_CHECK` | Skip TLS verification (`true`/`1` enables); same as `-k/--insecure`; see [Proxy Configuration](#proxy-configuration) | `false` |
| `NO_SLEEP` | Disable night sleep; run 24/7! See [Other](#other) | `false` |
| `WEB_ADDR` | Vulnerability board listen address, e.g. `127.0.0.1:8765` (empty = disabled) | |
| `WEB_PUBLIC_HOST` | Public host (IP or domain) for board/RSS/push links; port is taken from `WEB_ADDR` automatically | |
| `WEB_PUBLIC_URL` | Full public URL for the board (optional; overrides `WEB_PUBLIC_HOST`) | |
| `PUSHER_FILE` | Standalone yaml/json pusher list file (e.g. multiple DingTalk bots); see `pushers.example.yaml` | |

Example with a DingTalk bot:

```bash
docker run --restart always -d \
  -e DINGDING_ACCESS_TOKEN=xxxx \
  -e DINGDING_SECRET=xxxx \
  -e INTERVAL=30m \
  -e ENABLE_CVE_FILTER=true \
  ghcr.io/shellsec/watchvuln3:latest
```

<details><summary>Disable the “init finished” push on startup</summary>

On first start, a summary is pushed to the group by default (version, local vuln count, source list, etc.). To skip that message:

```bash
# Docker
docker run --restart always -d \
  -e NO_START_MESSAGE=true \
  -e DINGDING_ACCESS_TOKEN=xxxx \
  -e DINGDING_SECRET=xxxx \
  ghcr.io/shellsec/watchvuln3:latest

# Binary (Windows example; see “Recommended startup example” above)
.\watchvuln.exe --dt YOUR_DINGDING_ACCESS_TOKEN --ds "YOUR_DINGDING_SIGN_SECRET" -nm
```

Via config file: set `no_start_message: true` in the yaml.

</details>

<details><summary>Multiple DingTalk bots / multiple same-type pushers</summary>

**Config file (recommended)**: list multiple entries with the same `type` under `pusher`:

```yaml
pusher:
  - type: dingding
    access_token: "GROUP1_TOKEN"
    sign_secret: "GROUP1_SECRET"
  - type: dingding
    access_token: "GROUP2_TOKEN"
    sign_secret: "GROUP2_SECRET"
```

You can also run `watchvuln init-config` to generate a template. For CLI single-flag usage with a separate pusher file:

```bash
watchvuln --pusher-file pushers.example.yaml -c config.yaml
```

</details>

<details><summary>Vulnerability Intelligence Board (Web, no login)</summary>

Browse collected vulnerabilities from the local database — not limited to the “recently pushed” list. Supports search, filters, and sorting; one-click copy of an analysis prompt and jump to **ChatGPT / Gemini / DeepSeek**; no login.

Add `--web-addr` when starting monitoring. See the “Recommended startup example” above for a full example.

```bash
# Board only (no collection, no push)
watchvuln board --web-addr 127.0.0.1:8765

# Docker
-e WEB_ADDR=0.0.0.0:8765
```

Open `http://127.0.0.1:8765/` locally; if listening on `0.0.0.0`, use the LAN IP. The config key is `web_addr`. UI preview: [top of this doc](#vulnerability-intelligence-board).

**Board features**

| Capability | Description |
|------|------|
| Search | Title, CVE, description keywords |
| Filter | Severity, source |
| Sort | Default **by disclosure date** (most recently published CVEs first); can switch to **by DB update** (records most recently synced or changed by the program first) |
| AI analysis | **ChatGPT** / **Gemini** / **DS** next to titles: auto-copy analysis prompt (title, CVE, severity, disclosure date, original link), then open the corresponding site |
| RSS | `http://HOST:PORT/feed.xml`, latest **50 pushed** vulns; **RSS** button next to the title for feed readers |
| REST API | `http://CURRENT_HOST/api`, no login; **API / MCP** next to the title copies examples |
| MCP | `http://CURRENT_HOST/mcp` (Streamable HTTP); one-click copy of Cursor / Claude Code config |
| Detail | Click a table row for description, tags, remediation advice, and reference links |
| Pagination | 30 items per page |

**AI one-click analysis**

Each vulnerability title has three buttons on the right. Clicking one will:

1. **First copy** the prefilled analysis prompt to the clipboard (sync copy preferred; falls back to Clipboard API)
2. On success, the button briefly shows **“Copied”** (~2 seconds) and opens the corresponding AI site
3. In the AI chat box, press **Ctrl+V** to paste

| Button | Opens | Notes |
|------|----------|------|
| ChatGPT | [chatgpt.com](https://chatgpt.com/) | Short prompts may also go via URL query; longer ones rely on clipboard |
| Gemini | [gemini.google.com](https://gemini.google.com/app) | Mainly paste from clipboard |
| DS | [chat.deepseek.com](https://chat.deepseek.com/) | Mainly paste from clipboard |

If auto-copy fails, a dialog shows the full prompt so you can select/copy manually, then click “Open site”.

> Prefer `http://127.0.0.1:8765/` for the board; some browsers restrict the clipboard on LAN IPs (`http://192.168.x.x`) — check site permissions via the lock icon in the address bar.

**List column meanings**

| Column | Meaning |
|----|------|
| Disclosure date | Public disclosure time from the vulnerability source |
| DB update | Last update time of this record in the local DB (severity/tag changes, resync, etc. all update this) |

Whether page 1 shows the “newest” depends on the current sort: use **disclosure date** for recently published vulns; use **DB update** for recently changed records.

**RSS feed**

- URL: `http://HOST:8765/feed.xml` (same port as the board)
- Content: latest 50 **pushed** vulns, same scope as DingTalk pushes
- Startup logs print the reachable Feed URL (when `web_public_host` is set)

**Public board / RSS address (recommended when listening on `0.0.0.0`)**

```yaml
web_addr: "0.0.0.0:8766"
web_public_host: "192.168.1.100"   # port follows web_addr automatically
# or full URL: web_public_url: "http://192.168.1.100:8766"
```

After configuration, DingTalk pushes append a board link at the end; startup log example: `vuln board rss feed: http://192.168.1.100:8766/feed.xml`

**REST API and MCP (one-click connect)**

The board and the collector share the same HTTP port. **API / MCP URLs follow the address you used to open the board** (`window.location.origin` / request `Host`). Open `http://192.168.1.100:8765/` and that IP is the connect URL; after DHCP moves you to `192.168.1.200`, open the board on the new IP and the copy buttons update automatically. When listening on `0.0.0.0`, any reachable IP of this host hits the same `/api` and `/mcp`. Config already saved in Cursor does not rewrite itself — reopen the board and copy again after an IP change.

Click **API / MCP** next to the page title to copy the current origin, Cursor config, Claude Code command, and curl examples. You can also open `GET /api` for the catalog (`base_url` is the current Host).

| Endpoint | Description |
|------|------|
| `GET /api` | Catalog, MCP URL, examples |
| `GET /api/stats` | Totals and severity breakdown |
| `GET /api/sources` | Source list |
| `GET /api/vulns` | Paginated search: `q` `severity` `source` `sort` `page` `limit` |
| `GET /api/vuln` | One record: `id` or `cve` or `key` |
| `POST /mcp` | MCP Streamable HTTP (tools: `search_vulns` `get_vuln` `list_sources` `get_stats`) |

```bash
HOST=http://192.168.1.100:8765

curl -s "$HOST/api/vulns?q=CVE-2024&severity=严重&limit=5"
curl -s "$HOST/api/vuln?cve=CVE-2024-0001"
curl -s "$HOST/api/stats"

# Cursor mcp.json
# { "mcpServers": { "watchvuln": { "type": "http", "url": "http://192.168.1.100:8765/mcp" } } }

# Claude Code
claude mcp add --transport http watchvuln http://192.168.1.100:8765/mcp
```

</details>

### Common Commands

| Command | Description |
|------|------|
| `watchvuln list-sources` | List all source IDs, names, and URLs |
| `watchvuln init-config -o config.yaml` | Generate a commented config template |
| `watchvuln board` | Run the Vulnerability Intelligence Board only |

On startup, a config self-check is printed (sources, pusher count, board address, etc.).

You can also use this repo’s `docker-compose.yaml` and start with `docker compose`.

If the image is not published yet, build locally from the repo root:

```bash
docker build -t ghcr.io/shellsec/watchvuln3:latest .
```

On each update, pull or rebuild the image:

```bash
docker pull ghcr.io/shellsec/watchvuln3:latest
# or
docker build -t ghcr.io/shellsec/watchvuln3:latest .
```


<details><summary>Using a WeCom group bot</summary>

```bash
docker run --restart always -d \
  -e WECHATWORK_KEY=xxxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>Lark / Feishu group bot</summary>

```bash
docker run --restart always -d \
  -e LARK_ACCESS_TOKEN=xxxx \
  -e LARK_SECRET=xxxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>Using a Lanxin Webhook bot</summary>

```bash
docker run --restart always -d \
  -e LANXIN_DOMAIN=xxx \
  -e LANXIN_TOKEN=xxx \
  -e LANXIN_SECRET=xxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>Using ServerChan</summary>

```bash
docker run --restart always -d \
  -e SERVERCHAN_KEY=xxxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>Using PushPlus</summary>

```bash
docker run --restart always -d \
  -e PUSHPLUS_KEY=xxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>Using Slack Webhook</summary>

```bash
# Put your Slack Webhook here (do not commit a real URL to a public repo)
docker run --restart always -d \
  -e SLACK_WEBHOOK_URL=YOUR_SLACK_WEBHOOK_URL \
  -e SLACK_CHANNEL=#your-channel \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

> Important: Slack Incoming Webhooks are sensitive — never commit them to a public GitHub repo.

</details>

<details><summary>Using a Telegram bot</summary>

```bash
docker run --restart always -d \
  -e TELEGRAM_BOT_TOKEN=xxx \
  -e TELEGRAM_CHAT_IDS=1111,2222 \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>


<details><summary>Using a custom Bark service</summary>

```bash
docker run --restart always -d \
  -e BARK_URL=http://xxxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>Using a custom Webhook service</summary>

With a custom webhook server you can plug into other services. See: [example](./examples/webhook)

```bash
docker run --restart always -d \
  -e WEBHOOK_URL=http://xxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>



<details><summary>Using multiple services</summary>

If keys for multiple services are configured, each of them takes effect. Example with DingTalk and WeCom:

```bash
docker run --restart always -d \
  -e DINGDING_ACCESS_TOKEN=xxxx \
  -e DINGDING_SECRET=xxxx \
  -e WECHATWORK_KEY=xxxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>


On first run, a full local database is built (about 1 minute). Follow progress with `docker logs -f [containerId]`. When done, you will get a group message indicating the service is running normally.

### Using the Binary

Download the binary for your platform from [GitHub Releases](https://github.com/shellsec/watchvuln3/releases), or build from the repo root:

```bash
# Single platform (current OS)
go build -trimpath -ldflags "-s -w" -o watchvuln .

# One-shot cross-compile for multiple platforms on Windows (output to dist/, do not commit)
.\scripts\build-release.ps1

# Linux / macOS
./scripts/build-release.sh v3.1.0
```

Example `dist/` artifacts: `watchvuln-windows-amd64.exe`, `watchvuln-linux-amd64`, `watchvuln-linux-arm64`, etc. The project uses pure Go SQLite — **no CGO required** — so you can cross-compile Linux binaries on Windows.

**Recommended startup example (DingTalk + WeCom + disable init push + vuln board)**:

```powershell
.\watchvuln.exe --dt YOUR_DINGDING_ACCESS_TOKEN --ds "YOUR_DINGDING_SIGN_SECRET" --wk YOUR_WECHATWORK_KEY -nm --interval 30m --web-addr 0.0.0.0:8765
```

| Flag | Meaning |
|------|------|
| `--dt` / `--ds` | DingTalk bot Token and sign secret |
| `--wk` | WeCom group bot Key |
| `-nm` | Do not push the “init finished” message on startup (same as `--no-start-message`) |
| `--interval 30m` | Check every 30 minutes |
| `--web-addr 0.0.0.0:8765` | Enable the Vulnerability Intelligence Board; open `http://HOST:8765/` in a browser |
| `--web-public-host 192.168.1.100` | Public IP for board/RSS/push links (port follows `--web-addr`) |

> Replace `YOUR_*` with your own secrets; do not commit them to a public repo.

CLI flags map one-to-one with the Docker environment variables described above.

```bash
USAGE:
   watchvuln [global options] command [command options] [arguments...]

GLOBAL OPTIONS:
   --config value, -c value  config file path, support json or yaml

   [Push Options]

   --bark-url value, --bark value             your bark server url, ex: http://127.0.0.1:1111/DeviceKey
   --blacklist-file value, --bf value         specify a file that contains some keywords, vulns with these products will NOT be pushed
   --dingding-access-token value, --dt value  webhook access token of dingding bot
   --dingding-sign-secret value, --ds value   sign secret of dingding bot
   --lanxin-domain value, --lxd value         your lanxin server url, ex: https://apigw-example.domain
   --lanxin-hook-token value, --lxt value     lanxin hook token
   --lanxin-sign-secret value, --lxs value    sign secret of lanxin
   --lark-access-token value, --lt value      webhook access token/url of lark
   --lark-sign-secret value, --ls value       sign secret of lark
   --pushplus-key value, --pk value           send key for push plus
   --serverchan-key value, --sk value         send key for server chan
   --slack-channel value, --sc value          specify slack channel, eg, #security_vulns
   --slack-webhook-url value, --sw value      specify slack webhook url (use YOUR_SLACK_WEBHOOK_URL, do not commit secrets)
   --telegram-bot-token value, --tgtk value   telegram bot token, ex: 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
   --telegram-chat-ids value, --tgids value   chat ids want to send on telegram, ex: 123456,4312341,123123
   --webhook-url value, --webhook value       your webhook server url, ex: http://127.0.0.1:1111/webhook
   --wechatwork-key value, --wk value         webhook key of wechat work
   --whitelist-file value, --wf value         specify a file that contains some keywords, vulns with these keywords will be pushed

   [Launch Options]

   --db-conn value, --db value  database connection string (default: "sqlite3://vuln_v3.sqlite3")
   --diff                       skip init vuln db, push new vulns then exit (default: false)
   --enable-cve-filter          enable a filter that vulns from multiple sources with same cve id will be sent only once (default: true)
   --interval value, -i value   checking every [interval], supported format like 30s, 30m, 1h (default: "30m")
   --no-filter, --nf            ignore the valuable filter and push all discovered vulns (default: false)
   --no-github-search, --ng     don't search github repos and pull requests for every cve vuln (default: false)
   --no-sleep, --ns             don't sleep in night, run every interval (default: false)
   --no-start-message, --nm, --quiet
                                disable init finished push message when server starts (default: false)
   --proxy value, -x value      set request proxy, support socks5://xxx or http(s)://
   --sources value, -s value    set vuln sources (default: "avd,chaitin,nox,oscs,threatbook,seebug,struts2,kev,venustech")
   --web-addr value             vuln intelligence board listen address, e.g. 127.0.0.1:8765
   --web-public-url value       full public URL for board/RSS links in push messages
   --web-public-host value      public host/IP for board/RSS links; port from --web-addr

   [Other Options]

   --debug, -d     set log level to debug, print more details (default: false)
   --help, -h      show help (default: false)
   --insecure, -k  allow insecure server connections when using SSL/TLS (default: false)
   --test, -T      use to test message pusher, three mocked messages will be pushed (default: false)
   --version, -v   print the version (default: false)
```

Pass the relevant tokens as flags, e.g. with a DingTalk group bot:

```
$ ./watchvuln --dt DINGDING_ACCESS_TOKEN --ds DINGDING_SECRET -i 30m
```

<details><summary>Using a WeCom group bot</summary>

```
$ ./watchvuln --wk WECHATWORK_KEY -i 30m
```

</details>

<details><summary>Using a Lark / Feishu group bot</summary>

```bash
$ ./watchvuln --lt LARK_ACCESS_TOKEN --ls LARK_SECRET -i 30m
```

</details>

<details><summary>Using a Lanxin group bot</summary>

```
$ ./watchvuln --lxd xxxx --lxt xxx --lxs xxx -i 30m
```

</details>


<details><summary>Using ServerChan</summary>

```
$ ./watchvuln --sk xxxx -i 30m
```

</details>

<details><summary>Using PushPlus</summary>

```
$ ./watchvuln --pk xxxx -i 30m
```

</details>

<details><summary>Using Slack Webhook </summary>

```
# Put your Slack Webhook here (do not commit a real URL to a public repo)
$ ./watchvuln --sw YOUR_SLACK_WEBHOOK_URL --sc '#your_channel' -i 30m
```

> Important: Slack Incoming Webhooks are sensitive — never commit them to a public GitHub repo.

</details>


<details><summary>Using a Telegram bot</summary>

```
$ ./watchvuln --tgtk xxxx --tgids 1111,2222 -i 30m
```

</details>


<details><summary>Using a custom Bark service</summary>

```
$ ./watchvuln --bark http://xxxx -i 30m
```

</details>

<details><summary>Using a custom Webhook service</summary>

With a custom webhook server you can plug into other services. See: [example](./examples/webhook)

```
$ ./watchvuln --webhook http://xxxx -i 30m
```

</details>


<details><summary>Using multiple services</summary>

If keys for multiple services are configured, each of them takes effect. Example with DingTalk and WeCom:

```
$ ./watchvuln --dt DINGDING_ACCESS_TOKEN --ds DINGDING_SECRET --wk WECHATWORK_KEY -i 30m
```

</details>

## Config File

See details: [config file usage](CONFIG.md) (CONFIG.md is currently in Chinese).

## Database Connection

sqlite3 is used by default; the database file is `vuln_v3.sqlite3`. To use another database, set the connection string via `--db` or the `DB_CONN` env var. Supported databases:

- `sqlite3://filename`
- `mysql://user:pass@host:port/dbname`
- `postgres://user:pass@host:port/dbname`

Note: this project does not guarantee backward data compatibility. Upgrades may be incompatible; if errors occur, you may need to delete the DB and start over.

## Proxy Configuration

watchvuln supports an upstream proxy to work around network limits, in two ways:

- Env var `HTTPS_PROXY`
- CLI flag `--proxy`/`-x`

Both `socks5://xxxx` and `http(s)://xxkx` forms are supported.

The flag `-k/--insecure` or env var `GO_SKIP_TLS_CHECK=1` disables TLS verification (sets `InSecureSkipVerify` to `true`), which is useful when capturing traffic for debugging.

## Push Content Filtering

If you only want to push vulnerabilities for certain products, use a whitelist or blacklist. Both take a file with one product name per line, e.g.:

```txt
Apache
Weaver
```

Tip: if you run with `Docker`, mount the file into the container, e.g.:

```bash
echo "Apache" > whitelist.txt

docker run -v $(pwd):/config \
  -e WHITELIST_FILE=/config/whitelist.txt \
  -e xxxx=xxxxx
  ghcr.io/shellsec/watchvuln3:latest
```

### Whitelist filtering

Specify the whitelist file with CLI `-wf` or env var `WHITELIST_FILE`. When a new vuln is found, its **title** and **description** are checked against any line in the whitelist; if none match, it will not be pushed.

### Blacklist filtering

Specify the blacklist file with CLI `-bf` or env var `BLACKLIST_FILE`. When a new vuln is found, only its **title** is checked against any line in the blacklist; matches are not pushed. To avoid unexpected misses, the blacklist does **not** match against the **description**.

## Other

To reduce unnecessary overtime grind, the tool sleeps between 00:00 and 07:00 and will not run — make sure your server’s clock is correct!
