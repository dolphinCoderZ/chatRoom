package utils

import (
	"chatRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
)

type Transfer struct {
	Conn   net.Conn
	Buffer [8096]byte
}

// 一层反序列化->消息包
func (this *Transfer) ReadPkg() (mesRead message.Message, readErr error) {
	//buffer := make([]byte, 8096)
	//conn.Read只有在双方conn没有被关闭的情况下，才会阻塞等待
	//当客户端conn关闭，服务器端conn.Read就会返回io.EOF

	//读消息长度
	_, readErr = this.Conn.Read(this.Buffer[:4])
	if readErr != nil {
		if readErr != io.EOF {
			readErr = errors.New("read pkgLen error")
		}
		return
	}

	//读消息内容
	var pkgLen uint32
	//byte类型切片表示的信息包长度->uint32类型的信息包长度
	pkgLen = binary.BigEndian.Uint32(this.Buffer[:4])
	n, readErr := this.Conn.Read(this.Buffer[:pkgLen])
	if n != int(pkgLen) || readErr != nil {
		if readErr != io.EOF {
			readErr = errors.New("read pkgData error")
		}
		return
	}

	//反序列化
	readErr = json.Unmarshal(this.Buffer[:pkgLen], &mesRead)
	if readErr != nil {
		readErr = errors.New("ReadPkg json.Unmarshal err")
		return
	}
	return mesRead, readErr
}

func (this *Transfer) WritePkg(data []byte) (responseErr error) {
	//发送信息长度
	//n, err := conn.Write([]byte(len(mesData)))
	//mes切片的长度要用byte类型切片来表示
	var pkgLen uint32
	//用uint32来表示长度
	pkgLen = uint32(len(data))
	//将uint32表示的长度用byte类型的切片来表示
	binary.BigEndian.PutUint32(this.Buffer[:4], pkgLen) //长度只用了4个字节
	n, responseErr := this.Conn.Write(this.Buffer[:4])
	if n != 4 || responseErr != nil {
		responseErr = errors.New("WritePkg conn.Write(len) failed")
		return
	}

	//发送信息内容
	n, responseErr = this.Conn.Write(data)
	if n != int(pkgLen) || responseErr != nil {
		responseErr = errors.New("WritePkg conn.Write(data) failed")
		return
	}
	return
}
