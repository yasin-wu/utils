package test

import (
	"fmt"
	"testing"

	"github.com/yasin-wu/utils/similarity"
)

func TestExtractWithWeight(t *testing.T) {
	input := "最近有消息指出，MagSafe充电接口将回归 MacBook，彭博社的马克·古尔曼在另一份报告中继续表示，‌MagSafe ‌将是一个独立的充电端口，而不是 USB-C 端口，新端口位于 USB-C 端口的旁边。"
	topKey := 10
	ww, s := similarity.ExtractWithWeight(input, topKey, nil)
	fmt.Println(ww)
	fmt.Println(s)
}
