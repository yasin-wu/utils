package common

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

func CalendarDays(startTime, endTime time.Time, timeFormatTpl string) ([]*js.Json, error) {
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return nil, errors.New("startTime <= endTime")
	}
	if timeFormatTpl == "" {
		timeFormatTpl = "2006-01-02"
	}
	days := GetBetweenDates(startTime, endTime, timeFormatTpl)
	var data []*js.Json
	for i, v := range days {
		vt, _ := time.ParseInLocation("2006-01-02", v, time.Now().Location())
		var fromTime time.Time
		var toTime time.Time
		switch i {
		case 0:
			fromTime = startTime
			toTime = time.Date(vt.Year(), vt.Month(), vt.Day(), 23, 59, 59, 0, time.Now().Location())
		case len(days) - 1:
			fromTime = time.Date(vt.Year(), vt.Month(), vt.Day(), 0, 0, 0, 0, time.Now().Location())
			toTime = endTime
		default:
			fromTime = time.Date(vt.Year(), vt.Month(), vt.Day(), 0, 0, 0, 0, time.Now().Location())
			toTime = time.Date(vt.Year(), vt.Month(), vt.Day(), 23, 59, 59, 0, time.Now().Location())
		}
		jsonObj := js.New()
		jsonObj.Set("calendar", vt)
		jsonObj.Set("from_time", fromTime)
		jsonObj.Set("to_time", toTime)
		usageTime, _ := strconv.ParseFloat(strconv.FormatFloat(toTime.Sub(fromTime).Hours(), 'f', 1, 64), 64)
		jsonObj.Set("usage_time", usageTime)
		data = append(data, jsonObj)
	}
	return data, nil
}

func RemoveRepeatedElement(arr *[]string) {
	oldArr := *arr
	*arr = nil
	for i := 0; i < len(oldArr); i++ {
		repeat := false
		for j := i + 1; j < len(oldArr); j++ {
			if (oldArr)[i] == (oldArr)[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			*arr = append(*arr, oldArr[i])
		}
	}
}

func DelEmptyStrings(arr *[]string) {
	oldArr := *arr
	*arr = nil
	for _, v := range oldArr {
		if v == "" {
			continue
		}
		*arr = append(*arr, v)
	}
}

func GetBetweenDates(startTime, endTime time.Time, timeFormatTpl string) []string {
	var days []string
	if endTime.Before(startTime) {
		return nil
	}
	if timeFormatTpl == "" {
		timeFormatTpl = "2006-01-02"
	}
	endTimeStr := endTime.Format(timeFormatTpl)
	days = append(days, startTime.Format(timeFormatTpl))
	st := startTime.AddDate(0, 0, 1)
	stStr := st.Format(timeFormatTpl)
	if stStr > endTimeStr {
		return days
	}
	for {
		startTime = startTime.AddDate(0, 0, 1)
		startTimeStr := startTime.Format(timeFormatTpl)
		days = append(days, startTimeStr)
		if startTimeStr == endTimeStr {
			break
		}
	}
	return days
}

/**
 * @author: yasin
 * @date: 2021/4/2 15:36
 * @description：判断两个区间是否存在交集,true is mixed
 */
type Interval struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

func (i *Interval) IntervalMixed(interval *Interval) bool {
	startMax := math.Max(float64(i.Start), float64(interval.Start))
	endMin := math.Min(float64(i.End), float64(interval.End))
	return startMax <= endMin
}
