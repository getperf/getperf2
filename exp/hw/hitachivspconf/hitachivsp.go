package hitachivspconf

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
	"github.com/tidwall/gjson"
)

var debug = false

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

func (e *HitachiVSP) prepareClient() {
	e.client = resty.New().
		SetHeader("Accept", "application/json").
		RemoveProxy()
	if debug {
		e.client.SetDebug(true)
	}
	if e.Insecure {
		e.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
}

func (e *HitachiVSP) prepareNotAuthenticatedClient() error {
	if e.User == "" || e.Password == "" {
		err := fmt.Errorf("session info empty : %v", e)
		return HandleError(e.errFile, err, "prepare client")
	}
	e.prepareClient()
	e.client.SetBasicAuth(e.User, e.Password)
	return nil
}

func (e *HitachiVSP) prepareAuthorizedClient() error {
	if e.token == "" || e.sessionId == "" || e.storageDeviceId == "" {
		err := fmt.Errorf("session info empty : %v", e)
		return HandleError(e.errFile, err, "prepare authorized client")
	}
	e.prepareClient()
	e.client.SetHeader("Authorization", fmt.Sprintf("Session %s", e.token))
	return nil
}

func (e *HitachiVSP) createSession(ctx context.Context, storageDeviceId string) error {
	if err := e.prepareNotAuthenticatedClient(); err != nil {
		return HandleError(e.errFile, err, "create session")
	}
	url := e.url + "/ConfigurationManager/v1/objects/storages/" +
		storageDeviceId + "/sessions"
	resp, err := e.client.R().Post(url)
	if err != nil {
		return HandleError(e.errFile, err, "create session")
	}
	if code := resp.StatusCode(); code >= 400 {
		err := fmt.Errorf("POST %s, code:%d, body:%s", url, code, resp.String())
		return HandleError(e.errFile, err, "create session")
	}
	e.token = gjson.Get(resp.String(), "token").String()
	e.sessionId = gjson.Get(resp.String(), "sessionId").String()
	if e.token == "" || e.sessionId == "" {
		err := fmt.Errorf("POST %s, result:%s", url, resp.String())
		return HandleError(e.errFile, err, "create session")
	}
	e.storageDeviceId = storageDeviceId
	return nil
}

func (e *HitachiVSP) deleteSession(ctx context.Context) error {
	if err := e.prepareAuthorizedClient(); err != nil {
		return HandleError(e.errFile, err, "delete session")
	}
	url := e.url + "/ConfigurationManager/v1/objects/storages/" +
		e.storageDeviceId + "/sessions/" + e.sessionId
	resp, err := e.client.R().Delete(url)
	if err != nil {
		return HandleError(e.errFile, err, "delete session")
	}
	if code := resp.StatusCode(); code >= 400 {
		err := fmt.Errorf("delete %s,code : %d, body : %s", url, code, resp.String())
		return HandleError(e.errFile, err, "delete session")
	}
	return nil
}

func (e *HitachiVSP) getMetric(ctx context.Context, env *cfg.RunEnv,
	host string, metric string, requestUrl string) error {
	if err := e.prepareAuthorizedClient(); err != nil {
		return HandleError(e.errFile, err, metric)
	}
	requestUrl = strings.Replace(requestUrl, "{id}", e.storageDeviceId, -1)
	url := e.url + requestUrl
	resp, err := e.client.R().Get(url)
	if err != nil {
		return HandleError(e.errFile, err, metric)
	}
	if code := resp.StatusCode(); code >= 400 {
		err := fmt.Errorf("get %s,code : %d, body : %s", url, code, resp.String())
		return HandleError(e.errFile, err, metric)
	}
	logFile, err := env.OpenServerLog(host, metric)
	if err != nil {
		return HandleError(e.errFile, err, metric)
	}
	defer logFile.Close()
	logFile.Write(resp.Body())
	return nil
}

func existList(checkHost string, testHosts []string) string {
	for _, host := range testHosts {
		if checkHost == host {
			return host
		}
	}
	return ""
}

func (e *HitachiVSP) getSNMPConfigAndCheckHost(ctx context.Context, env *cfg.RunEnv,
	hosts []string) (string, error) {
	if err := e.prepareAuthorizedClient(); err != nil {
		return "", HandleError(e.errFile, err, "get snmp")
	}
	url := e.url + "/ConfigurationManager/v1/objects/storages/" +
		e.storageDeviceId + "/snmp-settings/instance"
	resp, err := e.client.R().Get(url)
	if err != nil {
		return "", HandleError(e.errFile, err, "get snmp")
	}
	if code := resp.StatusCode(); code >= 400 {
		err := fmt.Errorf("get %s,code : %d, body : %s", url, code, resp.String())
		return "", HandleError(e.errFile, err, "get snmp")
	}
	storageSystemName := gjson.Get(resp.String(), "systemGroupInformation.storageSystemName").String()
	log.Infof("host : %v", storageSystemName)
	host := existList(storageSystemName, hosts)
	if host == "" {
		return "", nil
	}
	logFile, err := env.OpenServerLog(host, "snmp")
	if err != nil {
		return "", HandleError(e.errFile, err, "get snmp")
	}
	defer logFile.Close()
	logFile.Write(resp.Body())

	return host, nil
}

func (e *HitachiVSP) getStorageDeviceIds(ctx context.Context, env *cfg.RunEnv) ([]string, error) {
	var storageDeviceIds []string
	if err := e.prepareNotAuthenticatedClient(); err != nil {
		return storageDeviceIds, HandleError(e.errFile, err, "get storages")
	}
	url := e.url + "/ConfigurationManager/v1/objects/storages"
	resp, err := e.client.R().Get(url)
	if err != nil {
		return storageDeviceIds, HandleError(e.errFile, err, url)
	}
	if code := resp.StatusCode(); code >= 400 {
		err := fmt.Errorf("get:%s,status code:%d", url, code)
		return storageDeviceIds, HandleError(e.errFile, err, "get storages")
	}

	logFile, err := env.OpenLog("storages")
	if err != nil {
		return storageDeviceIds, HandleError(e.errFile, err, "get storages")
	}
	defer logFile.Close()
	logFile.Write(resp.Body())
	results := gjson.Get(resp.String(), "data.#.storageDeviceId").Array()
	if len(results) == 0 {
		err := fmt.Errorf("get storage dievces is empty")
		return storageDeviceIds, HandleError(e.errFile, err, "get storages")
	}
	for _, id := range results {
		storageDeviceIds = append(storageDeviceIds, id.String())
	}
	return storageDeviceIds, nil
}

func (e *HitachiVSP) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()

	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare error.log")
	}
	defer errFile.Close()

	e.url, err = parseUrl(e.Url)
	if err != nil {
		return HandleError(errFile, err, "prepare rest url")
	}
	e.errFile = errFile

	storageDeviceIds, err := e.getStorageDeviceIds(ctx, env)
	if err != nil {
		return HandleError(errFile, err, "get storage drvices")
	}

	var servers []string
	// shift e.Server to target e.Servers
	if e.Server != "" {
		servers = append(servers, e.Server)
	}
	servers = append(servers, e.Servers...)
	for _, server := range servers {
		log.Infof("server:%v", server)
		datastore := filepath.Join(env.Datastore, server)
		if err := os.MkdirAll(datastore, 0755); err != nil {
			return HandleError(errFile, err, "create log directory")
		}
	}

	for _, storageDeviceId := range storageDeviceIds {
		err = e.createSession(ctx, storageDeviceId)
		if err != nil {
			return HandleError(errFile, err, "prepare storage session")
		}
		log.Infof("deviceId:%v, %v, %v", e.storageDeviceId, e.token, e.sessionId)

		host, err := e.getSNMPConfigAndCheckHost(ctx, env, servers)
		if err != nil {
			return HandleError(errFile, err, "get snmp to check host")
		}
		if host != "" {
			log.Infof("exist:%v", host)
			metrics = append(metrics, e.Metrics...)
			for _, metric := range metrics {
				if metric.Level > env.Level {
					continue
				}
				if metric.Id == "" || metric.Text == "" {
					continue
				}
				e.getMetric(ctx, env, host, metric.Id, metric.Text)
			}
		}
		err = e.deleteSession(ctx)
		if err != nil {
			return HandleError(errFile, err, "close storage session")
		}
	}
	log.Infof("run %s:elapse %s", e.Server, time.Since(startTime))

	return nil
}
