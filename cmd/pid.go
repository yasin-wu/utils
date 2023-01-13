package cmd

import (
	"os/exec"
	"strings"
)

func GetPID(name string) (string, error) {
	arg := `ps ux | awk '/` + name + `/ && !/awk/ {print $2}'`
	pid, err := exec.Command("/bin/sh", "-c", arg).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(pid)), nil
}
