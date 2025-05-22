package main

import (
	"flag"
	"fmt"

	"imd-seat-be/internal/config"
	"imd-seat-be/internal/handler"
	"imd-seat-be/internal/pkg/response"
	"imd-seat-be/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	httpx.SetErrorHandler(response.ErrHandler)

	DB := sqlx.NewMysql(c.DSN())
	ctx := svc.NewServiceContext(c, DB)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
