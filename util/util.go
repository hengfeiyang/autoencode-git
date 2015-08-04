package util

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

// 加密PHP代码
func EncodePhp(zendenc, sfile, dfile string) error {
	argv := make([]string, 4)
	argv[0] = "--no-header"
	argv[1] = "--silent"
	argv[2] = sfile
	argv[3] = dfile
	_, err := Command(zendenc, argv, path.Dir(zendenc))
	// 如果PHP编译错误，那就不管他了，可能文件本身有问题，直接COPY过去
	if err != nil {
		err = CopyFile(sfile, dfile)
	}
	return err
}

// 加密JS代码
func EncodeJs(java, compiler, sfile, dfile string) error {
	argv := make([]string, 8)
	argv[0] = "-jar"
	argv[1] = compiler
	argv[2] = "--warning_level"
	argv[3] = "QUIET"
	argv[4] = "--js"
	argv[5] = sfile
	argv[6] = "--js_output_file"
	argv[7] = dfile
	_, err := Command(java, argv, "")
	// 如果JS编译错误，那就不管他了，可能文件本身有问题，直接COPY过去
	if err != nil {
		err = CopyFile(sfile, dfile)
	}
	return err
}

// 加密JS代码 UglifyJS版本
func EncodeJsUglifyJS(bin, sfile, dfile string) error {
	argv := make([]string, 8)
	argv[0] = "-c"
	argv[1] = "warnings=false"
	argv[2] = "-m"
	argv[3] = "-r"
	argv[4] = "$,require,exports"
	argv[5] = "-o"
	argv[6] = dfile
	argv[7] = sfile
	res, err := Command(bin, argv, "")
	if err != nil {
		return errors.New(string(res) + err.Error())
	}
	return err
}

// 压缩css代码
func EncodeCss(sfile, dfile string) error {
	in, err := ioutil.ReadFile(sfile)
	if err != nil {
		return err
	}
	out, err := CssCompress(in)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dfile, out, 0644)
	return err
}

// 执行系统命令并返回结果
func Command(pro string, argv []string, baseDir string) ([]byte, error) {
	cmd := exec.Command(pro, argv...)
	// 设置命令运行时目录
	if baseDir != "" {
		cmd.Dir = baseDir
	}
	res, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取程序运行的目录
func GetDir() (string, error) {
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Dir(path), nil
}

// 判断一个文件或目录是否存在
func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	// Check if error is "no such file or directory"
	if _, ok := err.(*os.PathError); ok {
		return false, nil
	}
	return false, err
}

// 判断一个文件或目录是否有写入权限
func IsWritable(path string) (bool, error) {
	err := syscall.Access(path, syscall.O_RDWR)
	if err == nil {
		return true, nil
	}
	// Check if error is "no such file or directory"
	if _, ok := err.(*os.PathError); ok {
		return false, nil
	}
	return false, err
}

// 读取一个文件夹返回文件列表
func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 获取一个文件的文件后缀名
func GetExt(filename string) string {
	info := strings.Split(filename, ".")
	if len(info) < 2 {
		return ""
	}
	return info[len(info)-1]
}

// 复制文件，仅文件，不支持目录
func CopyFile(s, d string) error {
	// 坑爹啊，要先删除是不是link
	linfo, err := os.Readlink(s)
	if err == nil || len(linfo) > 0 {
		// 这货是link，创建link吧
		return os.Symlink(linfo, d)
	}
	// 不是link，创建文件
	sf, err := os.Open(s)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(d)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}

// 判断是否存在一个路径的父路径
func PathPrefix(m []string, path string) bool {
	for _, v := range m {
		if strings.HasPrefix(path, v+"/") {
			return true
		}
	}
	return false
}
