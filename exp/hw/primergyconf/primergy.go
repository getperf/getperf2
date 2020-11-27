package primergyconf

import (
	"context"
	"crypto/tls"

	// "encoding/json"
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
	// "github.com/tidwall/gjson"
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

func (e *Primergy) runSimple(ctx context.Context, client *resty.Client, metric *Metric) error {
	requestUrl := e.url + metric.Text
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetOutput(metric.Id).
		Get(requestUrl)
	if err != nil {
		return HandleError(e.errFile, err, metric.Text)
	}
	if code := resp.StatusCode(); code >= 400 {
		err := fmt.Errorf("get:%s,status code:%d", requestUrl, code)
		HandleError(e.errFile, err, metric.Id)
	}
	return nil
}

func (e *Primergy) runBatch(ctx context.Context, client *resty.Client, metric *Metric) error {
	log.Infof("run batch : %v\n", metric.Batch)
	return nil
}

func (e *Primergy) Run(ctx context.Context, env *cfg.RunEnv) error {
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

	e.url, err = parseUrl(e.Url)
	if err != nil {
		return HandleError(errFile, err, "init rest client")
	}
	e.errFile = errFile

	// requestUrl := url + "/rest/v1/Oem/eLCM/ProfileManagement/get?PARAM_PATH=Server/HWConfigurationIrmc/Adapters/RAIDAdapter"
	// fmt.Printf("url:%v\n", requestUrl)
	// resp, err := client.R().
	// 	SetHeader("Accept", "application/json").
	// 	Post(requestUrl)
	// if err != nil {
	// 	return HandleError(errFile, err, "metric.Text")
	// }
	// if code := resp.StatusCode(); code >= 400 {
	// 	err := fmt.Errorf("get:%s,status code:%d", requestUrl, code)
	// 	HandleError(errFile, err, "metric.Id")
	// }
	// // fmt.Printf("RESPONSE:%v\n", resp.Result)
	// fmt.Printf("RESPONSE2:%v\n", resp.String())
	// // bytes, err := json.Marshal(resp.Result)
	// // if err != nil {
	// // 	return HandleError(errFile, err, "decode json result")
	// // }
	// // fmt.Printf("RESPONSE:%v\n", string(bytes))
	// value := gjson.Get(resp.String(), "Session.Id").String()
	// fmt.Printf("Session.Id:%v\n", value)

	metrics = append(metrics, e.Metrics...)
	for _, metric := range metrics {
		if metric.Level > env.Level {
			continue
		}
		if metric.Id == "" || metric.Text == "" {
			continue
		}
		if metric.Batch == "" {
			e.runSimple(ctx, client, metric)
		} else {
			e.runBatch(ctx, client, metric)
		}
	}
	log.Infof("run %s:elapse %s", e.Server, time.Since(startTime))

	return nil
}
