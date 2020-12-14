package netappconf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// ToDo: SSH接続タイムアウトの実装
// net.Conn を介してアイドルタイムアウトの接続を作成する
// Reference:
// https://ja.coder.work/so/ssh/456706

var (
	defaultTimeoutDuration = 100 * time.Second
	timeoutKillAfter       = 1 * time.Second
	netAppSetCommand       = `set -showallfields true -rows 0 -showseparator "<|>" -units GB;`
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

func sshConnect(url, user, pass, keypath string) (*ssh.Client, error) {
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

// func (e *NetAPP) RunRemoteServer(ctx context.Context, env *cfg.RunEnv) error {
// 	log.Info("collect remote server : ", e.Server)
// 	e.datastore = filepath.Join(env.Datastore, e.Server)
// 	if err := os.MkdirAll(e.datastore, 0755); err != nil {
// 		return HandleError(e.errFile, err, "create log directory")
// 	}
// 	client, err := sshConnect(e.Url, e.User, e.Password, e.SshKeyPath)
// 	if err != nil {
// 		return HandleError(e.errFile, err, "connect remote server")
// 	}
// 	defer client.Close()
// 	for _, metric := range e.Metrics {
// 		if metric.Level == -1 || metric.Level > env.Level {
// 			continue
// 		}
// 		if metric.Id == "" || metric.Text == "" {
// 			continue
// 		}
// 		startTime := time.Now()
// 		outFile, err := env.OpenServerLog(e.Server, metric.Id)
// 		if err != nil {
// 			return HandleError(e.errFile, err, "prepare inventory log")
// 		}
// 		defer outFile.Close()
// 		log.Infof("run %s %s %s", outFile, time.Since(startTime))
// 		// if err := RunCommand(outFile, e.errFile, client, metric.Type, metric.Text); err != nil {
// 		// 	HandleError(e.errFile, err, fmt.Sprintf("run %s:%s", e.Server, metric.Id))
// 		// }
// 		// log.Infof("run %s:%s,elapse %s", e.Server, metric.Id, time.Since(startTime))
// 	}
// 	return nil
// }

func (e *NetAPP) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()
	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare NetAPP inventory error")
	}
	defer errFile.Close()
	e.errFile = errFile

	var servers []string
	if e.Server != "" {
		servers = append(servers, e.Server)
	}
	servers = append(servers, e.Servers...)
	for _, server := range servers {
		datastore := filepath.Join(env.Datastore, server)
		if err := os.MkdirAll(datastore, 0755); err != nil {
			return HandleError(errFile, err, "create log directory")
		}
	}
	client, err := sshConnect(e.Url, e.User, e.Password, "")
	if err != nil {
		return HandleError(e.errFile, err, "connect NetAPP management server")
	}
	defer client.Close()

	metrics = append(metrics, e.Metrics...)
	for _, metric := range metrics {
		if metric.Level > env.Level {
			continue
		}
		if metric.Id == "" || metric.Text == "" {
			continue
		}
		if !metric.Remote {
			log.Infof("get metric: %s", metric.Id)
			logFile, err := env.OpenLog(metric.Id)
			if err != nil {
				return HandleError(e.errFile, err, metric.Id)
			}
			defer logFile.Close()
			cmd := netAppSetCommand + metric.Text
			if err := RunCommand(logFile, e.errFile, client, metric.Type, cmd); err != nil {
				HandleError(e.errFile, err, metric.Id)
			}
		} else {
			for _, server := range servers {
				log.Infof("get metric: %s, node : %s", metric.Id, server)
				logFile, err := env.OpenServerLog(server, metric.Id)
				if err != nil {
					return HandleError(e.errFile, err, metric.Id)
				}
				defer logFile.Close()
				cmd := netAppSetCommand + metric.Text
				cmd = strings.Replace(cmd, "{host}", server, -1)
				if err := RunCommand(logFile, e.errFile, client, metric.Type, cmd); err != nil {
					HandleError(e.errFile, err, metric.Id)
				}
			}
		}
	}
	msg := fmt.Sprintf("Elapse %s", time.Since(startTime))
	log.Infof("Complete NetAPP inventory collection %s", msg)

	return err
}
