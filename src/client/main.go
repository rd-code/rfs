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

var host = flag.String("h", "127.0.0.1", "host")
var port = flag.Int("p", 5555, "port")
var account = flag.String("u", "", "user")
var passwrod = flag.String("a", "", "password")

func main() {
    flag.Parse()
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
    if err != nil {
        fmt.Println("create connection failed", err)
        return
    }
    defer conn.Close()
    var client *Client
    if client, err = NewClient(*account, *passwrod, conn); err != nil {
        fmt.Println("create client failed", err)
        return
    }
    if err = client.accept(); err != nil {
        fmt.Println("accept command failed", err)
    }
}
