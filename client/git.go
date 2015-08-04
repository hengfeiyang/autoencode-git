package client

import (
	"autoencode-git/util"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type GitEncoder struct {
	CommonEncoder
	CheckoutRev int       // 本地版本号,检出
	CommitRev   int       // 本地版本号,提交
	command     *util.Git // Git控制接口s
}

func (s *GitEncoder) Run(sign chan bool) {
	fmt.Println(time.Now())
	fmt.Println(s.Config)
	fmt.Println(s.ExcludeList)

	for {
		time.Sleep(time.Second)
		if s.IsDown {
			break
		}
		s.Info("runing..." + time.Now().Format("2006-01-02 15:04:05"))
	}

	sign <- true
}

func (s *GitEncoder) Shutdown(msg string) {
	s.IsDown = true
	s.Logger.Println(msg)
	s.Logger.Println(time.Now())
}

// 普通错误信息记录
func (s *GitEncoder) Info(msg interface{}) {
	s.Logger.Println(msg)
}

// 发生错误，锁定程序并退出执行，哦，还要发邮件通知管理员
func (s *GitEncoder) Fatal(msg interface{}) {
	// 创建锁定文件
	if ok, _ := util.IsExist(s.Pwd + util.LOCKFILE); !ok {
		os.Create(s.Pwd + util.LOCKFILE)
	}

	// 发送邮件
	mailConfig := &util.MailT{
		Addr:  s.Config["mailHost"],
		User:  s.Config["mailUser"],
		Pass:  s.Config["mailPass"],
		From:  s.Config["mailFrom"],
		To:    s.Config["mailTo"],
		Title: s.Config["mailTitle"],
		Type:  "",
	}
	if mailConfig.Title == "" {
		mailConfig.Title = "9466代码机器人"
	}
	mailConfig.Body = "autoencode is locked, please check! " + time.Now().Format("2006-01-02 15:04:05") + "\r\n\r\n" + msg.(string)
	err := util.SendMail(mailConfig)
	if err != nil {
		s.Logger.Println(mailConfig.Body)
		s.Logger.Println("autoencode warning send error: ", err.Error())
	}

	// 关闭程序
	s.Shutdown(msg.(string))
}

// 标记跳过的处理
func (s *GitEncoder) markBlank(rev int) (err error) {
	file := s.Config["commitDir"] + util.MARK_FILE
	if ok, _ := util.IsExist(file); !ok {
		// 创建
		_, err = os.Create(file)
		if err != nil {
			return err
		}
		// 添加到svn
		// err = s.command.Add(s.Config["commitDir"], util.MARK_FILE)
		// if err != nil {
		// 	err = errors.New(util.MARK_FILE + " " + err.Error())
		// 	return
		// }
	}
	// 更新
	err = ioutil.WriteFile(file, []byte(strconv.Itoa(rev)+"\n"), 0644)
	return err
}
