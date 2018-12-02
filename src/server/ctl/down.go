package ctl

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "os"
    "rfs/src/protocol"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-12-02 15:54
 **/

func down(server *Server, header *protocol.Header) error {
    if !server.isAuth {
        server.notAuth(header.Size)
        return nil
    }
    resHeader := &protocol.Header{
        Version: protocol.VERSION_0,
        Command: protocol.CTL_UP,
    }

    buffer := &bytes.Buffer{}
    if _, err := io.CopyN(buffer, server.conn, resHeader.Size); err != nil {
        body := []byte(err.Error())
        resHeader.Type = protocol.TYPE_ERR
        resHeader.Size = int64(len(body))
        if err := resp(server.conn, header, body); err != nil {
            fmt.Println("ersp failed", err)
        }
        return err
    }

    path := buffer.String()
    file, err := os.Open(path)
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

    var info os.FileInfo
    if info, err = file.Stat(); err != nil {
        body := []byte(err.Error())
        resHeader.Type = protocol.TYPE_ERR
        resHeader.Size = int64(len(body))
        if err := resp(server.conn, header, body); err != nil {
            fmt.Println("ersp failed", err)
        }
        return err
    }

    if info.IsDir() {
        err = errors.New("the file is directory")
        body := []byte(err.Error())
        resHeader.Type = protocol.TYPE_ERR
        resHeader.Size = int64(len(body))
        if err := resp(server.conn, header, body); err != nil {
            fmt.Println("ersp failed", err)
        }
        return err
    }
    resHeader.Type = protocol.TYPE_FILE
    resHeader.Size = info.Size()
    if _, err = server.conn.Write(resHeader.Encode()); err != nil {
        fmt.Println("write header failed", err)
        return err
    }
    if _, err = io.Copy(server.conn, file); err != nil {
        fmt.Println("write to conn failed", err)
        return err
    }

    return nil
}
