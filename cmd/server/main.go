// Command server 启动 mqx-lite Web 服务：K 线展示 + jqs 指标脚本引擎。
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/chainball/mqx-lite/internal/tdx"
	"github.com/chainball/mqx-lite/internal/web"
)

func main() {
	p, err := tdx.New()
	if err != nil {
		log.Fatalf("connect tdx: %v", err)
	}
	defer p.Close()

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("mqx-lite listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, web.New(p).Handler()))
}
