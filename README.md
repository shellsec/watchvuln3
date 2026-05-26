# WatchVuln 高价值漏洞采集与推送

[![GitHub Release](https://img.shields.io/github/v/release/shellsec/watchvuln3?label=release)](https://github.com/shellsec/watchvuln3/releases)
[![License](https://img.shields.io/github/license/shellsec/watchvuln3)](https://github.com/shellsec/watchvuln3)

**仓库**: [github.com/shellsec/watchvuln3](https://github.com/shellsec/watchvuln3) · **当前版本**: v3.0.0

## ☕ 请我喝可乐

开源不易，欢迎赞助支持：

👉 [爱发电](https://www.ifdian.net/a/shellsec)

## 关于本仓库

本项目基于原作者 **[zema1/watchvuln](https://github.com/zema1/watchvuln)** 维护与扩展，核心思路（多源高价值漏洞采集、过滤、多渠道推送）均来自上游，在此向原作者 **zema1** 表示感谢。

`shellsec/watchvuln3` 在 v3.0 中主要做了这些事：

**问题修复**

- 修复奇安信 TI 数据源 `Referer` / `Origin` 请求头错误
- 修复 `GO_SKIP_TLS_CHECK` 在 Docker / CLI 下不生效
- 修复 CVE 去重后仍每个周期重复处理同一漏洞
- 修复初始化多数据源时的并发数据竞争
- 统一默认 `sources` 与文档；布尔环境变量支持标准 `true`/`false` 解析

**功能增强**

- 漏洞情报看板（`--web-addr` / `watchvuln board`，本地浏览库内情报，支持按披露日期或入库更新排序，无登录）
- 配置文件 / `--pusher-file` 支持**多个同类型推送**（如多个钉钉群）
- `-nm` / `--quiet` 关闭启动时的「初始化完成」推送
- CLI：`list-sources`、`init-config`、启动配置自检

详细变更见 [CHANGELOG.md](./CHANGELOG.md)。若上游合并了同类修复，以 [zema1/watchvuln](https://github.com/zema1/watchvuln) 为准；本仓库也会持续跟进必要同步。

---

众所周知，CVE 漏洞库中 99% 以上的漏洞只是无现实意义的编号。我想集中精力看下当下需要关注的高价值漏洞有哪些，而不是被各类 RSS
和公众号的 ~~威胁情报~~ 淹没。 于是写了这个小项目来抓取部分高质量的漏洞信息源然后做推送。 `WatchVuln`意为**监测**
漏洞更新，同时也表示这些漏洞需要**注意**。

当前抓取了这几个站点的数据:

| 名称                         | 地址                                                                                              | 推送策略                                             |
|----------------------------|-------------------------------------------------------------------------------------------------|--------------------------------------------------|
| 阿里云漏洞库                     | https://avd.aliyun.com/high-risk/list                                                           | 等级为高危或严重                                         |
| 长亭漏洞库                      | https://stack.chaitin.com/vuldb/index                                                           | 等级为高危或严重**并且**标题含中文                              |
| OSCS开源安全情报预警               | https://www.oscs1024.com/cm                                                                     | 等级为高危或严重**并且**包含 `预警` 标签                         |
| 奇安信威胁情报中心                  | https://ti.qianxin.com/                                                                         | 等级为高危严重**并且**包含 `奇安信CERT验证` `POC公开` `技术细节公布`标签之一 |
| 微步在线研究响应中心(公众号)            | https://x.threatbook.com/v5/vulIntelligence                                                     | 等级为高危或严重                                         |
| 知道创宇Seebug漏洞库              | https://www.seebug.org/                                                                         | 等级为高危或严重                                         |
| 启明星辰漏洞通告                   | https://www.venustech.com.cn/new_type/aqtg/                                                     | 等级为高危或严重                                         |
| CISA KEV                   | https://www.cisa.gov/known-exploited-vulnerabilities-catalog                                    | 全部推送                                             |
| Struts2 Security Bulletins | [Struts2 Security Bulletins](https://cwiki.apache.org/confluence/display/WW/Security+Bulletins) | 等级为高危或严重                                         |

> 所有信息来自网站公开页面, 如果有侵权，请提交 issue, 我会删除相关源。
>
> 如果有更好的信息源也可以反馈给我，需要能够响应及时 & 有办法过滤出有价值的漏洞

具体来说，消息的推送有两种情况, 两种情况有内置去重，不会重复推送:

- 新建的漏洞符合推送策略，直接推送,
- 新建的漏洞不符合推送策略，但漏洞信息被更新后符合了推送策略，也会被推送

![app](./.github/assets/app.jpg)

## 快速使用

支持下列推送方式:

- [钉钉群组机器人](https://open.dingtalk.com/document/robots/custom-robot-access)
- [微信企业版群组机器人](https://open.work.weixin.qq.com/help2/pc/14931)
- [飞书群组机器人](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN)
- [蓝信群组机器人](https://developer.lanxin.cn/official/article?id=646ecae03d4e4adb7039c0e4&module=development-help&article_id=646f193b3d4e4adb7039c21c)
- [Server 酱](https://sct.ftqq.com/)
- [PushPlus](https://pushplus.plus/)
- [Slack Webhook](https://docs.slack.dev/messaging/sending-messages-using-incoming-webhooks/)
- [Telegram Bot](https://core.telegram.org/bots/tutorial)
- [自定义 Bark 服务](https://github.com/Finb/Bark)
- [自定义 Webhook 服务](./examples/webhook)

### 使用 Docker

Docker 方式推荐使用环境变量来配置服务参数

| 环境变量名                   | 说明                                                                                | 默认值                                               |
|-------------------------|-----------------------------------------------------------------------------------|---------------------------------------------------|
| `DB_CONN`               | 数据库链接字符串，详情见 [数据库连接](#数据库连接)                                                      | `sqlite3://vuln_v3.sqlite3`                       |
| `DINGDING_ACCESS_TOKEN` | 钉钉机器人 url 的 `access_token` 部分                                                     |                                                   |
| `DINGDING_SECRET`       | 钉钉机器人的加签值 （仅支持加签方式）                                                               |                                                   |
| `LARK_ACCESS_TOKEN`     | 飞书机器人 url 的 `/open-apis/bot/v2/hook/` 后的部分, 也支持直接指定完整的 url 来访问私有部署的飞书             |                                                   |
| `LARK_SECRET`           | 飞书机器人的加签值 （仅支持加签方式）                                                               |                                                   |
| `WECHATWORK_KEY `       | 微信机器人 url 的 `key` 部分                                                              |                                                   |
| `SERVERCHAN_KEY `       | Server酱的 `SCKEY`                                                                  |                                                   |
| `WEBHOOK_URL`           | 自定义 webhook 服务的完整 url                                                             |                                                   |
| `BARK_URL`              | Bark 服务的完整 url, 路径需要包含 DeviceKey                                                  |                                                   |
| `PUSHPLUS_KEY`          | PushPlus的token                                                                    |                                                   |
| `LANXIN_DOMAIN`         | 蓝信webhook机器人的域名                                                                   |                                                   |
| `LANXIN_TOKEN`          | 蓝信webhook机器人的hook token                                                           |                                                   |
| `LANXIN_SECRET`         | 蓝信webhook机器人的签名                                                                   |                                                   |
| `TELEGRAM_BOT_TOKEN`    | Telegram Bot Token                                                                |                                                   |
| `TELEGRAM_CHAT_IDS`     | Telegram Bot 需要发送给的 chat 列表，使用 `,` 分割                                             |                                                   |
| `SLACK_WEBHOOK_URL`     | slack webhook 完整 url                                                              |                                                   |
| `SLACK_CHANNEL`         | 要推送到的 slack 频道                                                                    |                                                   |
| `SOURCES`               | 启用哪些漏洞信息源，逗号分隔。可选：`avd`, `chaitin`, `nox`/`ti`, `oscs`, `threatbook`, `seebug`, `struts2`, `kev`, `venustech` | `avd,chaitin,nox,oscs,threatbook,seebug,struts2,kev,venustech` |
| `INTERVAL`              | 检查周期，支持秒 `60s`, 分钟 `10m`, 小时 `1h`, 最低 `1m`                                        | `30m`                                             |
| `ENABLE_CVE_FILTER`     | 启用 CVE 过滤，开启后多个数据源的统一 CVE 将只推送一次                                                  | `true`                                            |
| `NO_FILTER`             | 禁用上述推送过滤策略，所有新发现的漏洞都会被推送                                                          | `false`                                           |
| `NO_START_MESSAGE`      | 禁用服务启动后的「初始化完成」推送（设为 `true`/`1` 生效，`false` 关闭）                                      | `false`                                           |
| `WHITELIST_FILE`        | 指定推送漏洞的白名单列表文件, 详情见 [推送内容筛选](#推送内容筛选)                                             |                                                   |
| `BLACKLIST_FILE`        | 指定推送漏洞的黑名单列表文件, 详情见 [推送内容筛选](#推送内容筛选)                                             |                                                   |
| `DIFF`                  | 跳过初始化阶段，转而直接检查漏洞更新并推送                                                             |                                                   |
| `HTTPS_PROXY`           | 给所有请求配置代理, 详情见 [配置代理](#配置代理)                                                      |                                                   |
| `GO_SKIP_TLS_CHECK`     | 跳过 tls 校验（`true`/`1` 生效），等同 `-k/--insecure`，详情见 [配置代理](#配置代理)                      | `false`                                           |
| `NO_SLEEP`              | 禁用夜晚休眠，全天24小时无休！[其他](#其他)                                                         | `false`                                           |
| `WEB_ADDR`              | 漏洞情报看板监听地址，如 `127.0.0.1:8765`（空则关闭）                                                |                                                   |
| `PUSHER_FILE`           | 独立 yaml/json 推送列表文件（多钉钉等），见 `pushers.example.yaml`                                  |                                                   |

比如使用钉钉机器人

```bash
docker run --restart always -d \
  -e DINGDING_ACCESS_TOKEN=xxxx \
  -e DINGDING_SECRET=xxxx \
  -e INTERVAL=30m \
  -e ENABLE_CVE_FILTER=true \
  ghcr.io/shellsec/watchvuln3:latest
```

<details><summary>关闭启动后的「初始化完成」推送</summary>

首次启动默认会向群推送一条初始化摘要（版本、本地漏洞数、数据源列表等）。若不需要该消息：

```bash
# Docker
docker run --restart always -d \
  -e NO_START_MESSAGE=true \
  -e DINGDING_ACCESS_TOKEN=xxxx \
  -e DINGDING_SECRET=xxxx \
  ghcr.io/shellsec/watchvuln3:latest

# 二进制（Windows 示例见上文「推荐启动示例」）
.\watchvuln.exe --dt YOUR_DINGDING_ACCESS_TOKEN --ds "YOUR_DINGDING_SIGN_SECRET" -nm
```

配置文件方式：在 yaml 中设置 `no_start_message: true`。

</details>

<details><summary>多个钉钉 / 多个同类型推送</summary>

**配置文件（推荐）**：在 `pusher` 下写多条相同 `type`：

```yaml
pusher:
  - type: dingding
    access_token: "群1_TOKEN"
    sign_secret: "群1_SECRET"
  - type: dingding
    access_token: "群2_TOKEN"
    sign_secret: "群2_SECRET"
```

也可运行 `watchvuln init-config` 生成模板。CLI 单条参数方式可用独立推送文件：

```bash
watchvuln --pusher-file pushers.example.yaml -c config.yaml
```

</details>

<details><summary>漏洞情报看板（Web，无登录）</summary>

从本地数据库浏览已采集漏洞，不依赖「最近推送」列表。

与监测一起启动时加上 `--web-addr` 即可，完整示例见上文「推荐启动示例」。

```bash
# 仅启动看板（不采集、不推送）
watchvuln board --web-addr 127.0.0.1:8765

# Docker
-e WEB_ADDR=0.0.0.0:8765
```

本机访问 `http://127.0.0.1:8765/`；监听 `0.0.0.0` 时可用局域网 IP 访问。配置文件中对应项为 `web_addr`。

**看板功能**

| 能力 | 说明 |
|------|------|
| 搜索 | 标题、CVE、描述关键词 |
| 筛选 | 等级、数据来源 |
| 排序 | 默认**按披露日期**（最近公开的 CVE 在前）；可切换为**按入库更新**（最近被程序同步或变更的记录在前） |
| 详情 | 点击表格行查看描述、标签、修复建议、参考链接 |
| 分页 | 每页 30 条 |

**列表字段说明**

| 列 | 含义 |
|----|------|
| 披露日期 | 漏洞源上的公开披露时间 |
| 入库更新 | 本程序本地库中该条记录的最后更新时间（等级/标签变更、重新同步等都会更新） |

第一页是否为「最新」，取决于当前排序方式：查最近披露的漏洞用**披露日期**；查最近有变动的记录用**入库更新**。

</details>

### 常用命令

| 命令 | 说明 |
|------|------|
| `watchvuln list-sources` | 列出所有漏洞源 ID、名称、链接 |
| `watchvuln init-config -o config.yaml` | 生成带注释的配置模板 |
| `watchvuln board` | 仅运行漏洞情报看板 |

启动时会打印配置自检（数据源、推送通道数量、看板地址等）。

当然，你可以参考使用本仓库的 `docker-compose.yaml` 文件，使用 `docker compose` 来启动容器。

镜像尚未发布时，可在仓库根目录本地构建:

```bash
docker build -t ghcr.io/shellsec/watchvuln3:latest .
```

每次更新可重新拉取或构建镜像:

```bash
docker pull ghcr.io/shellsec/watchvuln3:latest
# 或
docker build -t ghcr.io/shellsec/watchvuln3:latest .
```


<details><summary>使用企业微信群组机器人</summary>

```bash
docker run --restart always -d \
  -e WECHATWORK_KEY=xxxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>飞书群组机器人</summary>

```bash
docker run --restart always -d \
  -e LARK_ACCESS_TOKEN=xxxx \
  -e LARK_SECRET=xxxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>使用蓝信 Webhook 机器人</summary>

```bash
docker run --restart always -d \
  -e LANXIN_DOMAIN=xxx \
  -e LANXIN_TOKEN=xxx \
  -e LANXIN_SECRET=xxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>使用 Server 酱</summary>

```bash
docker run --restart always -d \
  -e SERVERCHAN_KEY=xxxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>使用 PushPlus</summary>

```bash
docker run --restart always -d \
  -e PUSHPLUS_KEY=xxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>使用 Slack Webhook</summary>

```bash
# 此处填写你的 Slack Webhook（勿将真实 URL 写入公开仓库）
docker run --restart always -d \
  -e SLACK_WEBHOOK_URL=YOUR_SLACK_WEBHOOK_URL \
  -e SLACK_CHANNEL=#your-channel \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

> 重要：Slack Incoming Webhook 属于敏感信息，切勿提交到 GitHub 公开仓库。

</details>

<details><summary>使用Telegram 机器人</summary>

```bash
docker run --restart always -d \
  -e TELEGRAM_BOT_TOKEN=xxx \
  -e TELEGRAM_CHAT_IDS=1111,2222 \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>


<details><summary>使用自定义 Bark 服务</summary>

```bash
docker run --restart always -d \
  -e BARK_URL=http://xxxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>

<details><summary>使用自定义 Webhook 服务</summary>

通过自定义一个 webhook server，可以方便的接入其他服务, 实现方式可以参考: [example](./examples/webhook)

```bash
docker run --restart always -d \
  -e WEBHOOK_URL=http://xxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>



<details><summary>使用多种服务</summary>

如果配置了多种服务的密钥，那么每个服务都会生效， 比如使用钉钉和企业微信:

```bash
docker run --restart always -d \
  -e DINGDING_ACCESS_TOKEN=xxxx \
  -e DINGDING_SECRET=xxxx \
  -e WECHATWORK_KEY=xxxx \
  -e INTERVAL=30m \
  ghcr.io/shellsec/watchvuln3:latest
```

</details>


初次运行会在本地建立全量数据库，大约需要 1 分钟，可以使用 `docker logs -f [containerId]` 来查看进度,
完成后会在群内收到一个提示消息，表示服务已经在正常运行了。

### 使用二进制

前往 [GitHub Releases](https://github.com/shellsec/watchvuln3/releases) 下载对应平台的二进制，或在仓库根目录自行编译:

```bash
go build -o watchvuln.exe .
```

**推荐启动示例（钉钉 + 企业微信 + 关闭初始化推送 + 漏洞看板）**：

```powershell
.\watchvuln.exe --dt YOUR_DINGDING_ACCESS_TOKEN --ds "YOUR_DINGDING_SIGN_SECRET" --wk YOUR_WECHATWORK_KEY -nm --interval 30m --web-addr 0.0.0.0:8765
```

| 参数 | 含义 |
|------|------|
| `--dt` / `--ds` | 钉钉机器人 Token 与加签密钥 |
| `--wk` | 企业微信群机器人 Key |
| `-nm` | 不推送启动时的「初始化完成」消息（等同 `--no-start-message`） |
| `--interval 30m` | 每 30 分钟检查一次 |
| `--web-addr 0.0.0.0:8765` | 开启漏洞情报看板，浏览器访问 `http://本机IP:8765/` |

> 请将 `YOUR_*` 替换为你自己的密钥，勿提交到公开仓库。

命令行参数请参考 Docker 环境变量部分的说明，可以一一对应。

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

   [Other Options]

   --debug, -d     set log level to debug, print more details (default: false)
   --help, -h      show help (default: false)
   --insecure, -k  allow insecure server connections when using SSL/TLS (default: false)
   --test, -T      use to test message pusher, three mocked messages will be pushed (default: false)
   --version, -v   print the version (default: false)
```

在参数中指定相关 Token 即可, 比如使用钉钉群组机器人

```
$ ./watchvuln --dt DINGDING_ACCESS_TOKEN --ds DINGDING_SECRET -i 30m
```

<details><summary>使用企业微信群组机器人</summary>

```
$ ./watchvuln --wk WECHATWORK_KEY -i 30m
```

</details>

<details><summary>使用飞书群组机器人</summary>

```bash
$ ./watchvuln --lt LARK_ACCESS_TOKEN --ls LARK_SECRET -i 30m
```

</details>

<details><summary>使用蓝信群组机器人</summary>

```
$ ./watchvuln --lxd xxxx --lxt xxx --lxs xxx -i 30m
```

</details>


<details><summary>使用 Server 酱</summary>

```
$ ./watchvuln --sk xxxx -i 30m
```

</details>

<details><summary>使用 PushPlus</summary>

```
$ ./watchvuln --pk xxxx -i 30m
```

</details>

<details><summary>使用 Slack Webhook </summary>

```
# 此处填写你的 Slack Webhook（勿将真实 URL 写入公开仓库）
$ ./watchvuln --sw YOUR_SLACK_WEBHOOK_URL --sc '#your_channel' -i 30m
```

> 重要：Slack Incoming Webhook 属于敏感信息，切勿提交到 GitHub 公开仓库。

</details>


<details><summary>使用 Telegram 机器人</summary>

```
$ ./watchvuln --tgtk xxxx --tgids 1111,2222 -i 30m
```

</details>


<details><summary>使用自定义 Bark 服务</summary>

```
$ ./watchvuln --bark http://xxxx -i 30m
```

</details>

<details><summary>使用自定义 Webhook 服务</summary>

通过自定义一个 webhook server，可以方便的接入其他服务, 实现方式可以参考: [example](./examples/webhook)

```
$ ./watchvuln --webhook http://xxxx -i 30m
```

</details>


<details><summary>使用多种服务</summary>

如果配置了多种服务的密钥，那么每个服务都会生效， 比如使用钉钉和企业微信:

```
$ ./watchvuln --dt DINGDING_ACCESS_TOKEN --ds DINGDING_SECRET --wk WECHATWORK_KEY -i 30m
```

</details>

## 配置文件

进入查看详情 [使用配置文件](CONFIG.md)

## 数据库连接

默认使用 sqlite3 作为数据库，数据库文件为 `vuln_v3.sqlite3`，如果需要使用其他数据库，可以通过 `--db`
参数或是环境变量 `DB_CONN` 指定连接字符串，当前支持的数据库有:

- `sqlite3://filename`
- `mysql://user:pass@host:port/dbname`
- `postgres://user:pass@host:port/dbname`

注意：该项目不做数据向后兼容保证，版本升级可能存在数据不兼容的情况，如果报错需要删库重来。

## 配置代理

watchvuln 支持配置上游代理来绕过网络限制，支持两种方式:

- 环境变量 `HTTPS_PROXY`
- 命令行参数 `--proxy`/`-x`

支持 `socks5://xxxx` 或者 `http(s)://xxkx` 两种代理形式。

参数 `-k/--insecure` 或者环境变量 `GO_SKIP_TLS_CHECK=1` 可以禁用 tls 校验，即会设置 `InSecureSkipVerify` 为 `true`
，在抓包调试时会
比较有用。

## 推送内容筛选

如果你只想推送某些产品的漏洞，可以通过配置白名单或者黑名单来实现。这两个参数传入的都是一个文件，文件格式为每行一个产品名，比如:

```txt
Apache
泛微
```

温馨提示：如果你使用 `Docker` 来运行，可以通过挂载目录的方式将文件映射到容器内，比如:

```bash
echo "Apache" > whitelist.txt

docker run -v $(pwd):/config \
  -e WHITELIST_FILE=/config/whitelist.txt \
  -e xxxx=xxxxx
  ghcr.io/shellsec/watchvuln3:latest
```

### 白名单过滤

通过命令行参数 `-wf` 或者环境变量 `WHITELIST_FILE` 来指定白名单文件。在发现新漏洞时，将检查漏洞的 **标题** 和 **描述**
是否包含白名单的任意一行，全都不在的将不推送漏洞。

### 黑名单过滤

通过命令行参数 `-bf` 或者环境变量 `BLACKLIST_FILE` 来指定黑名单文件。在发现新漏洞时，将检查漏洞的 **标题** 是否包含黑名单的任意一行，
包含的将不推送漏洞。为了避免非预期的漏掉推送，黑名单**不会**检查漏洞的 **描述** 是否匹配。

## 其他

为了减少内卷，该工具在 00:00 到 07:00 间会去 sleep 不会运行, 请确保你的服务器是正确的时间！