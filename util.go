package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
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

/**
 * @author: yasin
 * @date: 2020/7/22 09:32
 * @description：
 */
func RemoveHtml(src string) string {
	re, _ := regexp.Compile(`\\<[\\S\\s]+?\\>`)
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	re, _ = regexp.Compile(`\\<style[\\S\\s]+?\\</style\\>`)
	src = re.ReplaceAllString(src, "")

	re, _ = regexp.Compile(`\\<script[\\S\\s]+?\\</script\\>`)
	src = re.ReplaceAllString(src, "")

	re, _ = regexp.Compile(`\\<[\\S\\s]+?\\>`)
	src = re.ReplaceAllString(src, "\n")

	re, _ = regexp.Compile(`\\s{2,}`)
	src = re.ReplaceAllString(src, "\n")

	return src
}

/**
 * @author: yasin
 * @date: 2020/7/22 13:57
 * @description：2 to 10
 */
func ConvertString2To10(input string) int64 {
	c := getInput(input)
	out := sq(c)
	sum := 0
	for o := range out {
		sum += o
	}
	return int64(sum)
}

func getInput(input string) <-chan int {
	out := make(chan int)
	go func() {
		for _, b := range stringToIntArray(input) {
			out <- b
		}
		close(out)
	}()

	return out
}
func sq(in <-chan int) <-chan int {
	out := make(chan int)

	var base, i float64 = 2, 0
	go func() {
		for n := range in {
			out <- (n - 48) * int(math.Pow(base, i))
			i++
		}
		close(out)
	}()
	return out
}
func stringToIntArray(input string) []int {
	output := []int{}
	for _, v := range input {
		output = append(output, int(v))
	}
	for i, j := 0, len(output)-1; i < j; i, j = i+1, j-1 {
		output[i], output[j] = output[j], output[i]
	}
	return output
}
