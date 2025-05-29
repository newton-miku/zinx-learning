package utils

import (
	"encoding/json"
	"log/slog"
	"os"
	"zinx/ziface"
)

type GlobalObj struct {

	//服务器端配置
	Server ziface.IServer
	Host   string
	Port   int

	//Zinx配置
	Name          string
	Version       string
	MaxPacketSize uint32
	MaxConn       uint32
}

var GlobalObject *GlobalObj

func (gb *GlobalObj) Reload() error {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		slog.Error("Reload", "err", err)
		return err
	}
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		slog.Error("Reload", "err", err)
		return err
	}
	return nil
}

func init() {
	GlobalObject = &GlobalObj{
		Name:          "ZinxServerApp",
		Version:       "v0.4",
		MaxPacketSize: 4096,
		MaxConn:       1000,
		Host:          "0.0.0.0",
		Port:          8999,
	}
	GlobalObject.Reload()
}
