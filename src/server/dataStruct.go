package server

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// message struct:
//   [header(F0)|body len(4bytes)|command type(1bytes)|body data(n bytes)|tail(FE)]

// command type:
//  1 user defined
//  2 user's message
//  3 query users list
//  4 user select other user

type NetInf interface {
	Parse([]byte) (interface{}, error)
	Build(interface{}) ([]byte, error)
}

type DataBuild interface {
	Build() ([]byte, error)
}

type NetPack struct {
	Header  byte
	Len     []byte
	Command byte
	Body    []byte
	Tail    byte
}

// 解析网络数据包返回解析后匹配的struct
func (np *NetPack) Parse(data []byte) (interface{}, error) {
	if data[0:1][0] != byte(0xF0) {
		return nil, errors.New("data struct parse failed, error data header.")
	}
	np.Header = byte(0xF0)
	if BytesToUInt32(data[1:5]) != uint32(len(data)) {
		return nil, errors.New("data struct parse failed, error data len.")
	}
	copy(np.Len, data[1:5])
	np.Command = data[5:6][0]
	copy(np.Body, data[6:len(data)-1-4-1-1])
	np.Tail = byte(0xFE)

	var inf interface{}
	if np.Command == byte(0x01) {
		inf = new(UserDefined)
	} else if np.Command == byte(0x02) {
		inf = new(UserMessage)
	} else if np.Command == byte(0x03) {
		inf = new(UserList)
	} else if np.Command == byte(0x04) {
		inf = new(UserSelect)
	} else {
		return nil, errors.New("do not recognize the command")
	}

	return inf, nil
}

// 封装网络数据包，将struct封装为字节流
func (np *NetPack) Build(obj interface{}) (result []byte, err error) {
	switch obj.(type) {
	case UserDefined:
		np.Command = byte(0x01)
	case UserMessage:
		np.Command = byte(0x02)
	case UserList:
		np.Command = byte(0x03)
	case UserSelect:
		np.Command = byte(0x04)
	default:
		return nil, errors.New("build network package failed!")
	}
	//对interface{}数据进行断言判断
	val, ok := obj.(DataBuild)
	if !ok {
		return nil, errors.New("parse interface data failed!")
	}
	body, err := val.Build()
	if err != nil {
		return nil, err
	}
	copy(np.Body, body)
	copy(np.Len, UInt32ToBytes(uint32(len(np.Body)+1+4+1+1)))
	np.Header = byte(0xF0)
	np.Tail = byte(0xFE)
	//结构类型序列化
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, np)

	return buf.Bytes(), nil
}

// 用户基本信息结构
type UserDefined struct {
	Id   string // 20bytes
	Name string
}

func (ud *UserDefined) Parse(data []byte) error {
	ud.Id = string(data[0:20])
	ud.Name = string(data[20:])

	return nil
}

func (ud *UserDefined) Build() (data []byte, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ud)
	return buf.Bytes(), nil
}

type UserMessage struct {
	Id   string //20bytes
	Data []byte
}

func (um *UserMessage) Parse(data []byte) error {
	um.Id = string(data[0:20])
	copy(um.Data, data[20:])

	return nil
}

func (um *UserMessage) Build() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, um)

	return buf.Bytes(), nil
}

type UserList struct {
	Ids []string
}

func (ul *UserList) Parse(data []byte) error {
	var idLen int = 20
	if len(data)%idLen != 0 {
		return errors.New("user lists data no normal!")
	}

	arrNum := len(data) / idLen

	ul.Ids = make([]string, arrNum)

	for i, j := 0, 0; i < arrNum; i += idLen {
		ul.Ids[j] = string(data[i : i+idLen])
		j++
	}

	return nil
}

func (ul *UserList) Build() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ul)

	return buf.Bytes(), nil
}

type UserSelect struct {
	Id string
}

func (us *UserSelect) Parse(data []byte) error {
	us.Id = string(data)
	return nil
}

func (us *UserSelect) Build() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, us)

	return buf.Bytes(), nil
}
