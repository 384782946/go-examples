package main

import (
	"os"
	"sync"

	"github.com/384782946/im/im_server/service/common/log"
	"github.com/384782946/im/im_server/service/common/rpc/server"
	"github.com/384782946/im/im_server/service/modules/tcp_service/tcpserver"
)

var (
	wg       sync.WaitGroup
	basePath = "/v1/rpcx_tcp_server"
	addr     = "127.0.0.1:8979"
)

func main() {
	addr, b := os.LookupEnv("TCP_SERVER_ADDR")
	if b != true {
		addr = "localhost:8080"
	}
	wg.Add(1)

	go func() {
		tcpserver.StartTcpService(addr)
		wg.Done()
	}()

	go func() {
		s := server.NewServer(basePath, addr)
		s.Run()
	}()

	wg.Wait()

	log.Info("tcp service exited")
}
