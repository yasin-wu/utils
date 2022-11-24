package tool

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const _png = "png"

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:31
 * @params: src string
 * @return: string
 * @description: 删除字符串中的HTML标签
 */
func RemoveHTML(src string) string {
	re := regexp.MustCompile(`\\<[\\S\\s]+?\\>`)
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	re = regexp.MustCompile(`\\<style[\\S\\s]+?\\</style\\>`)
	src = re.ReplaceAllString(src, "")

	re = regexp.MustCompile(`\\<script[\\S\\s]+?\\</script\\>`)
	src = re.ReplaceAllString(src, "")

	re = regexp.MustCompile(`\\<[\\S\\s]+?\\>`)
	src = re.ReplaceAllString(src, "\n")

	re = regexp.MustCompile(`\\s{2,}`)
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
	var output []int //nolint:prealloc
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
	return rand.Int63n(max-min) + min //nolint:gosec
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
	fileType := strings.ReplaceAll(path.Ext(path.Base(imgFile.Name())), ".", "")
	var staticImg image.Image
	switch fileType {
	case _png:
		staticImg, err = png.Decode(imgFile)
	default:
		staticImg, err = jpeg.Decode(imgFile)
	}
	if err != nil {
		return "", err
	}
	emptyBuff := bytes.NewBuffer(nil)
	switch fileType {
	case _png:
		_ = png.Encode(emptyBuff, staticImg)
	default:
		_ = jpeg.Encode(emptyBuff, staticImg, nil)
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
 * @date: 2022/1/13 14:36
 * @params: img image.Image, fileType string
 * @return: string, error
 * @description: image文件转base64
 */
func ImageToBase64(img image.Image, fileType string) (string, error) {
	var err error
	emptyBuff := bytes.NewBuffer(nil)
	switch fileType {
	case _png:
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

func Println(data any) {
	buffer, _ := json.MarshalIndent(data, "", "   ")
	fmt.Println(string(buffer))
}

func StringIn(target string, src []string) bool {
	sort.Strings(src)
	index := sort.SearchStrings(src, target)
	return index < len(src) && src[index] == target
}
