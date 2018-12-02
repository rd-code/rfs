package ctl

import (
    "fmt"
    "io"
    "net"
    "os"
    "path"
    "rfs/src/protocol"
    "rfs/src/server/context"
    "rfs/src/server/model"
    "strings"
    "sync"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-10 10:37
 **/
type Server struct {
    conn    net.Conn
    ctx     *context.Context
    isAuth  bool //是否认证
    account string
    workDir string
}

func NewServer(conn net.Conn, ctx *context.Context) *Server {
    return &Server{
        conn: conn,
        ctx:  ctx,
    }
}

func (s *Server) Handle() {
    defer func() {
        if msg := recover(); msg != nil {
            fmt.Println("happened panic", msg)
        }
        s.conn.Close()
    }()
    s.handle()

}

func (s *Server) handle() {
    for {
        h, err := protocol.GetHeader(s.conn)
        if err == io.EOF {
            return
        }
        if err != nil {
            fmt.Println("read header failed", err)
            return
        }
        if err := s.do(h); err != nil {
            fmt.Println("handle failed", err)
            return
        }
    }
}

func (s *Server) do(h *protocol.Header) (error) {
    switch h.Command {
    case protocol.CTL_LOGIN:
        return s.auth(h)
    case protocol.CTL_LS:
        return ls(s, h)
    case protocol.CTL_UP:
        return up(s, h)
    case protocol.CTL_DOWN:
        return down(s, h)
    }
    return nil
}

func (s *Server) auth(h *protocol.Header) (err error) {
    var body byte = 1
    defer func() {
        var byteArray []byte
        var h *protocol.Header
        if err != nil {
            s.isAuth = false
            s.responseErr(h, err)
            return
        }

        if body == 0 {
            s.isAuth = true
        } else {
            s.isAuth = false
        }
        h = &protocol.Header{
            Version: protocol.VERSION_0,
            Type:    protocol.TYPE_COMMON,
            Size:    1,
        }
        byteArray = []byte{body}

        if _, err = s.conn.Write(h.Encode()); err != nil {
            fmt.Println("write msg failed", err)
            return
        }
        if _, err = s.conn.Write(byteArray); err != nil {
            fmt.Println("write msg failed", err)
        }
    }()
    sb := &strings.Builder{}
    if _, err = io.CopyN(sb, s.conn, h.Size); err != nil {
        return err
    }
    slice := strings.Split(sb.String(), ";")
    if len(slice) != 2 {
        fmt.Println("the auth msg is not right")
        return nil
    }
    account, pass := slice[0], slice[1]
    var user *model.User
    if user, err = (&model.User{}).GetAsAccount(s.ctx.DB, account); err != nil {
        fmt.Println("query user message from db failed", err)
        return err
    }
    if user.Pass == pass {
        s.workDir = path.Join(s.ctx.Conf.File.Path, user.Account)
        if err = createUserDirectory(path.Join(s.ctx.Conf.File.Path, user.Account)); err != nil {
            fmt.Println("create user directory failed", err)
            return
        }
        body = 0
    }
    return
}

type writerImpl func([]byte) (int, error)

func (f writerImpl) Write(b []byte) (int, error) {
    return f(b)
}

func (s *Server) notAuth(size int64) {
    f := func([]byte) (int, error) {
        return 0, nil
    }
    io.CopyN(writerImpl(f), s.conn, size)
    body := "you need login"
    h := &protocol.Header{
        Version: protocol.VERSION_0,
        Type:    protocol.TYPE_ERR,
        Size:    int64(len(body)),
    }
    byteArray := h.Encode()
    s.conn.Write(byteArray)
    s.conn.Write([]byte(body))
}

func (s *Server) responseErr(h *protocol.Header, err error) {
    if err == nil {
        panic("the error must is nil")
    }

    body := []byte(err.Error())
    header := protocol.Header{
        Version: protocol.VERSION_0,
        Type:    protocol.TYPE_ERR,
        Command: h.Command,
        Size:    int64(len(body)),
    }
    if _, err := s.conn.Write(header.Encode()); err != nil {
        fmt.Println("write header to client failed", err)
        return
    }
    if _, err := s.conn.Write(body); err != nil {
        fmt.Println("write boody to client failed", err)
        return
    }
}

var userMutex sync.Mutex
//创建用户文件夹
func createUserDirectory(path string) error {
    userMutex.Lock()
    defer userMutex.Unlock()
    info, err := os.Stat(path)
    if os.IsNotExist(err) {
        if err = os.Mkdir(path, os.ModePerm); err != nil {
            return err
        }
        info, err = os.Stat(path)
    }
    if err != nil {
        return err
    }
    if !info.IsDir() {
        err = fmt.Errorf("the file:%s already exist", path)
        return err
    }
    return nil

}
