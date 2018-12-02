package ctl

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "path"
    "rfs/src/protocol"
    "strings"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-18 21:13
 **/

//up 命令，上传操作
func up(server *Server, header *protocol.Header) error {
    if !server.isAuth {
        server.notAuth(header.Size)
        return nil
    }
    resHeader := &protocol.Header{
        Version: protocol.VERSION_0,
        Command: protocol.CTL_UP,
    }
    buffer := &bytes.Buffer{}
    _, err := io.CopyN(buffer, server.conn, 25)
    if err != nil {
        body := []byte(err.Error())
        resHeader.Type = protocol.TYPE_ERR
        resHeader.Size = int64(len(body))
        if err := resp(server.conn, header, body); err != nil {
            fmt.Println("ersp failed", err)
        }
        return err
    }

    str := buffer.String()
    str = strings.TrimSpace(str)

    fileName := path.Join(server.workDir, str)
    _, err = os.Stat(fileName)
    if !(err != nil && os.IsNotExist(err)) {
        body := []byte("the file is invalid")
        resHeader.Type = protocol.TYPE_ERR
        resHeader.Size = int64(len(body))
        if err := resp(server.conn, header, body); err != nil {
            fmt.Println("ersp failed", err)
        }
        return err
    }
    file, err := os.Create(fileName)
    if err != nil {
        body := []byte(err.Error())
        resHeader.Type = protocol.TYPE_ERR
        resHeader.Size = int64(len(body))
        if err := resp(server.conn, header, body); err != nil {
            fmt.Println("ersp failed", err)
        }
        return err
    }
    defer file.Close()
    if _, err = io.CopyN(file, server.conn, header.Size-25); err != nil {
        body := []byte(err.Error())
        resHeader.Type = protocol.TYPE_ERR
        resHeader.Size = int64(len(body))
        if err := resp(server.conn, header, body); err != nil {
            fmt.Println("ersp failed", err)
        }
        return err
    }
    body := []byte("OK")
    resHeader.Type = protocol.TYPE_COMMON
    resHeader.Size = int64(len(body))
    if err := resp(server.conn, resHeader, body); err != nil {
        fmt.Println("ersp failed", err)
    }
    return err
}
func resp(writer io.Writer, header *protocol.Header, array []byte) error {
    byteArray := header.Encode()
    _, err := writer.Write(byteArray)
    if err != nil {
        return err
    }
    _, err = writer.Write(array)
    return err
}
