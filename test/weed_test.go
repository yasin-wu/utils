package test

import (
	"fmt"
	"testing"

	"github.com/yasin-wu/utils/weed"
)

func TestWeed(t *testing.T) {
	cli := weed.New("http://192.168.131.135:9333")
	err := cli.Delete("http://192.168.131.135:9080/2,0443e6463ca727")
	if err != err {
		fmt.Println(err)
		return
	}
}
