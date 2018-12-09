package context

import (
    "database/sql"
    "flag"
    "fmt"
    _ "github.com/lib/pq"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-12-02 18:22
 **/

type Context struct {
    Conf *Conf
    DB   *sql.DB
}

//根据配置文件路径生成context
func NewContext(configurationFilePath string) (*Context, error) {
    conf, err := newConf(configurationFilePath)
    if err != nil {
        return nil, err
    }
    res := &Context{
        Conf: conf,
    }
    if err = res.loadModule(); err != nil {
        return nil, err
    }
    return res, nil
}

//加载context各个模块
func (c *Context) loadModule() error {
    if err := c.loadDB(); err != nil {
        return err
    }

    return nil
}

//加载db模块
func (c *Context) loadDB() error {
    db, err := sql.Open("postgres", fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        c.Conf.DB.Host, c.Conf.DB.Port, c.Conf.DB.User, c.Conf.DB.Password, c.Conf.DB.DataBase))

    if err != nil {
        return err
    }
    c.DB = db
    return nil
}

func (c *Context) Close() error {
    return c.DB.Close()
}

//默认配置文件路径
var configurationFilePath = flag.String("conf", "./conf.yml", "the path of configuration file")

//生产默认Context信息
//配置文件路径使用默认路径
func Default() (*Context, error) {
    return NewContext(*configurationFilePath)
}
