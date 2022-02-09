package common

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	js "github.com/bitly/go-simplejson"
)

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:31
 * @params: src string
 * @return: string
 * @description: 删除字符串中的HTML标签
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
 * @author: yasinWu
 * @date: 2022/1/13 14:31
 * @params: input string
 * @return: int64
 * @description: 字符串转10进制
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
	var output []int
	for _, v := range input {
		output = append(output, int(v))
	}
	for i, j := 0, len(output)-1; i < j; i, j = i+1, j-1 {
		output[i], output[j] = output[j], output[i]
	}
	return output
}

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
	return rand.Int63n(max-min) + min
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:32
 * @params: startTime, endTime time.Time, timeFormatTpl string
 * @return: []*js.Json, error
 * @description: 切割时间区间为天
 */
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

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:33
 * @params: arr *[]string
 * @description: 删除[]string中的重复元素
 */
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

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:33
 * @params: arr *[]string
 * @description: 删除[]string中的空元素
 */
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

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:34
 * @params: startTime, endTime time.Time, timeFormatTpl string
 * @return: []string
 * @description: 获取时间区间中的天
 */
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
 * @author: yasinWu
 * @date: 2022/1/13 14:34
 * @params: file string
 * @return: string, error
 * @description: 图片转base64
 */
func ImageFileToBase64(file string) (string, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()
	fileType := strings.Replace(path.Ext(path.Base(imgFile.Name())), ".", "", -1)
	var staticImg image.Image
	switch fileType {
	case "png":
		staticImg, err = png.Decode(imgFile)
	default:
		staticImg, err = jpeg.Decode(imgFile)
	}
	if err != nil {
		return "", err
	}
	emptyBuff := bytes.NewBuffer(nil)
	switch fileType {
	case "png":
		err = png.Encode(emptyBuff, staticImg)
	default:
		err = jpeg.Encode(emptyBuff, staticImg, nil)
	}
	dist := make([]byte, 20*MB)
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes())
	index := bytes.IndexByte(dist, 0)
	baseImage := dist[0:index]
	return *(*string)(unsafe.Pointer(&baseImage)), nil
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
	var fileNames []string
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
	index := rand.Intn(len(fileNames))
	if index >= len(fileNames) {
		index = len(fileNames) - 1
	}
	return fileNames[index], err
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:36
 * @params: img image.Image, fileType string
 * @return: string, error
 * @description: image文件转base64
 */
func ImageToBase64(img image.Image, fileType string) (string, error) {
	var err error
	emptyBuff := bytes.NewBuffer(nil)
	switch fileType {
	case "png":
		err = png.Encode(emptyBuff, img)
	default:
		err = jpeg.Encode(emptyBuff, img, nil)
	}
	if err != nil {
		return "", err
	}
	dist := make([]byte, 20*1024*1024)
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes())
	index := bytes.IndexByte(dist, 0)
	baseImage := dist[0:index]
	return *(*string)(unsafe.Pointer(&baseImage)), nil
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:36
 * @params: data interface{}, result *map[string]interface{}
 * @return: error
 * @description: interface{}转map,只支持struct{}和*simplejson.Json
 */
func ToMap(data interface{}, result *map[string]interface{}) error {
	t := reflect.TypeOf(data)
	switch {
	case t.Kind() == reflect.Struct:
		v := reflect.ValueOf(data)
		for i := 0; i < t.NumField(); i++ {
			if v.Field(i).Type().Kind() == reflect.Struct {
				err := ToMap(v.Field(i).Interface(), result)
				if err != nil {
					log.Printf(err.Error())
				}
				continue
			}
			(*result)[t.Field(i).Tag.Get("json")] = v.Field(i).Interface()
		}
	case t.String() == "*simplejson.Json":
		var err error
		*result, err = data.(*js.Json).Map()
		if err != nil {
			return err
		}
	default:
		return errors.New("data type not supported")
	}
	return nil
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
		data = data / float64(1024)
		if data < 1024 {
			return strconv.FormatFloat(data, 'f', 2, 64) + unit[i]
		}
	}

	return strconv.FormatFloat(data, 'f', 2, 64) + unit[5]
}
