package tool

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:31
 * @params: min, max int64
 * @return: int64
 * @description: 生成随机数在min和max之间
 */
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min //nolint:gosec
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:35
 * @params: path string
 * @return: string, error
 * @description: 随机获取目录中的文件
 */
func RandFile(path string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	var fileNames []string //nolint:prealloc
	for _, v := range files {
		if strings.HasPrefix(v.Name(), ".") {
			continue
		}
		fileNames = append(fileNames, v.Name())
	}
	if len(fileNames) == 0 {
		return "", errors.New("dir is nil")
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(fileNames)) //nolint:gosec
	if index >= len(fileNames) {
		index = len(fileNames) - 1
	}
	return fileNames[index], err
}

/**
 * @author: yasinWu
 * @date: 2022/1/14 11:28
 * @params: value int64
 * @return: string
 * @description: 字节数单位换算
 */
func ByteWithUnit(value int64) string {
	if value < 1024 {
		return strconv.FormatInt(value, 10) + "B"
	}

	unit := [6]string{"KB", "MB", "GB", "TB", "PB", "EB"}

	data := float64(value)
	for i := 0; i < len(unit); i++ {
		data /= float64(1024)
		if data < 1024 {
			return strconv.FormatFloat(data, 'f', 2, 64) + unit[i]
		}
	}

	return strconv.FormatFloat(data, 'f', 2, 64) + unit[5]
}

func PrintlnFmt(data any) {
	buffer, _ := json.MarshalIndent(data, "", "   ")
	fmt.Println(string(buffer))
}
