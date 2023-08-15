package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min //nolint:gosec
}

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

func PrintlnFmt(data any) {
	buffer, _ := json.MarshalIndent(data, "", "   ")
	fmt.Println(string(buffer))
}
