package telnetx

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/ziutek/telnet"
)

const timeout = 100 * time.Second

type Telnetx struct {
	t          *telnet.Conn
	promptSize int
}

func expect(t *telnet.Conn, d ...string) error {
	if err := t.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return errors.Wrap(err, "expect read")
	}
	if err := t.SkipUntil(d...); err != nil {
		return errors.Wrap(err, "expect wait")
	}
	return nil
}

func sendln(t *telnet.Conn, s string) error {
	if err := t.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return errors.Wrap(err, "send start")
	}
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	if _, err := t.Write(buf); err != nil {
		return errors.Wrap(err, "send end")
	}
	return nil
}

func TelnetConnect(uri, user, passwd string) (*Telnetx, error) {
	return nil, fmt.Errorf("Windows environment does not support telnet")
}

func TelnetConnectOld(uri, user, passwd string) (*Telnetx, error) {

	t, err := telnet.Dial("tcp", uri)
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	t.SetUnixWriteMode(true)
	t.SetEcho(false)
	log.Debugf("login:%s\n", user)
	if err := expect(t, "login: "); err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	if err := sendln(t, user); err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	if err := expect(t, "ssword: "); err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	if err := sendln(t, passwd); err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	if err := expect(t, "% ", "$ "); err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	if err := sendln(t, "bash"); err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	if err := expect(t, "$ "); err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	if err := sendln(t, ""); err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	data, err := t.ReadUntil("$ ")
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	return &Telnetx{t, len(data)}, err
}

func (x *Telnetx) Close() {
	t := x.t
	sendln(t, "exit")
}

func (x *Telnetx) ExecCommand(cmd string) (string, error) {
	t := x.t
	var data []byte
	fmt.Printf("cmd:%d:%s\n", len(cmd), cmd)
	if err := sendln(t, cmd); err != nil {
		return "", errors.Wrap(err, "command")
	}
	// data, err = t.ReadBytes('%')
	data, err := t.ReadUntil("% ", "$ ")
	if err != nil {
		return "", errors.Wrap(err, "command")
	}
	return string(data[len(cmd) : len(data)-x.promptSize]), nil
}

func (x *Telnetx) ExecScript(cmd string) (string, error) {
	t := x.t
	var data []byte

	cmd = strings.TrimRight(cmd, "\r\n")
	cmd = convNewline(cmd, "\n")
	cmdLines := strings.Split(cmd, "\n")
	for no, cmdLine := range cmdLines {
		fmt.Printf("line:%d,%d:%s\n", no, len(cmdLines), cmdLine)
		if err := sendln(t, cmdLine); err != nil {
			return "", errors.Wrap(err, "script")
		}
		if no < len(cmdLines) {
			if err := expect(t, "> "); err != nil {
				return "", errors.Wrap(err, "script")
			}
		}
	}
	data, err := t.ReadUntil("% ", "$ ")
	if err != nil {
		return "", errors.Wrap(err, "script")
	}
	return string(data[0 : len(data)-x.promptSize]), nil
}
