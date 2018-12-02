package main

import (
    "bufio"
    "errors"
    "fmt"
    "net"
    "os"
    "strings"
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
        value, err := input()
        if err != nil {
            fmt.Println("get input failed", err)
            return err
        }
        if len(value) == 0 {
            continue
        }
        if len(value) == 1 && value[0] == "ls" {
            ls(c)
        } else if len(value) == 2 && value[0] == "up" {
            up(c, value[1])
        }

    }
    return nil
}

//获取用户输入信息
func input() ([]string, error) {
    inputReader := bufio.NewReader(os.Stdin)
    value, err := inputReader.ReadString('\n')
    if err != nil {
        return nil, err
    }
    value = value[:len(value)-1]
    items := strings.Split(value, " ")
    var slice []string
    for _, item := range items {
        if len(item) != 0 {
            slice = append(slice, item)
        }

    }
    return slice, nil
}
