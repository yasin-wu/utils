package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
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
	files, err := os.ReadDir(path)
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
	rand.NewSource(time.Now().UnixNano())
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

func NtToUnix(nt int64) time.Time {
	nt = (nt - 1.1644473600125e+17) / 1e+7
	return time.Unix(nt, 0)
}
