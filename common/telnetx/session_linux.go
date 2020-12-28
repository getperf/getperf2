package telnetx

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	expect "github.com/google/goexpect"
	"github.com/pkg/errors"
	"github.com/ziutek/telnet"
)

const timeout = 100 * time.Second

type Telnetx struct {
	exp    expect.Expecter
	prompt string
}

var (
	userRE   = regexp.MustCompile("login:")
	passRE   = regexp.MustCompile("ssword:")
	promptRE = regexp.MustCompile("[%$] ")
	shellRE  = regexp.MustCompile("> ")
)

// func expect(t *telnet.Conn, d ...string) error {
// 	if err := t.SetReadDeadline(time.Now().Add(timeout)); err != nil {
// 		return errors.Wrap(err, "expect read")
// 	}
// 	// if err := t.SkipUntil(d...); err != nil {
// 	// 	return errors.Wrap(err, "expect wait")
// 	// }
// 	data, err := t.ReadUntil(d...)
// 	if err != nil {
// 		return errors.Wrap(err, "expect wait")
// 	}
// 	fmt.Printf("expect:%v\n", string(data))
// 	return nil
// }

// func sendln(t *telnet.Conn, s string) error {
// 	if err := t.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
// 		return errors.Wrap(err, "send start")
// 	}
// 	buf := make([]byte, len(s)+1)
// 	copy(buf, s)
// 	buf[len(s)] = '\n'
// 	if _, err := t.Write(buf); err != nil {
// 		return errors.Wrap(err, "send end")
// 	}
// 	return nil
// }

func telnetSpawn(addr string, timeout time.Duration, opts ...expect.Option) (expect.Expecter, <-chan error, error) {
	conn, err := telnet.Dial("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	resCh := make(chan error)

	return expect.SpawnGeneric(&expect.GenOptions{
		In:  conn,
		Out: conn,
		Wait: func() error {
			return <-resCh
		},
		Close: func() error {
			close(resCh)
			return conn.Close()
		},
		Check: func() bool { return true },
	}, timeout, opts...)
}

func TelnetConnect(uri, user, passwd string) (*Telnetx, error) {
	// fmt.Println(term.Bluef("Telnet spawner connect"))
	// exp, _, err := telnetSpawn(uri, timeout, expect.Verbose(true))
	exp, _, err := telnetSpawn(uri, timeout)
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	exp.Expect(userRE, timeout)
	exp.Send(user + "\n")
	exp.Expect(passRE, timeout)
	exp.Send(passwd + "\n")
	exp.Expect(promptRE, timeout)
	exp.Send("bash\n")
	exp.Expect(promptRE, timeout)
	exp.Send("\n")
	result, _, err := exp.Expect(promptRE, timeout)
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	return &Telnetx{exp, result}, err
}

func (x *Telnetx) Close() {
	exp := x.exp
	exp.Send("exit\n")
}

func (x *Telnetx) ExecCommand(cmd string) (string, error) {
	exp := x.exp
	cmd = strings.TrimRight(cmd, "\r\n")
	exp.Send(cmd + "\n")
	result, _, err := exp.Expect(promptRE, timeout)
	if err != nil {
		return "", errors.Wrap(err, "command")
	}
	result = strings.ReplaceAll(result, x.prompt, "")
	return result, nil
}

func (x *Telnetx) ExecScript(cmd string) (string, error) {
	exp := x.exp
	cmd = strings.TrimRight(cmd, "\r\n")
	cmd = convNewline(cmd, "\n")
	cmdLines := strings.Split(cmd, "\n")
	for no, cmdLine := range cmdLines {
		fmt.Printf("line:%d,%d:%s\n", no, len(cmdLines), cmdLine)
		exp.Send(cmdLine + "\n")
		if no < len(cmdLines)-1 {
			exp.Expect(shellRE, timeout)
		}
	}
	result, _, err := exp.Expect(promptRE, timeout)
	if err != nil {
		return "", errors.Wrap(err, "script")
	}
	// return data[0 : len(data)-x.prompt], nil
	return result, nil
}
