package main

import (
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
 * @create 2018-12-02 16:41
 **/
func down(c *Client, fileName string) (err error) {

    if len(fileName) == 0 || len(fileName) > 25 {
        err = errors.New("invalid filename")
        return
    }

    header := &protocol.Header{
        Type:    protocol.TYPE_COMMON,
        Command: protocol.CTL_DOWN,
        Size:    int64(len(fileName)),
    }
    if _, err = c.conn.Write(header.Encode()); err != nil {
        fmt.Println("write header failed", err)
        return
    }
    if _, err = c.conn.Write([]byte(fileName)); err != nil {
        fmt.Println("write body failed", err)
        return
    }

    if header, err = protocol.GetHeader(c.conn); err != nil {
        fmt.Println("parse header failed", err)
        return
    }

    var file *os.File
    if file, err = os.Create(fileName); err != nil {
        fmt.Println("create file failed", err)
        return
    }
    defer file.Close()

    if _, err = io.CopyN(file, c.conn, header.Size); err != nil {
        fmt.Println("write msg to file failed", err)
        return
    }

    return
}

/**
func up(c *Client, fileName string) (err error) {
   if len(fileName) == 0 || len(fileName) > 25 {
       err = errors.New("invalid filename")
       return
   }
   var file *os.File
   if file, err = os.Open(fileName); err != nil {
       return err
   }
   defer file.Close()
   var info os.FileInfo
   if info, err = file.Stat(); err != nil {
       return
   }
   header := &protocol.Header{
       Version: protocol.VERSION_0,
       Type:    protocol.TYPE_COMMON,
       Command: protocol.CTL_UP,
       Size:    info.Size() + 25,
   }
   if _, err = c.conn.Write(header.Encode()); err != nil {
       fmt.Println("write header to protocol failed", err)
       return
   }
   for len(fileName) < 25 {
       fileName = fileName + " "
   }
   if _, err = c.conn.Write([]byte(fileName)); err != nil {
       return
   }
   if _, err = io.Copy(c.conn, file); err != nil {
       return
   }
   if header, err = protocol.GetHeader(c.conn); err != nil {
       fmt.Println("parse header failed", err)
       return
   }

   buffer := &bytes.Buffer{}
   if _, err = io.CopyN(buffer, c.conn, header.Size); err != nil {
       fmt.Println("get body failed", err)
       return err
   }
   if header.Type == protocol.TYPE_ERR {
       return errors.New(buffer.String())
   }
   if buffer.String() == "OK" {
       return
   }
   err = errors.New("unknow body:" + buffer.String())
   return

}

 */
