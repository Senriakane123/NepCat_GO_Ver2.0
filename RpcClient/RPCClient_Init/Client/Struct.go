package Client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

type CLient struct {
	ServiceID     int64    `json:"service_id"`       //客户端注册服务唯一ID
	Conn          net.Conn `json:"conn"`             //客户端连接
	ServerType    string   `json:"server_type"`      //客户端自己提供的服务类型
	MsgSnCodeList []int    `json:"msg_sn_code_list"` //存储客户端发送消息的SN码用于后续消息返回进行客户端回复匹配
	RemoteAddr    string   `json:"remote_addr"`      //连接客户端远端IP
}

var Init_Client Cient_Init

type Cient_Init struct {
	Clients sync.Map // 存储注册的服务: map[string]*Service
	//Listener net.Listener `json:"listener"` // 监听器
}

type VRTSProxyProtocolHeader struct {
	version int32
	size    int32
	msgType uint32 //消息类型
	msgSn   uint32 //消息SN号码
}

func (obj *VRTSProxyProtocolHeader) Parse(buffer []byte) int {
	bufio := bytes.NewReader(buffer)
	binary.Read(bufio, binary.BigEndian, &obj.version)
	binary.Read(bufio, binary.BigEndian, &obj.size)
	binary.Read(bufio, binary.BigEndian, &obj.msgType)
	binary.Read(bufio, binary.BigEndian, &obj.msgSn)
	return 0
}

func (obj *VRTSProxyProtocolHeader) Package(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, obj.version)
	binary.Write(buf, binary.BigEndian, obj.size)
	binary.Write(buf, binary.BigEndian, obj.msgType)
	binary.Write(buf, binary.BigEndian, obj.msgSn)
}

func (obj *VRTSProxyProtocolHeader) Size() int {
	return 16
}

type VRTSProxyRPCHeader struct {
	Method     int32
	ServerType string
	Host       string
	Caller     int32
	Service    int32
	BodySize   int32
}

func (obj *VRTSProxyRPCHeader) Size() int32 {
	return int32(4 + 4 + len(obj.ServerType) + 4 + len(obj.Host) + 12)
}

func (obj *VRTSProxyRPCHeader) Package(buffer *bytes.Buffer) int {
	binary.Write(buffer, binary.BigEndian, obj.Method)

	// ServerType
	binary.Write(buffer, binary.BigEndian, int32(len(obj.ServerType)))
	binary.Write(buffer, binary.BigEndian, []byte(obj.ServerType))

	// Host
	binary.Write(buffer, binary.BigEndian, int32(len(obj.Host)))
	binary.Write(buffer, binary.BigEndian, []byte(obj.Host))

	// 后续字段
	binary.Write(buffer, binary.BigEndian, obj.Caller)
	binary.Write(buffer, binary.BigEndian, obj.Service)
	binary.Write(buffer, binary.BigEndian, obj.BodySize)

	// 返回总长度
	return 4 + 4 + len(obj.ServerType) + 4 + len(obj.Host) + 12
}

func (obj *VRTSProxyRPCHeader) Parse(buf []byte) int32 {
	var retVal int32
	bufio := bytes.NewReader(buf)

	var strLen int32

	// Method
	if err := binary.Read(bufio, binary.BigEndian, &obj.Method); err != nil {
		fmt.Println("Parse Method Failed:", err)
		return -1
	}
	retVal += 4

	// ServerType
	if err := binary.Read(bufio, binary.BigEndian, &strLen); err != nil {
		fmt.Println("Parse ServerType length Failed:", err)
		return -1
	}
	retVal += 4

	if strLen > 0 {
		tmp := make([]byte, strLen)
		if err := binary.Read(bufio, binary.BigEndian, tmp); err != nil {
			fmt.Println("Parse ServerType Failed:", err)
			return -1
		}
		obj.ServerType = string(tmp)
		retVal += strLen
	}

	// Host
	if err := binary.Read(bufio, binary.BigEndian, &strLen); err != nil {
		fmt.Println("Parse Host length Failed:", err)
		return -1
	}
	retVal += 4

	if strLen > 0 {
		tmp := make([]byte, strLen)
		if err := binary.Read(bufio, binary.BigEndian, tmp); err != nil {
			fmt.Println("Parse Host Failed:", err)
			return -1
		}
		obj.Host = string(tmp)
		retVal += strLen
	}

	// Caller
	if err := binary.Read(bufio, binary.BigEndian, &obj.Caller); err != nil {
		fmt.Println("Parse Caller Failed:", err)
		return -1
	}
	retVal += 4

	// Service
	if err := binary.Read(bufio, binary.BigEndian, &obj.Service); err != nil {
		fmt.Println("Parse Service Failed:", err)
		return -1
	}
	retVal += 4

	// BodySize
	if err := binary.Read(bufio, binary.BigEndian, &obj.BodySize); err != nil {
		fmt.Println("Parse BodySize Failed:", err)
		return -1
	}
	retVal += 4

	return retVal
}
