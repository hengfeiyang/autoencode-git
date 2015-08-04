package util

// import (
// 	"testing"
// )

// var svnbin = "/usr/local/bin/svn"
// var svnRepos = "file:///Users/yanghengfei/coldstar/FYLM/SVNDB/test1"
// var dir1 = "/Users/yanghengfei/test/var/checkout"
// var dir2 = "/Users/yanghengfei/test/var/checkout2"
// var dir3 = "/Users/yanghengfei/test/var/commit"

// func TestCheckout1(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	err := svn.Checkout(dir1, svnRepos, 0)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestCheckout2(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	err := svn.Checkout(dir2, svnRepos, 0)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestUpdate1(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	err := svn.Update(dir1, 1)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestUpdate2(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	err := svn.Update(dir1, 0)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestUpdate3(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	err := svn.Update(dir3, 1)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestAdd1(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	err := svn.Add(dir1, "t2.php")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestDelete1(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	err := svn.Delete(dir1, "t.php")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestDelete2(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	err := svn.Delete(dir1, "t3.php")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestCommit1(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	err := svn.Commit(dir1, "micate", test1")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestCommit2(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	err := svn.Commit(dir3, "micate", "test3")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestLog1(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	_, err := svn.Log(dir1, 0)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestLog2(t *testing.T) {
// 	svn := new(SVNT)
// 	svn.Bin = svnbin
// 	_, err := svn.Log(dir3, 1)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
