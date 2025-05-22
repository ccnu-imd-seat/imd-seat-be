package main

import (
	"context"
	"flag"
	"fmt"

	"imd-seat-be/internal/config"
	"imd-seat-be/internal/handler"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/task"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/seat-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	DB := sqlx.NewMysql(c.DSN())
	ctx := svc.NewServiceContext(c, DB)
	handler.RegisterHandlers(server, ctx)
	//执行定时任务
	go task.RegisterTasks(context.Background(), ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
