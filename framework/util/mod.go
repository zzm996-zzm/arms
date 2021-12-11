package util

import (
	"bufio"
	"io"
	"os"
)

//获取go.mod文件的module
func GetModule(modPath string) string {
	//打开文件
	fi, err := os.Open(modPath)
	if err != nil {
		panic(err)
		return ""
	}
	defer fi.Close()

	//读取文件,仅仅读取第一行
	br := bufio.NewReader(fi)
	line, _, c := br.ReadLine()
	if c == io.EOF {
		return ""
	}
	return string(line)
}
