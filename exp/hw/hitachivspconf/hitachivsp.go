package hitachivspconf

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
	"github.com/tidwall/gjson"
)

func parseUrl(uri string) (string, error) {
	if !strings.HasPrefix(uri, "http://") &&
		!strings.HasPrefix(uri, "https://") {
		uri = "http://" + uri
	}
	_, err := url.Parse(uri)
	if err != nil {
		return uri, errors.Wrapf(err, "parse url %s", uri)
	}
	return uri, nil
}

func (e *HitachiVSP) runSimple(ctx context.Context, client *resty.Client, url string, metricId string) error {
	requestUrl := e.url + url
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetOutput(metricId).
		Get(requestUrl)
	if err != nil {
		return HandleError(e.errFile, err, url)
	}
	if code := resp.StatusCode(); code >= 400 {
		err := fmt.Errorf("get:%s,status code:%d", requestUrl, code)
		HandleError(e.errFile, err, metricId)
	}
	return nil
}

func (e *HitachiVSP) Run(ctx context.Context, env *cfg.RunEnv) error {
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
	requestUrl := e.url + "/ConfigurationManager/v1/objects/storages"
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetOutput("storage").
		Get(requestUrl)
	if err != nil {
		return HandleError(e.errFile, err, requestUrl)
	}
	if code := resp.StatusCode(); code >= 400 {
		err := fmt.Errorf("get:%s,status code:%d", requestUrl, code)
		HandleError(e.errFile, err, "storage")
	}
	fmt.Printf("RESPONSE:%v\n", resp.Result)
	fmt.Printf("RESPONSE2:%v\n", resp.String())
	// bytes, err := json.Marshal(resp.Result)
	// if err != nil {
	// 	return HandleError(errFile, err, "decode json result")
	// }
	// fmt.Printf("RESPONSE:%v\n", string(bytes))
	// https://github.com/tidwall/gjson
	value := gjson.Get(resp.String(), "Session.Id").String()

	log.Infof("resp:%v", value)
	if e.Server != "" {
		// shift e.Server to target e.Servers
		e.Servers, e.Servers[0] =
			append(e.Servers[:1], e.Servers[0:]...), e.Server
	}
	for _, sv := range e.Servers {
		log.Infof("server:%v", sv)
	}

	// metrics = append(metrics, e.Metrics...)
	// for _, metric := range metrics {
	// 	if metric.Level > env.Level {
	// 		continue
	// 	}
	// 	if metric.Id == "" || metric.Text == "" {
	// 		continue
	// 	}
	// 	if metric.Batch == "" {
	// 		e.runSimple(ctx, client, metric)
	// 	} else {
	// 		e.runBatch(ctx, client, metric)
	// 	}
	// }
	log.Infof("run %s:elapse %s", e.Server, time.Since(startTime))

	return nil
}
