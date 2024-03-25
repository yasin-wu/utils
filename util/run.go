package util

import (
	"bufio"
	"bytes"
	"os/exec"
	"runtime"
	"strings"
)

const (
	Kubepods = "kubepods"
	Docker   = "docker"
	Process  = "process"
	cgroup   = "/proc/1/cgroup"
	devices  = "devices"
	name     = "name"
)

func Runat() (string, error) {
	if strings.ToLower(runtime.GOOS) == "windows" || strings.ToLower(runtime.GOOS) == "darwin" {
		return Process, nil
	}
	cmd := exec.Command("cat", cgroup)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		texts := strings.SplitN(scanner.Text(), ":", 3)
		if len(texts) < 3 {
			continue
		}
		key := strings.ToLower(texts[1])
		value := strings.ToLower(texts[2])
		if strings.Contains(key, devices) || strings.Contains(key, name) {
			if strings.HasPrefix(value, "/"+Docker) {
				return Docker, nil
			} else if strings.HasPrefix(value, "/"+Kubepods) {
				return Kubepods, nil
			}
		}
	}
	return Process, nil
}
