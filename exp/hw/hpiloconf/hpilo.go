package hpiloconf

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func parseUrl(uri string) (string, error) {
	if !strings.HasPrefix(uri, "http://") &&
		!strings.HasPrefix(uri, "https://") {
		uri = "https://" + uri
	}
	_, err := url.Parse(uri)
	if err != nil {
		return uri, errors.Wrapf(err, "parse url %s", uri)
	}
	return uri, nil
}

func (e *HPiLO) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()

	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare error.log")
	}
	defer errFile.Close()

	datastore := filepath.Join(env.Datastore, e.Server)
	if err := os.MkdirAll(datastore, 0755); err != nil {
		return HandleError(errFile, err, "create log directory")
	}

	client := resty.New().
		SetBasicAuth(e.User, e.Password).
		SetOutputDirectory(datastore).
		RemoveProxy()

	if e.Insecure { // 自己証明書の許可 `true`
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	url, err := parseUrl(e.Url)
	if err != nil {
		return HandleError(errFile, err, "init rest client")
	}

	for _, command := range commands {
		if command.Level > env.Level {
			continue
		}
		if command.Id == "" {
			continue
		}
		requestUrl := url + command.Text
		resp, err := client.R().
			SetOutput(command.Id).
			Get(requestUrl)
		if err != nil {
			return HandleError(errFile, err, command.Text)
		}
		if code := resp.StatusCode(); code >= 400 {
			err := fmt.Errorf("get:%s,status code:%d", requestUrl, code)
			HandleError(errFile, err, command.Id)
		}
	}
	log.Infof("run %s:elapse %s", e.Server, time.Since(startTime))

	return nil
}
