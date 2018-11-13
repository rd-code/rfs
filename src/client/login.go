package main

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "net"
    "rfs/src/protocol"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-11 23:22
 **/

func login(account, pass string, conn net.Conn) (bool, error) {
    body := fmt.Sprintf("%s;%d", account, pass)
    header := &protocol.Header{
        Version: protocol.VERSION_0,
        Type:    protocol.TYPE_COMMON,
        Command: protocol.CTL_LOGIN,
        Size:    int64(len(body)),
    }
    _, err := conn.Write(header.Encode())
    if err != nil {
        return false, err
    }
    if _, err = conn.Write([]byte(body)); err != nil {
        return false, err
    }
    if header, err = protocol.GetHeader(conn); err != nil {
        return false, err
    }

    buffer := &bytes.Buffer{}
    if _, err = io.CopyN(buffer, conn, header.Size); err != nil {
        return false, err
    }
    if header.Type == protocol.TYPE_ERR {
        return false, errors.New(buffer.String())
    }

    slice := buffer.Bytes()
    if len(slice) != 1 {
        return false, errors.New("cannot handle data")
    }
    return slice[0] == 0, nil
}
