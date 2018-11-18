package context

import (
    "database/sql"
    "errors"
    "flag"
    "fmt"
    _ "github.com/lib/pq"
    "gopkg.in/yaml.v2"
    "os"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-13 22:09
 **/

type Conf struct {
    DB     DB     `yaml:"db"`
    Listen Listen `yaml:"listen"`
    File   File   `yaml:"file"`
}

func (*Conf) nilErr() error {
    return errors.New("the conf is nil")
}

func NewConf(path string) (*Conf, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    decoder := yaml.NewDecoder(f)
    res := &Conf{}
    err = decoder.Decode(&res)
    return res, err

}

type DB struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    DataBase string `yaml:"database"`
}

type Listen struct {
    Addr string `yaml:"addr"`
}

type File struct {
    Path string `yaml:"path"`
}

var path = flag.String("conf", "src/server/conf/conf.yml", "the path of configuration file")

type Context struct {
    Conf *Conf
    DB   *sql.DB
}

func (c *Context) loadConf(path string) error {
    conf, err := NewConf(path)
    if err != nil {
        return err
    }
    c.Conf = conf
    return nil
}

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

func (c *Context) loadModule() error {
    if err := c.loadDB(); err != nil {
        return err
    }
    return nil
}

func (c *Context) createDirectory() error {
    info, err := os.Stat(c.Conf.File.Path)
    if os.IsNotExist(err) {
        if err = os.Mkdir(c.Conf.File.Path, os.ModePerm); err != nil {
            return err
        }
        info, err = os.Stat(c.Conf.File.Path)
    }
    if err != nil {
        return err
    }
    if !info.IsDir() {
        err = fmt.Errorf("the path:%s must is directory", c.Conf.File.Path)
        return err
    }
    return nil
}

func Load() (*Context, error) {
    ctx := &Context{}
    if err := ctx.loadConf(*path); err != nil {
        return nil, err
    }
    if err := ctx.loadModule(); err != nil {
        return nil, err
    }

    if err := ctx.createDirectory(); err != nil {
        return nil, err
    }
    return ctx, nil
}
