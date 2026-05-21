package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/zema1/watchvuln/ctrl"
	"github.com/zema1/watchvuln/grab"
	"github.com/zema1/watchvuln/web"
	"gopkg.in/yaml.v3"
)

func listSourcesAction(c *cli.Context) error {
	fmt.Print(grab.FormatSourcesTable())
	fmt.Println("\n提示: 多个同类型推送（如多个钉钉）请使用配置文件，在 pusher 下列出多条相同 type。")
	fmt.Println("      运行 watchvuln init-config 可生成带注释的 config.yaml。")
	return nil
}

func initConfigAction(c *cli.Context) error {
	out := c.String("output")
	if out == "" {
		out = "config.yaml"
	}
	if _, err := os.Stat(out); err == nil && !c.Bool("force") {
		return fmt.Errorf("%s already exists, use --force to overwrite", out)
	}
	if err := os.WriteFile(out, []byte(ctrl.InitConfigTemplate), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", out)
	fmt.Println("edit pusher/sources then run: watchvuln -c", out)
	return nil
}

func boardAction(c *cli.Context) error {
	addr := c.String("web-addr")
	if addr == "" {
		addr = "127.0.0.1:8765"
	}
	db := c.String("db-conn")
	if os.Getenv("DB_CONN") != "" {
		db = os.Getenv("DB_CONN")
	}
	config := &ctrl.WatchVulnAppConfig{
		DBConn:  db,
		Version: Version,
		WebAddr: addr,
	}
	client, closeFn, err := ctrl.OpenDatabase(config)
	if err != nil {
		return err
	}
	defer closeFn()

	ctx, cancel := signalCtx()
	defer cancel()
	fmt.Printf("漏洞情报看板: http://%s/ (Ctrl+C 退出)\n", addr)
	return web.NewServer(client, addr).Start(ctx)
}

func loadPusherConfig(c *cli.Context) ([]map[string]string, error) {
	path := c.String("pusher-file")
	if os.Getenv("PUSHER_FILE") != "" {
		path = os.Getenv("PUSHER_FILE")
	}
	if path == "" {
		return initPusher(c)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "read pusher-file")
	}
	var pushers []map[string]string
	if strings.HasSuffix(strings.ToLower(path), ".json") {
		err = json.Unmarshal(data, &pushers)
	} else {
		err = yaml.Unmarshal(data, &pushers)
	}
	if err != nil {
		return nil, errors.Wrap(err, "parse pusher-file")
	}
	if len(pushers) == 0 {
		return nil, fmt.Errorf("pusher-file %s is empty", path)
	}
	log.Infof("loaded %d pusher(s) from %s", len(pushers), path)
	return pushers, nil
}
