package ciscoucsconf

import (
	// "bytes"
	"context"
	"fmt"
	// "io"
	// "io/ioutil"
	// "net/url"
	"os"
	"path/filepath"
	// "strings"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/getperf/getperf2/common/sshx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	// "golang.org/x/crypto/ssh"
)

var (
	netAppSetCommand = "set cli output yaml\n"
)

// // ToDo: SSH接続タイムアウトの実装
// // net.Conn を介してアイドルタイムアウトの接続を作成する
// // Reference:
// // https://ja.coder.work/so/ssh/456706

// var (
// 	defaultTimeoutDuration = 100 * time.Second
// 	timeoutKillAfter       = 1 * time.Second
// )

// func getSshKey(keypath string) (ssh.Signer, error) {
// 	buf, err := ioutil.ReadFile(keypath)
// 	if err != nil {
// 		return nil, errors.Wrap(err, fmt.Sprintf("read key %s", keypath))
// 	}
// 	key, err := ssh.ParsePrivateKey(buf)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "private key parse")
// 	}
// 	return key, nil
// }

// func parseSshUrl(uri string) (string, string, error) {
// 	var ip, port string
// 	if !strings.HasPrefix(uri, "ssh://") {
// 		uri = "ssh://" + uri
// 	}
// 	u, err := url.Parse(uri)
// 	if err != nil {
// 		return ip, port, errors.Wrap(err, fmt.Sprintf("parse ssh url %s", uri))
// 	}
// 	ip = u.Hostname()
// 	port = u.Port()
// 	if port == "" {
// 		port = "22"
// 	}
// 	return ip, port, nil
// }

// func sshConnect(url, user, pass, keypath string) (*ssh.Client, error) {
// 	ip, port, err := parseSshUrl(url)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "prepare ssh connect")
// 	}
// 	auths := make([]ssh.AuthMethod, 0, 2)
// 	if pass != "" {
// 		auths = append(auths, ssh.Password(pass))
// 	}
// 	if keypath != "" {
// 		key, err := getSshKey(keypath)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "get ssh key")
// 		}
// 		auths = append(auths, ssh.PublicKeys(key))
// 	}
// 	config := &ssh.ClientConfig{
// 		User:            user,
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 		Auth:            auths,
// 	}

// 	conn, err := ssh.Dial("tcp", ip+":"+port, config)
// 	if err != nil {
// 		return nil, fmt.Errorf("ssh connect failed : %s", err)
// 	}
// 	log.Infof("connected : %s", url)
// 	return conn, nil
// }

// // 改行コードを統一する。
// func convNewline(str, nlcode string) string {
// 	return strings.NewReplacer(
// 		"\r\n", nlcode,
// 		"\r", nlcode,
// 		"\n", nlcode,
// 	).Replace(str)
// }

// // CIMC オプションにYAML出力設定追加
// func addCommandOption(str string) string {
// 	return "set cli output yaml\n" + str
// }

// func RunCommand(stdOut, stdErr io.Writer, conn *ssh.Client, execType ExecType, command string) error {
// 	session, err := conn.NewSession()
// 	if err != nil {
// 		return errors.Wrap(err, "prepare command")
// 	}
// 	defer session.Close()

// 	session.Stdout = stdOut
// 	session.Stderr = stdErr
// 	// 「予期しないファイル終了（EOF）」エラー回避のため、
// 	// 改行コードは LF に統一する

// 	command = addCommandOption(command)
// 	command = convNewline(command, "\n")
// 	if execType == "Cmd" || execType == "" {
// 		err = session.Run(command)
// 		if err != nil {
// 			return errors.Wrap(err, "run command")
// 		}
// 	} else if execType == "Script" {
// 		session.Stdin = bytes.NewBufferString(command + "\n")
// 		if err := session.Shell(); err != nil {
// 			return errors.Wrap(err, "run shell")
// 		}
// 		if err := session.Wait(); err != nil {
// 			return errors.Wrap(err, "run shell")
// 		}
// 	}

// 	return nil
// }

func (e *CiscoUCS) RunRemoteServer(ctx context.Context, env *cfg.RunEnv) error {
	log.Info("collect remote server : ", e.Server)
	e.datastore = filepath.Join(env.Datastore, e.Server)
	if err := os.MkdirAll(e.datastore, 0755); err != nil {
		return HandleError(e.errFile, err, "create log directory")
	}
	client, err := sshx.SshConnect(e.Url, e.User, e.Password, e.SshKeyPath)
	if err != nil {
		return HandleError(e.errFile, err, "connect remote server")
	}
	defer client.Close()
	for _, metric := range e.Metrics {
		if metric.Level == -1 || metric.Level > env.Level {
			continue
		}
		if metric.Id == "" || metric.Text == "" {
			continue
		}
		startTime := time.Now()
		outFile, err := env.OpenServerLog(e.Server, metric.Id)
		if err != nil {
			return HandleError(e.errFile, err, "prepare inventory log")
		}
		defer outFile.Close()
		cmd := netAppSetCommand + metric.Text
		if err := sshx.RunCommand(outFile, e.errFile, client, metric.Type, cmd); err != nil {
			HandleError(e.errFile, err, fmt.Sprintf("run %s:%s", e.Server, metric.Id))
		}
		log.Infof("run %s:%s,elapse %s", e.Server, metric.Id, time.Since(startTime))
	}
	return nil
}

func (e *CiscoUCS) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()
	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare CiscoUCS inventory error")
	}
	defer errFile.Close()
	e.errFile = errFile
	e.Env = env

	// if e.LocalExec == true {
	// 	log.Info("collect local server : ", e.LocalExec)
	// 	if err = e.RunLocalServer(ctx, env, e.Server); err != nil {
	// 		msg := fmt.Sprintf("run local server '%s'", e.Server)
	// 		HandleErrorWithAlert(e.errFile, err, msg)
	// 	}
	// }
	if err = e.RunRemoteServer(ctx, env); err != nil {
		msg := fmt.Sprintf("run remote server '%s'", e.Server)
		HandleErrorWithAlert(e.errFile, err, msg)
	}
	msg := fmt.Sprintf("Elapse %s", time.Since(startTime))
	log.Infof("Complete CiscoUCS inventory collection %s", msg)

	return err
}
