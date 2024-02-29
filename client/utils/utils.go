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

func (this *Transfer) ReadPkg() (mesRead message.Message, readErr error) {
	//buffer := make([]byte, 8096)
	//conn.Read只有在双方conn没有被关闭的情况下，才会阻塞等待
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
	//把客户端发送的byte类型切片表示的pkgLen->uint32类型的pkgLen
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
		readErr = errors.New("json.Unmarshal err")
		return
	}
	return mesRead, readErr
}

func (this *Transfer) WritePkg(data []byte) (responseErr error) {
	//发送信息长度
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buffer [4]byte
	binary.BigEndian.PutUint32(this.Buffer[:4], pkgLen) //长度只用了4个字节
	n, responseErr := this.Conn.Write(this.Buffer[:4])
	if n != 4 || responseErr != nil {
		responseErr = errors.New("conn.Write(len) failed")
		return
	}

	//发送信息内容
	n, responseErr = this.Conn.Write(data)
	if n != int(pkgLen) || responseErr != nil {
		responseErr = errors.New("conn.Write(data) failed")
		return
	}
	return responseErr
}
