package sshx

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type ExecType string

const (
	Cmd    = ExecType("Cmd")
	Script = ExecType("Script")
)

var (
	defaultTimeoutDuration = 100 * time.Second
	timeoutKillAfter       = 1 * time.Second
)

func getSshKey(keypath string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(keypath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("read key %s", keypath))
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, errors.Wrap(err, "private key parse")
	}
	return key, nil
}

func parseSshUrl(uri string) (string, string, error) {
	var ip, port string
	if !strings.HasPrefix(uri, "ssh://") {
		uri = "ssh://" + uri
	}
	u, err := url.Parse(uri)
	if err != nil {
		return ip, port, errors.Wrap(err, fmt.Sprintf("parse ssh url %s", uri))
	}
	ip = u.Hostname()
	port = u.Port()
	if port == "" {
		port = "22"
	}
	return ip, port, nil
}

func SshConnect(url, user, pass, keypath string) (*ssh.Client, error) {
	ip, port, err := parseSshUrl(url)
	if err != nil {
		return nil, errors.Wrap(err, "prepare ssh connect")
	}
	auths := make([]ssh.AuthMethod, 0, 2)
	if pass != "" {
		auths = append(auths, ssh.Password(pass))
	}
	if keypath != "" {
		key, err := getSshKey(keypath)
		if err != nil {
			return nil, errors.Wrap(err, "get ssh key")
		}
		auths = append(auths, ssh.PublicKeys(key))
	}
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            auths,
	}

	conn, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		return nil, fmt.Errorf("ssh connect failed : %s", err)
	}
	log.Infof("connected : %s", url)
	return conn, nil
}

// 改行コードを統一する。
func convNewline(str, nlcode string) string {
	return strings.NewReplacer(
		"\r\n", nlcode,
		"\r", nlcode,
		"\n", nlcode,
	).Replace(str)
}

func RunCommand(stdOut, stdErr io.Writer, conn *ssh.Client, execType ExecType, command string) error {
	session, err := conn.NewSession()
	if err != nil {
		return errors.Wrap(err, "prepare command")
	}
	defer session.Close()

	session.Stdout = stdOut
	session.Stderr = stdErr
	// 「予期しないファイル終了（EOF）」エラー回避のため、
	// 改行コードは LF に統一する
	command = convNewline(command, "\n")
	// fmt.Printf("command:%v,type:%v\n", command, execType)
	if execType == "Cmd" || execType == "" {
		err = session.Run(command)
		if err != nil {
			return errors.Wrap(err, "run command")
		}
	} else if execType == "Script" {
		session.Stdin = bytes.NewBufferString(command + "\n")
		if err := session.Shell(); err != nil {
			return errors.Wrap(err, "run shell")
		}
		if err := session.Wait(); err != nil {
			return errors.Wrap(err, "run shell")
		}
	}

	return nil
}
