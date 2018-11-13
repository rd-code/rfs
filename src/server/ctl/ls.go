package ctl

import (
    "fmt"
    "github.com/rd-code/common/filter"
    _ "github.com/rd-code/common/filter"
    "io/ioutil"
    "rfs/src/protocol"
    "strings"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-11 19:47
 **/

//ls命令，列出当前文件夹下面所有信息
func ls(server *Server, header *protocol.Header) error {
    if !server.isAuth {
        server.notAuth(header.Size)
        return nil
    }
    infos, err := ioutil.ReadDir(server.workDir)
    if err != nil {
        server.responseErr(header, err)
        fmt.Println("read user directory failed", err)
        return err
    }

    names := filter.MapString(infos, func(index int) string {
        return infos[index].Name()
    })

    body := []byte(strings.Join(names, "\n"))

    h := protocol.Header{
        Version: protocol.VERSION_0,
        Type:    protocol.TYPE_COMMON,
        Command: protocol.CTL_LS,
        Size:    int64(len(body)),
    }

    if _, err = server.conn.Write(h.Encode()); err != nil {
        fmt.Println("write header failed", err)
        return err
    }
    if _, err = server.conn.Write(body); err != nil {
        fmt.Println("write header failed", err)
    }
    return err
}
