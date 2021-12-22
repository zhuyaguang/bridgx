package utils

import (
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// SshCheck 测试目的机器SSH连通性
func SshCheck(ip, user, pwd string) bool {
	addr := strings.TrimSpace(ip) + ":22"
	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            strings.TrimSpace(user),
		Auth:            []ssh.AuthMethod{ssh.Password(strings.TrimSpace(pwd))},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 3,
	})
	if err != nil {
		return false
	}
	session, err := client.NewSession()
	if err != nil {
		return false
	}
	session.Close()
	return true
}
