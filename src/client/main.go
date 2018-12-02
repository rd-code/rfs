package main

import (
    "flag"
    "fmt"
    "net"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-11 23:14
 **/

var host = flag.String("h", "47.105.163.97", "host")
var port = flag.Int("p", 8888, "port")
var account = flag.String("u", "god", "user")
var passwrod = flag.String("a", "123321", "password")

func main() {
    flag.Parse()
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
    if err != nil {
        fmt.Println("create connection failed", err)
        return
    }
    defer conn.Close() //关闭链接

    var client *Client
    if client, err = NewClient(*account, *passwrod, conn); err != nil {
        fmt.Println("create client failed", err)
        return
    }
    if err = client.accept(); err != nil {
        fmt.Println("accept command failed", err)
    }
}
