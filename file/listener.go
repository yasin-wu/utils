package file

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func Listener(file string) {
	readChannel := make(chan string)
	go readFile(readChannel, file)
	for c := range readChannel {
		fmt.Println(c)
	}
}

func readFile(readChannel chan string, file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err = f.Seek(0, os.SEEK_END); err != nil {
		log.Println(err)
		return
	}
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(time.Second)
			continue
		} else if err != nil {
			log.Println(err)
		}
		lineStr := strings.TrimSpace(string(line))
		readChannel <- lineStr
	}
}
