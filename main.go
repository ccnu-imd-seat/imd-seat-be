package main

import (
	"flag"
	"fmt"

	"imd-seat-be/internal/config"
	"imd-seat-be/internal/handler"
	"imd-seat-be/internal/pkg/resp"
	"imd-seat-be/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/seat-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	httpx.SetErrorHandler(resp.ErrHandler)

	DB := sqlx.NewMysql(c.DSN())
	ctx := svc.NewServiceContext(c, DB)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
