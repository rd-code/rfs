package context

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "os"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-12-02 18:11
 **/

type DB struct {
    Host     string `yaml:"yaml"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    DataBase string `yaml:"database"`
}

//后期需要tls相关信息
type Listen struct {
    Addr string `yaml:"addr"`
}

//后期增加容量限制等信息
type Storage struct {
    Path string `yaml:"path"`
}

type Conf struct {
    DB      DB      `yaml:"db"`
    Listen  Listen  `yaml:"listen"`
    Storage Storage `yaml:"storage"`
}

func newConf(path string) (*Conf, error) {
    f, err := os.Open(path)
    if err != nil {
        fmt.Printf("open file:%s failed, err:%v", path, err)
        return nil, err
    }
    defer f.Close()
    decoder := yaml.NewDecoder(f)
    res := &Conf{}
    err = decoder.Decode(&res)
    if err != nil {
        fmt.Println("decode file failed", err)
    }
    return res, err
}
