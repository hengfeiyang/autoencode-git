package client

import (
	"autoencode-svn/util"
	"bufio"
	"errors"
	"github.com/9466/daemon"
	"github.com/9466/goconfig"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type ACST struct {
	Config      map[string]string // 配置文件
	CheckoutRev int               // 本地版本号,检出
	CommitRev   int               // 本地版本号,提交
	Pwd         string            // 当前工作目录
	svn         *util.SVNT        // SVN控制接口
	Logger      *log.Logger       // 日志处理接口
	excludeList []string          // 文件处理，排除列表
}

var configField = [18]string{"svnBin", "svnReposCheckout", "svnReposCommit",
	"svnReposBaseDir", "checkoutDir", "commitDir",
	"mailHost", "mailUser", "mailPass", "mailFrom", "mailTo", "mailTitle",
	"encodeExt", "logFile", "encodeBinPHP", "encodeBinJAVA", "encodeLibCompiler", "encodeBinUglifyJS",
}

const (
	MARK_FILE string = ".touch" // 用于标记跳过空白提交
)

func main() {
	var err error

	// 先设定工作目录
	var dir string
	dir, err = util.GetDir()
	if err != nil {
		log.Fatalln("My God, GetDir Fatal!")
	}
	acs := new(ACST)
	acs.Pwd = dir

	// 检测是否存在错误锁定
	if ok, _ := util.IsExist(acs.Pwd + util.ACS_LOCKFILE); ok {
		log.Fatalln("Sorry, acs had locked, please check error!")
	}

	// 启动daemon模式
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		_, err = daemon.Daemon(1, 0)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// 加载配置文件
	var conFile string
	conFile = "/conf/acs.conf"
	acs.parseConfig(acs.Pwd + conFile)

	// 加载过滤列表
	conFile = "/conf/exclude.conf"
	acs.parseExcludeList(acs.Pwd + conFile)

	// 初始化日志
	logFileHandle, err := os.OpenFile(acs.Config["logFile"], os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer logFileHandle.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	acs.Logger = log.New(logFileHandle, "", log.Ldate|log.Ltime)

	// 判断工作目录是否存在，如果不存在则创建
	if ok, err := util.IsExist(acs.Config["checkoutDir"]); !ok {
		err = os.Mkdir(acs.Config["checkoutDir"], 0777)
		if err != nil {
			log.Fatalln("checkoutDir create error: " + err.Error())
		}
	}
	if ok, err := util.IsExist(acs.Config["commitDir"]); !ok {
		err = os.Mkdir(acs.Config["commitDir"], 0777)
		if err != nil {
			log.Fatalln("commitDir create error: " + err.Error())
		}
	}

	// 初始化SVN结构
	acs.svn = new(util.SVNT)
	acs.svn.Bin = acs.Config["svnBin"]

	// 初始化SVN目录
	if ok, err := util.IsExist(acs.Config["checkoutDir"] + ".svn"); !ok {
		err = acs.svn.Checkout(acs.Config["checkoutDir"], acs.Config["svnReposCheckout"], 0)
		if err != nil {
			log.Fatalln("checkoutDir init svn error: " + err.Error())
		}
	}
	if ok, err := util.IsExist(acs.Config["commitDir"] + ".svn"); !ok {
		err = acs.svn.Checkout(acs.Config["commitDir"], acs.Config["svnReposCommit"], 0)
		if err != nil {
			log.Fatalln("commitDir init svn error: " + err.Error())
		}
	}

	// 检查错误是否已经被修正
	if ok, err := acs.svn.Status(acs.Config["commitDir"]); ok {
		log.Fatalln("oh! there's some error need you check, to be use [svn status] error: " + err.Error())
	}

	// 获取当前的工作的SVN版本
	var svnInfo = new(util.SvnInfo)
	svnInfo, err = acs.svn.Info(acs.Config["checkoutDir"])
	if err != nil {
		log.Fatalln("checkoutDir init svn info error: " + err.Error())
	}
	acs.CheckoutRev = svnInfo.Rev
	svnInfo, err = acs.svn.Info(acs.Config["commitDir"])
	if err != nil {
		log.Fatalln("commitDir init svn info error: " + err.Error())
	}
	acs.CommitRev = svnInfo.Rev

	// 开始工作啦
	for {
		time.Sleep(time.Second)
		acs.run()
	}
}

// 解析配置文件
func (c *ACST) parseConfig(configFile string) {
	conf, err := goconfig.ReadConfigFile(configFile)
	if err != nil {
		log.Fatalln("ReadConfigFile Err: ", err.Error(), "\nConfigFile:", configFile)
	}
	c.Config = make(map[string]string)
	var field, value string
	for _, field = range configField {
		value, err = conf.GetString("default", field)
		if err != nil {
			log.Fatalln("baby, your config is error: ", err.Error())
		}
		c.Config[field] = value
	}
	// 处理相对路径问题
	for _, field = range [3]string{"checkoutDir", "commitDir", "logFile"} {
		if c.Config[field][0] != '/' {
			c.Config[field] = c.Pwd + "/" + c.Config[field]
		}
	}
}

// 加载处理文件的排除列表
func (c *ACST) parseExcludeList(configFile string) error {
	fHandle, err := os.Open(configFile)
	if err != nil {
		return err
	}
	c.excludeList = make([]string, 0)
	bufHandle := bufio.NewReader(fHandle)
	for {
		l, err := bufHandle.ReadString('\n')
		l = strings.TrimSpace(l)
		if len(l) == 0 || l[0] == ';' || l[0] == '#' {
			//continue
		} else {
			c.excludeList = append(c.excludeList, l)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

// 执行工作
func (c *ACST) run() {
	// 获取远程新版本
	var svnInfo = new(util.SvnInfo)
	var err error
	svnInfo, err = c.svn.Info(c.Config["svnReposCheckout"])
	if err != nil {
		c.Logger.Fatalln(c.Config["svnReposCheckout"] + " svn info error: " + err.Error())
		return
	}

	// 判断当于工作版本与远程版本
	if c.CheckoutRev >= svnInfo.Rev {
		return
	}

	// 版本较小，开始更新
	var opRev int = c.CheckoutRev + 1
	err = c.svn.Update(c.Config["checkoutDir"], opRev)
	if err != nil {
		c.Logger.Fatalln(c.Config["checkoutDir"] + " svn update error: " + err.Error())
		return
	}
	c.Logger.Println("checkoutDir svn updated  Rev " + strconv.Itoa(opRev))

	// 定义失败时的回滚
	var checkoutRollBack = true
	defer func() {
		// 无错误
		if err == nil {
			return
		}
		// 回滚操作
		if checkoutRollBack {
			c.Logger.Println("checkoutDir has RollBack Rev " + strconv.Itoa(opRev-1) + " !!!")
			c.svn.Update(c.Config["checkoutDir"], opRev-1)
		}
		// 终止执行
		c.lockQuit(err.Error())
	}()

	// 获取版本更新详情
	var opLog *util.SvnLog
	opLog = new(util.SvnLog)
	opLog, err = c.svn.Log(c.Config["checkoutDir"], opRev)
	if err != nil {
		err = errors.New("checkoutDir svn log error: " + err.Error())
		return
	}

	// 遍历文件进行处理
	// 标记SVN操作
	// 记录一个数组，将A和D的操作记录下来，先A再D
	// 加入标示数组时，应先判断父路径是否存在，如果父路径存在子目录就不要再添加了，svn有目录继承
	// 因为添加时存在一个目录级问题，如果父目录是copy而来，后面的add文件无法生成，故，add操作直接执行，不再做标记
	svnCommit := make(map[string][]string)
	svnCommit["D"] = make([]string, 0)
	for _, item := range opLog.Paths {
		// 处理相对文件路径
		item["path"] = item["path"][len(c.Config["svnReposBaseDir"]):]
		switch item["action"] {
		case "A":
			// 创建文件或目录，立即执行创建操作
			if item["copyfrom"] != "" {
				item["copyfrom"] = item["copyfrom"][len(c.Config["svnReposBaseDir"]):]
				// 构造源文件 http://repos.com/db/file@rev
				var sfile string
				if c.Config["svnReposCommit"][len(c.Config["svnReposCommit"])-1] == '/' {
					sfile = c.Config["svnReposCommit"] + item["copyfrom"]
				} else {
					sfile = c.Config["svnReposCommit"] + "/" + item["copyfrom"]
				}
				sfile += "@" + item["copyfrom-rev"]
				err = c.svn.Copy(c.Config["commitDir"], sfile, item["path"])
				if err != nil {
					err = errors.New(item["path"] + " (copyfrom " + item["copyfrom"] + " ) " + err.Error())
					return
				}
				// 判断，如果是文件还要复制一次，因为有可能在copy之后又修改了内容
				if item["kind"] != "dir" {
					// 创建文件
					err = c.encodeFile(item["path"])
					if err != nil {
						return
					}
				}
			} else {
				if item["kind"] == "dir" {
					// 创建目录
					err = os.Mkdir(c.Config["commitDir"]+item["path"], 0755)
					if err != nil {
						err = errors.New(item["path"] + " mkdir error: " + err.Error())
						return
					}
				} else {
					// 创建文件
					err = c.encodeFile(item["path"])
					if err != nil {
						return
					}
				}
				err = c.svn.Add(c.Config["commitDir"], item["path"])
				if err != nil {
					err = errors.New(item["path"] + " " + err.Error())
					return
				}
			}
		case "D":
			// 删除文件或目录，加入SVN标记
			if ok := util.PathPrefix(svnCommit["D"], item["path"]); !ok {
				svnCommit["D"] = append(svnCommit["D"], item["path"])
			}
		case "M":
			// 不处理目录，但文件改变了还是要处理一下文件的
			if item["kind"] != "dir" {
				err = c.encodeFile(item["path"])
				if err != nil {
					return
				}
			}
		case "R":
			// 先delete掉原文件，再创建新文件，再加入SVN标记
			err := c.svn.Delete(c.Config["commitDir"], item["path"])
			if err != nil {
				err = errors.New(item["path"] + " " + err.Error())
			}
			if item["copyfrom"] != "" {
				item["copyfrom"] = item["copyfrom"][len(c.Config["svnReposBaseDir"]):]
				// 构造源文件 http://repos.com/db/file@rev
				var sfile string
				if c.Config["svnReposCommit"][len(c.Config["svnReposCommit"])-1] == '/' {
					sfile = c.Config["svnReposCommit"] + item["copyfrom"]
				} else {
					sfile = c.Config["svnReposCommit"] + "/" + item["copyfrom"]
				}
				sfile += "@" + item["copyfrom-rev"]
				err = c.svn.Copy(c.Config["commitDir"], sfile, item["path"])
				if err != nil {
					err = errors.New(item["path"] + " (copyfrom " + item["copyfrom"] + " ) " + err.Error())
					return
				}
				// 判断，如果是文件还要复制一次，因为有可能在copy之后又修改了内容
				if item["kind"] != "dir" {
					// 创建文件
					err = c.encodeFile(item["path"])
					if err != nil {
						return
					}
				}
			} else {
				if item["kind"] == "dir" {
					// 创建目录
					err = os.Mkdir(c.Config["commitDir"]+item["path"], 0755)
					if err != nil {
						err = errors.New(item["path"] + " mkdir error: " + err.Error())
						return
					}
				} else {
					// 创建文件
					err = c.encodeFile(item["path"])
					if err != nil {
						return
					}
				}
				err = c.svn.Add(c.Config["commitDir"], item["path"])
				if err != nil {
					err = errors.New(item["path"] + " " + err.Error())
					return
				}
			}
		default:
			// 其它处理，暂无
		}
	}
	// 执行SVN操作
	// path是经过排序的，故这里不再排序
	var vpath string
	for _, vpath = range svnCommit["D"] {
		err = c.svn.Delete(c.Config["commitDir"], vpath)
		if err != nil {
			err = errors.New(vpath + " " + err.Error())
			return
		}
	}

	// 提交前，先st检测一下，如果没有变化，退出执行，应该发生问题了
	if ok, _ := c.svn.Status(c.Config["commitDir"]); !ok {
		// 通常这里只是因为某个提交只修改了注释或空格之类，在代码压缩后结果没有变化，导致st为空
		// 此时，标记一下，继续进行
		err = c.markBlank(opRev)
		if err != nil {
			err = errors.New("oh! i will commit, but there is nothing! i try to markBlank but error again: " + err.Error())
			return
		}
	}

	// 提交文件，若失败则通知并退出
	err = c.svn.Commit(c.Config["commitDir"], opLog.Author, opLog.Msg)
	if err != nil {
		checkoutRollBack = false // 不要回滚，已经走到了提交这一部，如果出错，管理员手工提交一下，不再回滚
		err = errors.New("commitDir commit error: " + err.Error())
		return
	}

	// 检查提交是否OK，不知为什么有文件不能被提交，特意检查一下
	if ok, err := c.svn.Status(c.Config["commitDir"]); ok {
		err = errors.New("oh! there's some error need you check, to be use [svn status] " + err.Error())
		return
	}

	// 提交成功，update操作目录 commitDir
	c.CommitRev++
	c.CheckoutRev++
	c.svn.Update(c.Config["commitDir"], c.CommitRev)
	c.Logger.Println("commitDir   svn commited Rev " + strconv.Itoa(c.CommitRev))
	err = nil
	return
}

// 文件编码处理
// 判断文件类型，如果是php,js要处理，如果是目录，跳过处理，直接创建
func (c *ACST) encodeFile(spath string) error {
	// for debug
	c.Logger.Println("beggin encode file: " + spath)

	// 排除处理，默认需要处理
	need := true
	for _, item := range c.excludeList {
		if strings.Contains(spath, item) {
			// 发现排除，设为false
			need = false
			break
		}
	}

	var err error
	fileExt := util.GetExt(spath)
	fileDir := path.Dir(spath)
	if ok, _ := util.IsExist(c.Config["commitDir"] + fileDir); !ok {
		err = os.MkdirAll(c.Config["commitDir"]+fileDir, 0755)
		if err != nil {
			return err
		}
	}
	if need && fileExt != "" && strings.Contains(c.Config["encodeExt"], fileExt) {
		switch fileExt {
		case "php":
			err = util.EncodePhp(c.Config["encodeBinPHP"], c.Config["checkoutDir"]+spath, c.Config["commitDir"]+spath)
			if err != nil {
				return errors.New(c.Config["checkoutDir"] + spath + " encodePHP error: " + err.Error())
			}
		case "js":
			//err = util.EncodeJs(c.Config["encodeBinJAVA"], c.Config["encodeLibCompiler"], c.Config["checkoutDir"]+spath, c.Config["commitDir"]+spath)
			err = util.EncodeJsUglifyJS(c.Config["encodeBinUglifyJS"], c.Config["checkoutDir"]+spath, c.Config["commitDir"]+spath)
			if err != nil {
				if err.Error() == "exit status 127" {
					// 这个错误未知，重试一下看看
					c.Logger.Println(c.Config["checkoutDir"] + spath + " encodeJS error: " + err.Error())
					c.Logger.Println("encodeJS retry")
					err = util.EncodeJsUglifyJS(c.Config["encodeBinUglifyJS"], c.Config["checkoutDir"]+spath, c.Config["commitDir"]+spath)
					if err != nil {
						return errors.New(c.Config["checkoutDir"] + spath + " encodeJS error: " + err.Error())
					}
					// end retry
				}
				return errors.New(c.Config["checkoutDir"] + spath + " encodeJS error: " + err.Error())
			}
		case "css":
			err = util.EncodeCss(c.Config["checkoutDir"]+spath, c.Config["commitDir"]+spath)
			if err != nil {
				return errors.New(c.Config["checkoutDir"] + spath + " encodeCss error: " + err.Error())
			}
		default:
			// 其它类型暂不支持，直接将文件COPY过去
			err = util.CopyFile(c.Config["checkoutDir"]+spath, c.Config["commitDir"]+spath)
			if err != nil {
				return errors.New(c.Config["checkoutDir"] + spath + " copy file error: " + err.Error())
			}
		}
	} else {
		// 不进行处理的文件，直接COPY过去
		util.CopyFile(c.Config["checkoutDir"]+spath, c.Config["commitDir"]+spath)
		if err != nil {
			return errors.New(c.Config["checkoutDir"] + spath + " copy file error: " + err.Error())
		}
	}

	return nil
}

// 标记跳过的处理
func (c *ACST) markBlank(rev int) (err error) {
	file := c.Config["commitDir"] + MARK_FILE
	if ok, _ := util.IsExist(file); !ok {
		// 创建
		_, err = os.Create(file)
		if err != nil {
			return err
		}
		// 添加到svn
		err = c.svn.Add(c.Config["commitDir"], MARK_FILE)
		if err != nil {
			err = errors.New(MARK_FILE + " " + err.Error())
			return
		}
	}
	// 更新
	err = ioutil.WriteFile(file, []byte(strconv.Itoa(rev)+"\n"), 0644)
	return err
}

// 发生错误，锁定程序并退出执行，哦，还要发邮件通知管理员
func (c *ACST) lockQuit(msg string) {
	if ok, _ := util.IsExist(c.Pwd + util.ACS_LOCKFILE); !ok {
		os.Create(c.Pwd + util.ACS_LOCKFILE)
	}
	// 发送邮件
	mailConfig := &util.MailT{
		Addr:  c.Config["mailHost"],
		User:  c.Config["mailUser"],
		Pass:  c.Config["mailPass"],
		From:  c.Config["mailFrom"],
		To:    c.Config["mailTo"],
		Title: c.Config["mailTitle"],
		Type:  "",
	}
	if mailConfig.Title == "" {
		mailConfig.Title = "9466代码机器人"
	}
	mailConfig.Body = "acs is locked, please check! " + time.Now().Format("2006-01-02 15:04:05") + "\r\n\r\n" + msg
	err := util.SendMail(mailConfig)
	if err != nil {
		c.Logger.Println(mailConfig.Body)
		c.Logger.Println("acs warning send error: ", err.Error())
	}
	c.Logger.Println(msg)
	c.Logger.Fatalln("acs stopped!")
}
