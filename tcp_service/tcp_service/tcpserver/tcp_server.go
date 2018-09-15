package tcpserver

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/384782946/im/im_server/service/common/log"
	"github.com/384782946/im/im_server/service/modules/tcp_service/proto"
	"github.com/384782946/im/im_server/service/modules/tcp_service/tcpserver/handler"
)

var (
	connectors    map[string]*net.TCPConn = make(map[string]*net.TCPConn)
	connectorsMux sync.Mutex
	handlers      = handler.NewPBHandler()
)

func handleConnect(conn *net.TCPConn) {
	scanner := bufio.NewScanner(conn)
	scanner.Split(proto.Split)

	for scanner.Scan() {
		pkg, _ := proto.NewPackage()
		err := pkg.UnPacket(scanner.Bytes())
		if err != nil {
			log.Warning("tcp service unpackage error", err)
		} else {
			log.Debug(pkg)

			r, err := handlers.Handle(pkg)
			if err != nil {
				log.Warning("tcp service process error", err)
			}

			conn.Write(r.Packet())
		}
	}
}

func StartTcpService(addr string) {

	address, _ := net.ResolveTCPAddr("tcp4", addr)       //定义一个本机IP和端口。
	var tcpListener, err = net.ListenTCP("tcp", address) //在刚定义好的地址上进监听请求。
	if err != nil {
		fmt.Println("监听出错：", err)
		return
	} else {
		fmt.Println("start tcp service at:", address)
	}
	defer tcpListener.Close()

	fmt.Println("正在等待连接...")

	for {
		var conn, err2 = tcpListener.AcceptTCP() //接受连接。
		if err2 != nil {
			fmt.Println("接受连接失败：", err2)
			continue
		}

		var remoteAddr = conn.RemoteAddr() //获取连接到的对像的IP地址。
		log.Info("接受到一个连接：", remoteAddr)

		go handleConnect(conn)
	}
}
