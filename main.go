package main

import (
	"autoencode-git/client"
	"autoencode-git/util"
	"bufio"
	"github.com/9466/daemon"
	"github.com/9466/goconfig"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

var configField = [21]string{
	"versionType", "logFile",
	"svnBin", "svnReposCheckout", "svnReposCommit", "svnReposBaseDir",
	"gitBin", "gitReposCheckout", "gitReposCommit", "gitBranches",
	"checkoutDir", "commitDir",
	"encodeExt", "encodeBinPHP", "encodeBinUglifyJS",
	"mailHost", "mailUser", "mailPass", "mailFrom", "mailTo", "mailTitle",
}

const (
	CONFIG_FILE  string = "/conf/common.conf"  // 通用配置文件
	EXCLUDE_FILE string = "/conf/exclude.conf" // 过滤配置文件
)

func main() {
	var err error

	// 先设定工作目录
	var dir string
	dir, err = util.GetDir()
	if err != nil {
		log.Fatalln("My God, GetDir Fatal!")
	}
	dir = path.Dir(dir)

	// 检测是否存在错误锁定
	if ok, _ := util.IsExist(dir + util.LOCKFILE); ok {
		log.Fatalln("Sorry, encoder had locked, please check error!")
	}

	// 启动daemon模式
	var isDaemon bool
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		_, err = daemon.Daemon(1, 0)
		if err != nil {
			log.Fatalln(err)
		}
		isDaemon = true
	}

	// 加载配置文件
	config := parseConfig(dir+CONFIG_FILE, dir)

	// 加载过滤列表
	excludeList := parseExcludeList(dir + EXCLUDE_FILE)

	// 初始化日志
	var logFileHandle *os.File
	if isDaemon {
		logFileHandle, err = os.OpenFile(config["logFile"], os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		defer logFileHandle.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
	} else {
		logFileHandle = os.Stderr
	}
	logger := log.New(logFileHandle, "", log.Ldate|log.Ltime)

	var encoder client.Encoder
	if config["versionType"] == "git" {
		// GIT
		git := new(client.GitEncoder)
		git.Config = config
		git.ExcludeList = excludeList
		git.Pwd = dir
		git.Logger = logger
		encoder = git
	} else {
		// SVN
		svn := new(client.SvnEncoder)
		svn.Config = config
		svn.ExcludeList = excludeList
		svn.Pwd = dir
		svn.Logger = logger
		encoder = svn
	}

	// 开始工作啦
	sign := make(chan bool)
	encoder.Info("autoencode starting...")
	go encoder.Run(sign)

	// 启动信号监听，等待所有Process都运行结束
	trapSignal(encoder)
	<-sign

	// 停止工作
	encoder.Info("autoencode stopped.")
}

// 解析配置文件
func parseConfig(configFile string, pwd string) map[string]string {
	conf, err := goconfig.ReadConfigFile(configFile)
	if err != nil {
		log.Fatalln("ReadConfigFile Err: ", err.Error(), "\nConfigFile:", configFile)
	}
	config := make(map[string]string)
	var field, value string
	for _, field = range configField {
		value, err = conf.GetString("default", field)
		if err != nil {
			log.Fatalln("baby, your config is error: ", err.Error())
		}
		config[field] = value
	}
	// 处理相对路径问题
	for _, field = range [3]string{"checkoutDir", "commitDir", "logFile"} {
		if config[field][0] != '/' {
			config[field] = pwd + "/" + config[field]
		}
	}
	return config
}

// 加载处理文件的排除列表
func parseExcludeList(configFile string) []string {
	fHandle, err := os.Open(configFile)
	if err != nil {
		log.Fatalln("ReadConfigFile Err: ", err.Error(), "\nConfigFile:", configFile)
	}
	excludeList := make([]string, 0)
	bufHandle := bufio.NewReader(fHandle)
	for {
		l, err := bufHandle.ReadString('\n')
		l = strings.TrimSpace(l)
		if len(l) == 0 || l[0] == ';' || l[0] == '#' {
			//continue
		} else {
			excludeList = append(excludeList, l)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln("ReadConfigFile Err: ", err.Error(), "\nConfigFile:", configFile)
		}
	}
	return excludeList
}

// 处理系统信号
// 监听系统信号，重启或停止服务
func trapSignal(server client.Encoder) {
	sch := make(chan os.Signal, 10)
	signal.Notify(sch, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT,
		syscall.SIGHUP, syscall.SIGSTOP, syscall.SIGQUIT)
	go func(ch <-chan os.Signal) {
		sig := <-ch
		server.Shutdown("signal recieved " + sig.String() + ", at: " + time.Now().String())
		if sig == syscall.SIGHUP {
			server.Info("autoencode restart now...")
			procAttr := new(os.ProcAttr)
			procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
			procAttr.Dir = os.Getenv("PWD")
			procAttr.Env = os.Environ()
			process, err := os.StartProcess(os.Args[0], os.Args, procAttr)
			if err != nil {
				server.Info("autoencode restart process failed:" + err.Error())
				return
			}
			waitMsg, err := process.Wait()
			if err != nil {
				server.Info("autoencode restart wait error:" + err.Error())
			}
			server.Info(waitMsg)
		} else {
			server.Info("autoencode shutdown now...")
		}
	}(sch)
}
