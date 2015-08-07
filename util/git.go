package util

import (
	"os"
	"path"
)

type Git struct {
	Bin    string
	Path   string
	Branch string
}

type GitLog struct {
	Commit  string     // 提交编号
	Author  string     // 作者
	Email   string     // 邮箱
	Message string     // 提交说明
	Date    string     // 提交时间
	Files   []*GitFile // 文件列表
}

type GitFile struct {
	File   string // 文件路径
	Status string // 文件状态，新增 A 修改 M 删除 D
}

type GitTag struct {
	Tag     string // 标签名称
	Message string // 标签说明
	Commit  string // 标签对应的提交编号
}

func NewGit(bin string, remote string, path string) (*Git, error) {
	git := &Git{
		Bin:  bin,
		Path: path,
	}
	return git
}

// 检出一个代码库
func (s *Git) Clone() error {

}

// 切换到一个分支
func (s *Git) Checkout(branch string, isNew bool) error {

}

// 更新一个分支的源数据
func (s *Git) Fetch() error {

}

// 合并到指定的版本
func (s *Git) Merge(commit string) error {

}

// 添加变化的文件待提交
func (s *Git) Add(path string) error {

}

// 提交变化 并返回提交后的版本号
func (s *Git) Commit() error {

}

// 推送更新
func (s *Git) Push() error {

}

/*
 * 查看日志 如果指定commit则查看指定的版本日志
 *
 * 包括：username, email, message, file list, file status
 */
func (s *Git) Log(commit string) ([]*GitLog, error) {

}

// 获取标签列表
// 返回标签和commit对应的一个列表
func (s *Git) Tags() ([]*GitTag, error) {

}

// 创建标签
func (s *Git) Tag(name string, commit string, message string) error {

}
