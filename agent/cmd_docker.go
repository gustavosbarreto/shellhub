// +build docker

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func newCmd(u *User, shell, term, host string, command ...string) *exec.Cmd {
	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)

	nscommand, _ := nsenterCommandWrapper(uid, gid, fmt.Sprintf("/host/%s", u.HomeDir), command...)

	cmd := exec.Command(nscommand[0], nscommand[1:]...)
	cmd.Env = []string{
		"TERM=" + term,
		"HOME=" + u.HomeDir,
		"SHELL=" + shell,
		"SHELLHUB_HOST=" + host,
	}

	return cmd
}

func nsenterCommandWrapper(uid, gid int, home string, command ...string) ([]string, error) {
	wrappedCommand := []string{}

	if _, err := os.Stat("/usr/bin/nsenter"); err == nil {
		wrappedCommand = append([]string{
			"/usr/bin/nsenter",
			"-t", "1",
			"-a",
			"-S", strconv.Itoa(uid),
			"-G", strconv.Itoa(gid),
			fmt.Sprintf("--wd=%s", home),
		}, wrappedCommand...)
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	wrappedCommand = append(wrappedCommand, command...)

	return wrappedCommand, nil
}
