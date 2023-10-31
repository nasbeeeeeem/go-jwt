package main

import (
	"flag"

	"go-jwt/pkg/interfaces/api/server"
)

var (
	addr string
)

// main関数の前に処理される
func init() {
	// コマンドライン引数の設定
	flag.StringVar(&addr, "addr", ":8080", "tcp host:port to connect")
	flag.Parse()
}

func main() {
	server.Server(addr)
}
