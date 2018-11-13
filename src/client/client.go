package main

import (
    "errors"
    "fmt"
    "net"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-11 23:20
 **/

type Client struct {
    conn   net.Conn
    prefix string
}

func NewClient(account, password string, conn net.Conn) (*Client, error) {
    ok, err := login(account, password, conn)
    if err != nil {
        fmt.Println("login failed", err)
        return nil, err
    }
    if !ok {
        return nil, errors.New("account or password invalid")
    }
    return &Client{
        conn:   conn,
        prefix: account + ":>",
    }, nil
}

func (c *Client) accept() error {
    for {
        fmt.Print(c.prefix)
        //接收用户输入信息
    }
    return nil
}
