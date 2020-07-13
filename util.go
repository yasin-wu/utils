package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

/**
 * @author: yasin
 * @date: 2020/2/25 14:02
 * @description：print json to file
 */
func PrintJson(to string, data interface{}) error {
	buf, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	err = json.Indent(&out, buf, "", "\t")
	if err != nil {
		return err
	}

	file, err := os.Create(to)
	if err != nil {
		return err
	}

	if _, err := out.WriteTo(file); err != nil {
		return err
	}
	return nil
}

/**
 * @author: yasin
 * @date: 2020/5/25 14:41
 * @description：put file to url
 */
func PutFileToUrl(url string, file *os.File) ([]byte, error) {
	req, err := http.NewRequest("PUT", url, file)
	if err != nil {
		return nil, err
	}

	resp, err := ctxhttp.Do(context.Background(), http.DefaultClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code %v", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

/**
 * @author: yasin
 * @date: 2020/5/25 15:55
 * @description：get file from url
 */
func GetFileFromUrl(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

/**
 * @author: yasin
 * @date: 2020/7/13 10:30
 * @description：analysis file basic info
 */
func ParseFileInfo(file *os.File) *FileInfo {
	fileInfo := &FileInfo{
		Name:     path.Base(file.Name()),
		Path:     file.Name(),
		FileType: path.Ext(path.Base(file.Name())),
	}
	return fileInfo
}
