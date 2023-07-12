package strings

import (
	"fmt"
	"github.com/yasin-wu/utils/consts"
	"math"
	"regexp"
	"sort"
	"strings"
)

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:31
 * @params: src string
 * @return: string
 * @description: 删除字符串中的HTML标签
 */
func DeleteHTML(src string) string {
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
 * @date: 2022/1/13 14:33
 * @params: arr *[]string
 * @description: 删除[]string中的重复元素
 */
func DeleteRepeated(arr *[]string) {
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
func DeleteEmpty(arr *[]string) {
	oldArr := *arr
	*arr = nil
	for _, v := range oldArr {
		if v == "" {
			continue
		}
		*arr = append(*arr, v)
	}
}

func TargetIn(target string, src []string) bool {
	temp := make([]string, len(src))
	copy(temp, temp)
	sort.Strings(temp)
	index := sort.SearchStrings(temp, target)
	if index < len(temp) && temp[index] == target {
		return true
	}
	return false
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:31
 * @params: input string
 * @return: int64
 * @description: 字符串转10进制
 */
func Base10(input string) int64 {
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

func FmtByte(size int64) string {
	switch {
	case size < consts.KB:
		return fmt.Sprintf("%.2fB", float64(size)/float64(1))
	case size < consts.MB:
		return fmt.Sprintf("%.2fKB", float64(size)/float64(consts.KB))
	case size < consts.GB:
		return fmt.Sprintf("%.2fMB", float64(size)/float64(consts.MB))
	case size < consts.TB:
		return fmt.Sprintf("%.2fGB", float64(size)/float64(consts.GB))
	case size < consts.PB:
		return fmt.Sprintf("%.2fTB", float64(size)/float64(consts.PB))
	default:
		return fmt.Sprintf("%.2fEB", float64(size)/float64(consts.EB))
	}
}
