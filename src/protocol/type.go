package protocol

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "io"
    "sync"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-11 19:30
 **/

const (
    HEADER_SIZE = 10
)
const (
    VERSION_0 = iota
)
const (
    TYPE_COMMON = iota //传递普通消息
    TYPE_ERR           //传递错误信息
    TYPE_FILE          //传递文件
)

const (
    CTL_LOGIN = iota //登录命令
    CTL_LS
    CTL_EXIT
    CTL_UP  //upload
    CTL_DOWN
)

type Header struct {
    Version uint8 //版本
    Type    uint8 //类型
    Command uint8 //命令
    Size    int64 //body长度
}

func (h *Header) Encode() []byte {
    buffer := &bytes.Buffer{}
    vt := h.Version
    vt = vt << 4
    vt = vt | h.Type
    buffer.WriteByte(vt)
    buffer.WriteByte(h.Command)
    binary.Write(buffer, endian, h.Size)
    return buffer.Bytes()
}

var endian = binary.BigEndian

//解析请求头
func GetHeader(writer io.Reader) (*Header, error) {
    buffer := &bytes.Buffer{}
    if _, err := io.CopyN(buffer, writer, HEADER_SIZE); err != nil {
        return nil, err
    }
    byteArraay := buffer.Bytes()

    value := byteArraay[0]

    h := &Header{
        Type: value & 0xF,
    }
    value = value >> 4
    h.Version = value
    h.Command = byteArraay[1]
    if err := binary.Read(bytes.NewReader(byteArraay[2:]), endian, &h.Size); err != nil {
        fmt.Println("parse int from bytes failed", err)
        return nil, err
    }
    return h, nil
}

var userMutex sync.Mutex
