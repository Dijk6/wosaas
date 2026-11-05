package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"wosaas/bootstrap"
	"wosaas/servers"
)

// appConfig 主配置文件。
//
//	go run . -conf ./conf/app.toml
var appConfig = flag.String("conf", "./conf/app.toml", "app config file")

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	bootstrap.MustInit(ctx, *appConfig)

	if err := servers.Start(ctx); err != nil {
		log.Fatalln("server exit:", err)
	}
}
