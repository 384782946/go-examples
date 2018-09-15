package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	StartLength = 2
	LenLength   = 4
	TypeLength  = 2
	HeadLen     = StartLength + TypeLength + LenLength
	EndLength   = 2
)

//自定义tcp通信协议
type Package struct {
	StartFlag uint16 `json:"startflag"` //同步标记 0x
	length    uint32 `json:length`      //数据总长度
	Type      uint16 `json:type`        //数据类型
	headLen   uint16 `json:headlen`     //包头长度
	HeadData  string `json:headdata`    //头部保留位
	Data      string `json:data`        //数据
	EndFlag   uint16 `json:endflag`     //结束标志
}

func NewPackage() (*Package, error) {
	return &Package{
		StartFlag: 0x8080,
		length:    12,
		Type:      0,
		headLen:   10,
		HeadData:  "",
		Data:      "",
		EndFlag:   0x8181,
	}, nil
}

func (h *Package) UnPacket(b []byte) error {
	buf := bytes.NewReader(b)
	var startFlog uint16

	binary.Read(buf, binary.BigEndian, &startFlog)
	if h.StartFlag != startFlog {
		return errors.New("start flag can not match")
	}

	binary.Read(buf, binary.BigEndian, &h.length)
	binary.Read(buf, binary.BigEndian, &h.Type)
	binary.Read(buf, binary.BigEndian, &h.headLen)

	if h.headLen > 10 {
		headData := make([]byte, h.headLen-10)
		binary.Read(buf, binary.BigEndian, &headData)
		h.HeadData = string(headData)
	}

	data := make([]byte, h.length-uint32(h.headLen+2))

	binary.Read(buf, binary.BigEndian, &data)
	h.Data = string(data)

	var endFlag uint16
	binary.Read(buf, binary.BigEndian, &endFlag)
	if endFlag != h.EndFlag {
		return errors.New("end flag can not match")
	}

	return nil
}

func (h *Package) Packet() []byte {
	byteBuffer := bytes.NewBuffer([]byte{})
	h.headLen = uint16(len(h.HeadData) + 10)
	h.length = uint32(len(h.Data)) + uint32(h.headLen) + 2
	binary.Write(byteBuffer, binary.BigEndian, h.StartFlag)
	binary.Write(byteBuffer, binary.BigEndian, h.length)
	binary.Write(byteBuffer, binary.BigEndian, h.Type)
	binary.Write(byteBuffer, binary.BigEndian, h.headLen)
	binary.Write(byteBuffer, binary.BigEndian, []byte(h.HeadData))
	binary.Write(byteBuffer, binary.BigEndian, []byte(h.Data))
	binary.Write(byteBuffer, binary.BigEndian, h.EndFlag)
	return byteBuffer.Bytes()
}

type Handler interface {
	Handle(h *Package) (*Package, error)
}

type SimpleHandler struct {
}

func (s *SimpleHandler) Handle(h *Package) (*Package, error) {
	reponse, _ := NewPackage()
	reponse.Data = h.Data
	return reponse, nil
}

func Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !atEOF && len(data) > HeadLen && data[0] == 0x80 && data[1] == 0x80 { // 由于我们定义的数据包头
		length := uint32(0)
		binary.Read(bytes.NewReader(data[StartLength:StartLength+LenLength]), binary.BigEndian, &length) // 读取数据包第3-4字节(int16)=>数据部分长度
		if int(length) <= len(data) {                                      // 如果读取到的数据正文长度+2字节版本号+2字节数据长度不超过读到的数据(实际上就是成功完整的解析出了一个包)
			if data[length-2] == 0x81 && data[length-1] == 0x81 {
				return int(length), data[:int(length)], nil
			}
		}
	}
	return
}
