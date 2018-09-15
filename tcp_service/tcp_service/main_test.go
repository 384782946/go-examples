package main_test

import (
	"net"
	"testing"

	"github.com/384782946/im/im_server/service/modules/tcp_service/proto"
)

func TestTcpService(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		t.Fatal("connect error", err)
	}

	pack, _ := proto.NewPackage()
	pack.Type = 0
	pack.HeadData = "version=1.0"
	pack.Data = `{"id":123123,"msg":"hello world"}`

	for i := 0; i < 10; i++ {
		len, err := conn.Write(pack.Packet())

		if err != nil {
			t.Error(err)
		} else {
			t.Log("tcp send data len:", len)

			buf := make([]byte, 1024)
			len, _ := conn.Read(buf)

			recPack, _ := proto.NewPackage()
			recPack.UnPacket(buf[:len])
			t.Log("tcp read:", recPack)
		}
	}

	conn.Close()
}
