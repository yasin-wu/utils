package frp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	timeoutSignal        = "Read file timeout"
	newProxySignal       = regexp.MustCompile(`new proxy \[[\w\W]+] type \[tcp] success`)
	closeProxyIDSignal   = regexp.MustCompile(`\[FRP.*\] listener is closed: accept tcp`)
	closeProxyPortSignal = regexp.MustCompile(`.*?(\d+): use of closed network connection`)
)

type Listener struct {
	file    string
	timeout time.Duration
}

type CloseProxyHandle func(line []byte)

func New(file string, timeout time.Duration) *Listener {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Listener{
		file:    file,
		timeout: timeout,
	}
}

func (l *Listener) NewProxy() bool {
	readChannel := make(chan string)
	defer close(readChannel)
	go l.newProxy(readChannel)
	for c := range readChannel {
		if c == timeoutSignal {
			return false
		}
		if newProxySignal.MatchString(c) {
			return true
		}
	}
	return false
}

func (l *Listener) CloseProxy(handlers ...CloseProxyHandle) {
	file, err := os.Open(l.file)
	if err != nil {
		return
	}
	if _, err = file.Seek(0, io.SeekEnd); err != nil {
		return
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			continue
		} else if err != nil {
			continue
		}
		for _, f := range handlers {
			f(line)
		}
	}
}

func (l *Listener) GetID() CloseProxyHandle {
	return func(line []byte) {
		resp := closeProxyIDSignal.FindAllStringSubmatch(strings.TrimSpace(string(line)), 1)
		if len(resp) > 0 && len(resp[0]) > 0 {
			id := strings.ReplaceAll(resp[0][0], "] listener is closed: accept tcp", "")
			id = strings.ReplaceAll(id, "[FRP", "")
			fmt.Println(id)
		}
	}
}

func (l *Listener) GetPort() CloseProxyHandle {
	return func(line []byte) {
		resp := closeProxyPortSignal.FindAllStringSubmatch(strings.TrimSpace(string(line)), 1)
		if len(resp) > 0 && len(resp[0]) > 1 {
			port := resp[0][1]
			fmt.Println(port)
		}
	}
}

func (l *Listener) newProxy(readChannel chan string) {
	file, err := os.Open(l.file)
	if err != nil {
		return
	}
	if _, err = file.Seek(0, io.SeekEnd); err != nil {
		return
	}
	reader := bufio.NewReader(file)
	timeoutChan := time.After(l.timeout)
	for {
		select {
		case <-timeoutChan:
			readChannel <- timeoutSignal
			return
		default:
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				continue
			} else if err != nil {
				continue
			}
			readChannel <- strings.TrimSpace(string(line))
			if newProxySignal.MatchString(strings.TrimSpace(string(line))) {
				return
			}
		}
	}
}
