package client

import (
	"autoencode-git/util"
	"fmt"
	"os"
	"time"
)

type SvnEncoder struct {
	CommonEncoder
	CheckoutRev int       // 本地版本号,检出
	CommitRev   int       // 本地版本号,提交
	command     *util.Git // Git控制接口s
}

func (s *SvnEncoder) Run(sign chan bool) {
	fmt.Println(time.Now())
	fmt.Println(s.Config)
	fmt.Println(s.ExcludeList)
}

func (s *SvnEncoder) Shutdown(msg string) {
	s.IsDown = true
	s.Logger.Println("autoencode is shutdown...")
	s.Logger.Println(time.Now())
}

func (s *SvnEncoder) Info(msg interface{}) {
	s.Logger.Println(msg)
}

// 发生错误，锁定程序并退出执行，哦，还要发邮件通知管理员
func (s *SvnEncoder) Fatal(msg interface{}) {
	// 创建锁定文件
	if ok, _ := util.IsExist(s.Pwd + util.LOCKFILE); !ok {
		os.Create(s.Pwd + util.LOCKFILE)
	}
}
