package proto_test

import (
	"testing"

	"github.com/384782946/im/im_server/service/modules/tcp_service/proto"
)

func TestPacakge(t *testing.T) {
	p, _ := proto.NewPackage()
	p.Data = "123123abcABC"
	p.HeadData = "version=V1.0"

	pack := p.Packet()

	t.Log(p.Packet())

	p1, _ := proto.NewPackage()
	p1.UnPacket(pack)

	t.Log(p1)
}

func TestSplit(t *testing.T) {
	p, _ := proto.NewPackage()
	p.Data = "123123abcABC"
	p.HeadData = "version=V1.0"

	pack := p.Packet()
	pack = append(pack, []byte("haha")...)
	t.Log(pack)

	a, token, _ := proto.Split(pack, false)
	t.Log(a, token)
}
