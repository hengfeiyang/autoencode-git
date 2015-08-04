package util

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
	"time"
)

type SVNT struct {
	Bin string
}

type SvnInfo struct {
	Rev        int       // 最新版本号
	LastAuthor string    // 最后提交者
	LastDate   time.Time // 最后提交时间
}

type SvnLog struct {
	Rev    int                 // 本次提交版本
	Author string              // 提交者
	Msg    string              // 提交信息
	Date   time.Time           // 提交时间
	Paths  []map[string]string // 文件列表
}

// 解析info用的xml结构
type xmlInfo struct {
	XMLName xml.Name `xml:"info"`
	Entry   xmlEntry `xml:"entry"`
}

type xmlEntry struct {
	Rev    int       `xml:"revision,attr"`
	Author string    `xml:"commit>author"`
	Date   time.Time `xml:"commit>date"`
}

// 解析log用的xml结构
type xmlLog struct {
	XMLName  xml.Name    `xml:"log"`
	Logentry xmlLogentry `xml:"logentry"`
}

type xmlLogentry struct {
	Rev    int       `xml:"revision,attr"`
	Author string    `xml:"author"`
	Date   time.Time `xml:"date"`
	Msg    string    `xml:"msg"`
	Paths  []xmlPath `xml:"paths>path"`
}

type xmlPath struct {
	Kind        string `xml:"kind,attr"`
	Action      string `xml:"action,attr"`
	CopyFrom    string `xml:"copyfrom-path,attr"`
	CopyFromRev string `xml:"copyfrom-rev,attr"`
	Path        string `xml:",chardata"`
}

// svn info
func (s *SVNT) Info(repos string) (*SvnInfo, error) {
	argv := make([]string, 3)
	argv[0] = "info"
	argv[1] = "--xml"
	argv[2] = repos
	data, err := Command(s.Bin, argv, "")
	if err != nil {
		return nil, err
	}
	v := new(xmlInfo)
	err = xml.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	info := new(SvnInfo)
	info.Rev = v.Entry.Rev
	info.LastAuthor = v.Entry.Author
	info.LastDate = v.Entry.Date
	return info, nil
}

// svn status
func (s *SVNT) Status(path string) (bool, error) {
	argv := make([]string, 2)
	argv[0] = "status"
	argv[1] = path
	data, err := Command(s.Bin, argv, path)
	if err != nil {
		return true, err
	}
	if len(data) > 0 {
		return true, errors.New(string(data))
	} else {
		return false, nil
	}
}

// svn checkout
func (s *SVNT) Checkout(path, repos string, rev int) error {
	argv := make([]string, 5)
	argv[0] = "checkout"
	argv[1] = "-r"
	argv[2] = strconv.Itoa(rev)
	argv[3] = repos
	argv[4] = path
	_, err := Command(s.Bin, argv, "")
	return err
}

// svn update
func (s *SVNT) Update(path string, rev int) error {
	argv := make([]string, 4)
	argv[0] = "update"
	argv[1] = "-r"
	argv[2] = strconv.Itoa(rev)
	argv[3] = path
	_, err := Command(s.Bin, argv, path)
	return err
}

// svn log
func (s *SVNT) Log(repos string, rev int) (*SvnLog, error) {
	argv := make([]string, 6)
	argv[0] = "log"
	argv[1] = "-v"
	argv[2] = "-r"
	argv[3] = strconv.Itoa(rev)
	argv[4] = "--xml"
	argv[5] = repos
	data, err := Command(s.Bin, argv, "")
	if err != nil {
		return nil, err
	}
	v := new(xmlLog)
	err = xml.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	info := new(SvnLog)
	info.Rev = v.Logentry.Rev
	info.Author = v.Logentry.Author
	info.Date = v.Logentry.Date
	info.Msg = v.Logentry.Msg
	info.Paths = make([]map[string]string, len(v.Logentry.Paths))
	path := make(map[string]string)
	for k, v := range v.Logentry.Paths {
		path = map[string]string{
			"kind":         v.Kind,
			"action":       v.Action,
			"path":         v.Path,
			"copyfrom":     v.CopyFrom,
			"copyfrom-rev": v.CopyFromRev,
		}
		info.Paths[k] = path
	}
	// 排序
	MapSort(info.Paths, func(p1, p2 map[string]string) bool {
		return p1["path"] < p2["path"]
	})
	return info, nil
}

// svn add
// ** svn add时必须在当前目录下操作
func (s *SVNT) Add(path, file string) error {
	file = svnFileName(file)
	argv := make([]string, 2)
	argv[0] = "add"
	argv[1] = file
	data, err := Command(s.Bin, argv, path)
	if err != nil {
		return errors.New("svn add error: " + string(data) + err.Error())
	}
	return nil
}

// svn delete
// ** svn add时必须在当前目录下操作
func (s *SVNT) Delete(path, file string) error {
	file = svnFileName(file)
	argv := make([]string, 2)
	argv[0] = "delete"
	argv[1] = file
	data, err := Command(s.Bin, argv, path)
	if err != nil {
		return errors.New("svn delete error: " + string(data) + err.Error())
	}
	return nil
}

// svn rename
// ** svn add时必须在当前目录下操作
func (s *SVNT) Copy(path, sfile, dfile string) error {
	argv := make([]string, 3)
	argv[0] = "copy"
	argv[1] = sfile
	argv[2] = svnFileName(dfile)
	data, err := Command(s.Bin, argv, path)
	if err != nil {
		return errors.New("svn copy error: " + string(data) + err.Error())
	}
	return nil
}

// svn revert
// ** svn add时必须在当前目录下操作
func (s *SVNT) Revert(path, file string) error {
	file = svnFileName(file)
	argv := make([]string, 3)
	argv[0] = "revert"
	argv[1] = "-R"
	argv[2] = file
	data, err := Command(s.Bin, argv, path)
	if err != nil {
		return errors.New("svn revert error: " + string(data) + err.Error())
	}
	return nil
}

// svn commit
// ** svn add时必须在当前目录下操作
func (s *SVNT) Commit(path, author, msg string) error {
	argv := make([]string, 4)
	argv[0] = "commit"
	argv[1] = "-m"
	argv[2] = msg
	argv[3] = "--username=" + author
	data, err := Command(s.Bin, argv, path)
	if err != nil {
		return errors.New("svn commit error: " + string(data) + err.Error())
	}
	return nil
}

// svn有个问题，如果文件名中含有@,需要在结尾再加个@
func svnFileName(name string) string {
	if strings.Contains(name, "@") {
		return name + "@"
	}
	return name
}
