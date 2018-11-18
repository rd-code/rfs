package main

import (
    "fmt"
    "net"
    "rfs/src/server/context"
    "rfs/src/server/ctl"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-13 22:08
 **/

func main() {
    ctx, err := context.Load()
    if err != nil {
        fmt.Println("load context failed", err)
        return
    }
    if err = listen(ctx); err != nil {
        fmt.Println("listen failed", err)
    }

}

func listen(ctx *context.Context) error {
    listener, err := net.Listen("tcp", ctx.Conf.Listen.Addr)
    if err != nil {
        return err
    }
    fmt.Println("listen to: ", ctx.Conf.Listen.Addr)
    for {
        var conn net.Conn
        if conn, err = listener.Accept(); err != nil {
            fmt.Println("accept connection failed", err)
            continue
        }
        fmt.Println("accept new connection")
        s := ctl.NewServer(conn, ctx)
        go s.Handle()
    }
    return nil
}
