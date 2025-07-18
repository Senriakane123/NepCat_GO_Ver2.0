package Server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"unsafe"
)

type Server struct {
	ServiceID     int64    `json:"service_id"`
	Conn          net.Conn `json:"conn"`
	ServerType    string   `json:"server_type"`
	MsgSnCodeList []int    `json:"msg_sn_code_list"`
	RemoteAddr    string   `json:"remote_addr"`
}

var Init_Server Service_Init

type Service_Init struct {
	Services sync.Map     // 存储注册的服务: map[string]*Service
	Listener net.Listener `json:"listener"` // 监听器
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

type VRTSProxyRegService struct {
	serviceType int32
	name        string
}

func (obj *VRTSProxyRegService) size() int32 {
	return int32(4*2 + len(obj.name))
}

func (obj *VRTSProxyRegService) Package(buf *bytes.Buffer) int {
	size := 8 + len(obj.name)
	err := binary.Write(buf, binary.BigEndian, obj.serviceType)
	if err != nil {
		fmt.Println(err.Error())
	}
	//binary.Write(buf, binary.BigEndian, obj.serviceType)
	err = binary.Write(buf, binary.BigEndian, int32(len(obj.name)))
	if err != nil {
		fmt.Println(err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, []byte(obj.name))
	if err != nil {
		fmt.Println(err.Error())
	}
	return size

}

func (obj *VRTSProxyRegService) Parse(buf []byte) int32 {
	var retVal int32
	bufio := bytes.NewReader(buf)

	var nameLen int32
	retVal = 0

	// 解析 serviceType
	if err := binary.Read(bufio, binary.BigEndian, &obj.serviceType); err != nil {
		fmt.Println("Parse serviceType failed:", err)
		return -1
	}
	retVal += int32(unsafe.Sizeof(obj.serviceType))

	// 解析 name 的长度
	if err := binary.Read(bufio, binary.BigEndian, &nameLen); err != nil {
		fmt.Println("Parse name length failed:", err)
		return -1
	}
	retVal += int32(unsafe.Sizeof(nameLen))

	// 读取 name 本体
	if nameLen > 0 {
		temp := make([]byte, nameLen)
		if err := binary.Read(bufio, binary.BigEndian, temp); err != nil {
			fmt.Println("Parse name failed:", err)
			return -1
		}
		obj.name = string(temp)
		retVal += nameLen
	}

	return retVal
}
