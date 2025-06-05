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
	// 服务器名称
	Name string
	// 服务器版本
	Version string
	// 最大包大小
	MaxPacketSize uint32
	// 最大连接数
	MaxConn uint32
	// 工作池数量
	WorkerPoolSize uint32
	// 任务队列最大长度
	MaxWorkerTaskQueueSize uint32
	// 日志等级
	LogLevel string
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
		Name:                   "ZinxServerApp",
		Version:                "v0.9",
		MaxPacketSize:          4096,
		MaxConn:                1000,
		Host:                   "0.0.0.0",
		Port:                   8999,
		WorkerPoolSize:         10,
		MaxWorkerTaskQueueSize: 1024,
		LogLevel:               "INFO",
	}
	GlobalObject.Reload()
	var level slog.Level
	switch GlobalObject.LogLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	slog.SetLogLoggerLevel(level)
}
