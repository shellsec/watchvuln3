package ctrl

import (
	"context"
	"time"

	"entgo.io/ent/dialect"
	entSql "entgo.io/ent/dialect/sql"
	"github.com/pkg/errors"
	"github.com/zema1/watchvuln/ent"
	"github.com/zema1/watchvuln/ent/migrate"
)

// OpenDatabase opens the ent client and ensures schema exists.
func OpenDatabase(config *WatchVulnAppConfig) (*ent.Client, func(), error) {
	config.Init()
	drvName, connStr, err := config.DBConnForEnt()
	if err != nil {
		return nil, nil, err
	}
	drv, err := entSql.Open(drvName, connStr)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed opening connection to db")
	}
	db := drv.DB()
	maxConn := 5
	if drvName == dialect.SQLite {
		maxConn = 1
	}
	db.SetMaxOpenConns(maxConn)
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(maxConn)
	client := ent.NewClient(ent.Driver(drv))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if err := client.Schema.Create(ctx, migrate.WithDropIndex(true), migrate.WithDropColumn(true)); err != nil {
		_ = client.Close()
		return nil, nil, errors.Wrap(err, "failed creating schema resources")
	}
	closeFn := func() { _ = client.Close() }
	return client, closeFn, nil
}
