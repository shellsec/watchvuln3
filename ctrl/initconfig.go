package ctrl

const InitConfigTemplate = `# WatchVuln 配置文件
# 使用: watchvuln -c config.yaml
# 仓库: https://github.com/shellsec/watchvuln3
#
# CLI 等价示例:
#   .\watchvuln.exe --dt YOUR_DINGDING_ACCESS_TOKEN --ds "YOUR_DINGDING_SIGN_SECRET" --wk YOUR_WECHATWORK_KEY -nm --interval 30m --web-addr 0.0.0.0:8765

db_conn: sqlite3://vuln_v3.sqlite3
sources: [ "avd", "chaitin", "nox", "oscs", "threatbook", "seebug", "struts2", "kev", "venustech" ]
interval: 30m
no_start_message: true
enable_cve_filter: true
no_github_search: false
no_sleep: false
diff_mode: false
skip_tls_verify: false
proxy: ""
web_addr: "0.0.0.0:8765"
# REST API: /api    MCP: /mcp（地址跟随访问 Host）
# 推送末尾看板链接：端口自动跟随 web_addr
web_public_host: "192.168.1.100"
white_keywords: [ ]
black_keywords: [ ]

pusher:
  - type: dingding
    access_token: "YOUR_DINGDING_ACCESS_TOKEN"
    sign_secret: "YOUR_DINGDING_SIGN_SECRET"
  - type: wechatwork
    key: "YOUR_WECHATWORK_KEY"
`
