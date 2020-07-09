package gcmain

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/getperf/getperf2/agent"
	"github.com/getperf/getperf2/cfg"
	_ "github.com/getperf/getperf2/exp/all"

	// . "github.com/getperf/getperf2/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type InventoryRetriever struct {
	Base          *cfg.ExportBase
	Common        *cfg.CommonConfig
	Env           *cfg.RunEnv
	WorkDir       string
	DatastoreSets []agent.DatastoreSet
	HttpClient    *http.Client
}

func NewInventoryRetriever(base *cfg.ExportBase, common *cfg.CommonConfig, env *cfg.RunEnv) *InventoryRetriever {
	if base == nil || common == nil || env == nil {
		return nil
	}
	return &InventoryRetriever{
		Base:   base,
		Common: common,
		Env:    env,

		DatastoreSets: []agent.DatastoreSet{},
	}
}

func (c *InventoryRetriever) Validate() error {
	if err := c.Common.Check(); err != nil {
		return errors.Wrap(err, "check config")
	}
	if err := c.Env.Check(); err != nil {
		return errors.Wrap(err, "check config")
	}
	if c.Env.Retrieve.FromUrl == "" {
		return errors.New("--from-url must specifiend")
	}
	return nil
}

func CreateHttpClient(conf *cfg.RetrieveConfig) (*http.Client, error) {
	SetHostIpLookupTransport()
	client := &http.Client{}
	log.Infof("client cert : %s", conf.ClientCertPath)
	log.Infof("CA cert : %s", conf.CaPath)
	if conf.ClientCertPath != "" && conf.CaPath != "" {
		cert, err := tls.LoadX509KeyPair(conf.ClientCertPath, conf.ClientCertPath)
		if err != nil {
			return client, errors.Wrap(err, "create http client")
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(conf.CaPath)
		if err != nil {
			return client, errors.Wrap(err, "create http client")
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}
	return client, nil
}

func (c *InventoryRetriever) GetDatastoreSets() error {
	urlBase := c.Env.Retrieve.FromUrl
	urlDownloadSets := urlBase + "/store"

	// client, err := CreateHttpClient(c.Env.Retrieve)
	client := c.HttpClient

	request, err := http.NewRequest("GET", urlDownloadSets, nil)
	if err != nil {
		return errors.Wrap(err, "create download keys request")
	}
	resp, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "post get download keys")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read download keys response")
	}
	err = json.Unmarshal(body, &c.DatastoreSets)
	if err != nil {
		return errors.Wrapf(err, "json parse download keys response '%v'", string(body))
	}
	return nil
}

func (c *InventoryRetriever) UnzipDatastore(dsSet agent.DatastoreSet) error {
	zipPath := filepath.Join(c.WorkDir, dsSet.ZipFile)
	log.Infof("unzip %s", dsSet.ZipFile)
	return agent.Unzip(zipPath, c.Env.Datastore)
}

func (c *InventoryRetriever) RetrieveDatastore(dsSet agent.DatastoreSet) error {
	baseUrl := c.Env.Retrieve.FromUrl
	downloadUrl := fmt.Sprintf("%s/zip/%s", baseUrl, dsSet.ZipFile)

	// client, err := CreateHttpClient(c.Env.Retrieve)
	client := c.HttpClient

	request, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return errors.Wrap(err, "create download request")
	}
	response, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "post download request")
	}
	defer response.Body.Close()

	zipPath := filepath.Join(c.WorkDir, dsSet.ZipFile)
	zipSaveFile, err := os.Create(zipPath)
	if err != nil {
		return errors.Wrap(err, "pepare zip file")
	}
	defer zipSaveFile.Close()

	_, err = io.Copy(zipSaveFile, response.Body)
	if err != nil {
		return errors.Wrap(err, "read download keys response")
	}
	return nil
}

// TODO:
// エージェントの接続先は、https://{ホスト名}:59443/ 形式で
// ホスト名を指定する必要がある。ホスト名からIPの名前解決が必要
// Transport.DialContext をカスタマイズする記事があり、動作検証中
// 他にネームサービスにIPを登録する方法を調査中
//
// リファレンス：
//
// https://stackoverflow.com/questions/40624248/golang-force-http-request-to-specific-ip-similar-to-curl-resolve
// You can provide a custom Transport.DialContext function.
// https://golang.org/pkg/net/http/#Transport

func SetHostIpLookupTransport() error {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	// or create your own transport, there's an example on godoc.
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		log.Info("address original =", addr)
		if addr == "centos80:59443" {
			addr = "192.168.0.5:59443"
			log.Info("address modified =", addr)
		}
		return dialer.DialContext(ctx, network, addr)
	}
	return nil
}

// POST url+/downloadkeys/  して、host と stat_name のリストを取得
// リストを順に実行
// 	POST url+/download/{host}/{stat_name} をしてzip ファイル取得
// 	取得したzipから 日付ディレクトリのパスを除いて解凍

func (c *InventoryRetriever) Run() error {
	workDir, err := ioutil.TempDir("", "ptune")
	if err != nil {
		return errors.Wrap(err, "prepare retriever temp")
	}
	defer os.RemoveAll(workDir)
	c.WorkDir = workDir

	log.Info("TESTTESTTEST")
	log.Info("set inventory datastore : ", c.Env.Datastore)
	log.Info("set work dir : ", c.WorkDir)
	log.Info("get inventory url : ", c.Env.Retrieve.FromUrl)

	// if err := c.SetHostIpLookupTransport(); err != nil {
	// 	return errors.Wrap(err, "prepare ip lookup")
	// }
	client, err := CreateHttpClient(c.Env.Retrieve)
	if err != nil {
		return errors.Wrap(err, "prepare http client")
	}
	c.HttpClient = client
	if err := c.GetDatastoreSets(); err != nil {
		return errors.Wrap(err, "request find download list")
	}
	log.Infof("datastore found : %d", len(c.DatastoreSets))
	for _, datastoreSet := range c.DatastoreSets {
		if err := c.RetrieveDatastore(datastoreSet); err != nil {
			log.Error(err)
		}
		if err := c.UnzipDatastore(datastoreSet); err != nil {
			log.Error(err)
		}
	}

	return nil
}
