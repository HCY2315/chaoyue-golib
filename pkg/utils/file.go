package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"strings"
)

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	//PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, fi.Name())
		}
	}

	return files, nil
}

func MD5SumOfFile(filePath string) (string, error) {
	content, errRead := ioutil.ReadFile(filePath)
	if errRead != nil {
		return "", errRead
	}
	return fmt.Sprintf("%x", md5.Sum(content)), nil
}
