package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	js "github.com/bitly/go-simplejson"
	"golang.org/x/net/context/ctxhttp"
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

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func CalendarDays(startTime, endTime time.Time) ([]*js.Json, error) {
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return nil, errors.New("startTime <= endTime")
	}
	u1 := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.Now().Location())
	u2 := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, time.Now().Location())
	durationDays := int(u2.Sub(u1).Hours() / 24)
	if int(endTime.Sub(startTime).Hours())%24 > 0 {
		durationDays += 1
	}
	var data []*js.Json
	for i := 0; i < durationDays; i++ {
		var stt time.Time
		var ett time.Time
		stt = startTime.AddDate(0, 0, i)
		ett = endTime
		jsonObj := js.New()
		if i == 0 && i != durationDays-1 {
			ett = time.Date(stt.Year(), stt.Month(), stt.Day(), 23, 59, 59, 0, stt.Location())
		} else if i == durationDays-1 {
			stt = time.Date(stt.Year(), stt.Month(), stt.Day(), 0, 0, 0, 0, stt.Location())
		} else {
			stt = time.Date(stt.Year(), stt.Month(), stt.Day(), 0, 0, 0, 0, stt.Location())
			ett = time.Date(stt.Year(), stt.Month(), stt.Day(), 23, 59, 59, 0, stt.Location())
		}
		jsonObj.Set("calendar", stt)
		jsonObj.Set("from_time", stt)
		jsonObj.Set("to_time", ett)
		usageTime, _ := strconv.ParseFloat(strconv.FormatFloat(ett.Sub(stt).Hours(), 'f', 1, 64), 64)
		jsonObj.Set("usage_time", usageTime)
		data = append(data, jsonObj)
	}
	return data, nil
}

func RemoveRepeatedElement(arr []string) []string {
	newArr := make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}
