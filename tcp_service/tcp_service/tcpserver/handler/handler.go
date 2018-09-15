package handler

import (
	"github.com/384782946/im/im_server/service/common/log"
	"github.com/384782946/im/im_server/service/modules/tcp_service/proto"
	"github.com/384782946/im/im_server/service/protol/tcp_service"
	protobuf "github.com/golang/protobuf/proto"

	"errors"
)

type Handler interface {
	Handle(*proto.Package) (proto.Package, error)
}

type PBHandler struct {
}

func NewPBHandler() *PBHandler {
	return &PBHandler{}
}

func (pbh *PBHandler) Handle(p *proto.Package) (r *proto.Package, e error) {
	if p == nil {
		return nil, errors.New("empty package:")
	}

	resType, response, err := pbh.deal(p.HeadData, p.Type, p.Data)
	if err != nil {

	}

	resPackage, _ := proto.NewPackage()
	resPackage.Type = resType
	resPackage.HeadData = p.HeadData
	resPackage.Data = string(response)

	return resPackage, nil

}

func (pbh *PBHandler) deal(head string, t uint16, data string) (uint16, []byte, error) {
	log.Debug("deal with", head, t, data)
	if t == uint16(tcpprotol.RequestType_REQ_HEART_BEAT) {
		return t, []byte(data), nil
		req := &tcpprotol.Request{}

		err := protobuf.Unmarshal([]byte(data), req)
		if err != nil {
			log.Error("tcp unmarshal error:", err)
		}
		pData, err := protobuf.Marshal(req)
		return uint16(req.GetType()), pData, err
	}

	return 0, []byte(""), nil
}
