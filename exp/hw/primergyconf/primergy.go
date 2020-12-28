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
	"github.com/tidwall/gjson"
)

var debug = false

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

	// 1. プロファイル生成リクエスト
	// POST rest/v1/Oem/eLCM/ProfileManagement/get?PARAM_PATH=Server/HWConfigurationIrmc/Adapters/RAIDAdapter

	makeProfileUrl := e.url + metric.Text
	resp, err := client.R().Post(makeProfileUrl)
	if err != nil {
		return HandleError(e.errFile, err, metric.Text)
	}

	// 409コンフリクトで既に作成ずみの場合は次の終了待ちをスキップ
	profileExist := false
	if code := resp.StatusCode(); code >= 400 {
		if code == 409 {
			profileExist = true
		} else {
			err := fmt.Errorf("post:%s,status code:%d", makeProfileUrl, code)
			return HandleError(e.errFile, err, metric.Id)
		}
	}

	sessionId := ""
	if !profileExist {
		// 2. レスポンスからセッションID取得
		sessionId = gjson.Get(resp.String(), "Session.Id").String()
		if sessionId == "" {
			return HandleError(e.errFile, fmt.Errorf("not found Session.Id"), metric.Id)
		}

		// 3. セッションステータスを検索してプロファイル生成完了を待つ
		sessionComplete := false
		getSessionUrl := e.url + fmt.Sprintf("/sessionInformation/%s/status", sessionId)
		for i := 1; i < 60; i++ {
			time.Sleep(time.Second * 5)
			// GET sessionInformation/${session_id}/status
			resp, err := client.R().Get(getSessionUrl)
			if err != nil {
				return HandleError(e.errFile, err, metric.Id)
			}
			if code := resp.StatusCode(); code >= 400 {
				err := fmt.Errorf("get:%s,status code:%d", getSessionUrl, code)
				return HandleError(e.errFile, err, metric.Id)
			}
			sessionStatus := gjson.Get(resp.String(), "Session.Status").String()
			// 'running'の状態ならスリープして繰り返す
			if strings.ToLower(sessionStatus) != "running" {
				sessionComplete = true
				profileExist = true
				break
			}
		}
		if !sessionComplete {
			return HandleError(e.errFile, fmt.Errorf("profile session timeout"), metric.Id)
		}
	}

	// 4. プロファイル取得して、結果保存

	// GET /rest/v1/Oem/eLCM/ProfileManagement/RAIDAdapter
	getProfileUrl := e.url + metric.Batch
	resp, err = client.R().SetOutput(metric.Id).Get(getProfileUrl)
	if err != nil {
		return HandleError(e.errFile, err, metric.Text)
	}
	if code := resp.StatusCode(); code >= 400 {
		err := fmt.Errorf("get:%s,status code:%d", getProfileUrl, code)
		return HandleError(e.errFile, err, metric.Id)
	}

	// 5. セッション削除。セッションを取得した場合

	// DELETE /sessionInformation/${session_id}/remove
	if sessionId != "" {
		deleteSessionUrl := e.url + fmt.Sprintf("/sessionInformation/%s/remove", sessionId)
		resp, err = client.R().Delete(deleteSessionUrl)
	}

	// // 6. プロファイル削除

	// DELETE /rest/v1/Oem/eLCM/ProfileManagement/RAIDAdapter
	if profileExist {
		resp, err = client.R().Delete(getProfileUrl)
		if err != nil {
			return HandleError(e.errFile, err, metric.Text)
		}
		if code := resp.StatusCode(); code >= 400 {
			err := fmt.Errorf("delete:%s,status code:%d", getProfileUrl, code)
			return HandleError(e.errFile, err, metric.Id)
		}
	}

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
		SetHeader("Accept", "application/json").
		SetOutputDirectory(datastore).
		RemoveProxy()

	if e.Insecure { // 自己証明書の許可 `true`
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if debug {
		client.SetDebug(true)
	}

	e.url, err = parseUrl(e.Url)
	if err != nil {
		return HandleError(errFile, err, "init rest client")
	}
	e.errFile = errFile

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
