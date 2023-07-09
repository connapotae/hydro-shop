package main

import (
	"os"

	"github.com/connapotae/hydro-shop/config"
	"github.com/connapotae/hydro-shop/modules/servers"
	"github.com/connapotae/hydro-shop/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewServer(cfg, db).Start()
}
