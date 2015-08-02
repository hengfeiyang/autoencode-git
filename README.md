# 自动编码 for Git

该程序从指定的git库中读取数据，同步制定的分支到新的git库中。在同步的过程中，可以执行制定的处理。目前，已内置的处理有：

1. ZendGuard加密PHP
2. UglifyJS加密JS
3. 内置压缩CSS

并与源库的提交保持一致的提交记录（提交人，提交时间，提交信息和文件编号）。同时支持创建同步的TAG标签。

## 安装

直接从git库中检出执行go install即可：

```
git clone http://git.cmstop.cc:10080/mediacloud/go-autoencode-git.git
cd go-autoencode-git
go build
```

编译即可得到主程序。

配置文件示例：

1. acs.conf 主配置文件，主要指定来源git和目标git，以及相关依赖的路径，mail通知设置等。
2. exclude.conf 辅助配置文件，在处理的过程中，如果某些文件不需要处理，可指定跳过。

## 依赖

该程序依赖几个系统的东西：

1. git命令，用于git操作
2. node和uglify模块，用于混淆JS
3. ZendGuard，用于加密PHP

依赖的东东，需要提前安装好，然后在配置文件acs.conf中指定。
