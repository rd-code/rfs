package main

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "rfs/src/protocol"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-18 10:38
 **/

func ls(cli *Client) (err error) {
    fmt.Println("-----")
    header := &protocol.Header{
        Version: protocol.VERSION_0,
        Type:    protocol.TYPE_COMMON,
        Command: protocol.CTL_LS,
        Size:    0,
    }

    if _, err = cli.conn.Write(header.Encode()); err != nil {
        fmt.Println("write header to protocol failed", err)
        return
    }
    if header, err = protocol.GetHeader(cli.conn); err != nil {
        fmt.Println("parse header failed", err)
        return
    }

    buffer := &bytes.Buffer{}
    if _, err = io.CopyN(buffer, cli.conn, header.Size); err != nil {
        fmt.Println("get body failed", err)
        return err
    }
    if header.Type == protocol.TYPE_ERR {
        return errors.New(buffer.String())
    }
    fmt.Println("==========")
    fmt.Println(cli.prefix, buffer.String())

    return
}
