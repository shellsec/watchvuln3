# 配置文件

仓库: [github.com/shellsec/watchvuln3](https://github.com/shellsec/watchvuln3)

Watchvuln 从 `v2.0.0` 版本开始支持从文件加载配置, 使用时需要使用 `-c` 参数指定配置文件路径, 如:

```
./watchvuln -c /path/to/config.yaml
./watchvuln -c /path/to/config.json
```

同时为了简化开发和减低理解成本，我们约定，**如果指定了配置文件，那么命令行指定的任何参数将不再生效**

## 文件格式

支持 `yaml` 和 `json` 两种格式的配置文件，这两个格式本质上是互通的，你可以选择自己喜欢的后缀。
一般，你只需将 `config.example.yaml` 的内容改一下即可，一个最简单的配置大概如下:

完整示例见仓库根目录 [`config.example.yaml`](./config.example.yaml)（含钉钉+企微+关闭启动推送+看板）。

最简示例:

```yaml
interval: 30m
no_start_message: true
web_addr: "0.0.0.0:8765"
pusher:
  - type: dingding
    access_token: "YOUR_DINGDING_ACCESS_TOKEN"
    sign_secret: "YOUR_DINGDING_SIGN_SECRET"
  - type: wechatwork
    key: "YOUR_WECHATWORK_KEY"
```

等价命令行:

```powershell
.\watchvuln.exe --dt YOUR_DINGDING_ACCESS_TOKEN --ds "YOUR_DINGDING_SIGN_SECRET" --wk YOUR_WECHATWORK_KEY -nm --interval 30m --web-addr 0.0.0.0:8765
```

聪明的你一定发现了，配置文件里的字段和命令行参数是一一对应的，这里就不再赘述了。

实际上，配置文件的出现主要是为了解决**多推送、多钉钉**的问题，比如你有两个钉钉群需要推送，那么可以写成这样

```yaml
db_conn: sqlite3://vuln_v3.sqlite3
sources: [ "avd", "chaitin", "nox", "oscs", "threatbook", "seebug", "struts2", "kev", "venustech" ]
interval: 30m
pusher:
  - type: dingding
    access_token: "xxxx"
    sign_secret: "yyyy"
   
  - type: dingding
    access_token: "pppp"
    sign_secret: "qqqq"
```

## 漏洞情报看板

在配置中增加 `web_addr`（或环境变量 `WEB_ADDR`、参数 `--web-addr`）即可在运行监测的同时提供本地 Web 看板，**无需登录**：

```yaml
web_addr: "127.0.0.1:8765"
```

仅查看看板时可执行：`./watchvuln board --web-addr 127.0.0.1:8765`

## 辅助命令

- `./watchvuln list-sources` — 查看可用漏洞源
- `./watchvuln init-config` — 生成 `config.yaml` 模板
