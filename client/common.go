package client

import (
	"log"
)

type Encoder interface {
	Run(chan bool)
	Info(interface{})
	Fatal(interface{})
	Shutdown(string)
}

type CommonEncoder struct {
	Config      map[string]string // 配置文件
	ExcludeList []string          // 文件处理，排除列表
	Pwd         string            // 当前工作目录
	Logger      *log.Logger       // 日志处理接口
	IsDown      bool              // 是否关闭标识
}
