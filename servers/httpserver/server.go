package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"wosaas/config"
)

// NewServer 创建主 HTTP Server。
func NewServer(addr string, cfg config.ServerItem) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      NewRouter(cfg),
		ReadTimeout:  toDuration(cfg.ReadTimeout),
		WriteTimeout: toDuration(cfg.WriteTimeout),
		IdleTimeout:  toDuration(cfg.IdleTimeout),
	}
}

// NewAdminServer 创建管理 HTTP Server。
// 管理接口包括健康检查、指标占位等，生产环境不应暴露到外网。
func NewAdminServer(addr string, cfg config.ServerItem) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintln(w, "# wosaas metrics placeholder")
	})

	return &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  toDuration(cfg.ReadTimeout),
		WriteTimeout: toDuration(cfg.WriteTimeout),
		IdleTimeout:  toDuration(cfg.IdleTimeout),
	}
}

func toDuration(milliseconds int) time.Duration {
	if milliseconds <= 0 {
		return 0
	}
	return time.Duration(milliseconds) * time.Millisecond
}
