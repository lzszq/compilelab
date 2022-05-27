# 词法分析程序使用说明

## 环境

windows 11
go version go1.17.7 windows/amd64

## 所需依赖

若要生成 NFA，DFA等图，须先执行 **go get github.com/goccy/go-graphviz** 以安装所需依赖。

## 编译方法

位于 **lab1** 项目根目录，执行 **cd ./run/main** ，然后 **go build main.go draw.go** 。

## 使用方法

要使用序列化文件时，请保证 **main.exe** 同级目录下存在 **gob.data** 。

**./main.exe** ，并根据提示操作。