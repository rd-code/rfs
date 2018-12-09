package main

import (
    "fmt"
    "rfs/context"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-12-02 18:06
 **/

func main() {
    context, err := context.Default()
    if err != nil {
        fmt.Println("get default context failed", err)
        return
    }
    defer context.Close()
}
