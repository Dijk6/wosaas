package servers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"wosaas/config"
	"wosaas/library/resource"
	"wosaas/servers/httpserver"
)

// Start 启动所有 server，并在收到退出信号后优雅关闭。
func Start(ctx context.Context) error {
	app := config.App()
	serverCfg, err := config.LoadServer(config.Path("server.toml"))
	if err != nil {
		return fmt.Errorf("load server.toml: %w", err)
	}

	mainAddr := app.ExpandListen(serverCfg.HTTPServer.Listen)
	adminAddr := app.ExpandListen(serverCfg.AdminServer.Listen)

	mainServer := httpserver.NewServer(mainAddr, serverCfg.HTTPServer)
	adminServer := httpserver.NewAdminServer(adminAddr, serverCfg.AdminServer)

	errCh := make(chan error, 2)

	go func() {
		log.Printf("HTTP server listening on %s", mainAddr)
		if err := mainServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server: %w", err)
		}
	}()

	go func() {
		log.Printf("admin server listening on %s", adminAddr)
		if err := adminServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("admin server: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-errCh:
		log.Printf("server error: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = mainServer.Shutdown(shutdownCtx)
	_ = adminServer.Shutdown(shutdownCtx)

	if resource.RedisClient != nil {
		_ = resource.RedisClient.Close()
	}
	resource.CloseDB()

	log.Println("servers stopped")
	return nil
}
