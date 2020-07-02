package agent

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func RunNettest(ctx context.Context, argv []string, stdout, stderr io.Writer) error {
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	var (
		configPath = flag.String("config", "", "It performs by the specified directory.")
		serviceUrl = flag.String("u", "", "dowonload service url")
		backupDir  = flag.String("b", "testdata/webservice/backup", "dowonload dir")
		tlsConfig  = flag.String("t", "", "TLS config file")
	)
	flag.Parse()

	hostName, err := GetHostname()
	if err != nil {
		log.Fatal("get hostname ", err)
	}
	if *configPath == "" {
		home, err := GetParentAbsPath(os.Args[0], 2)
		if err != nil {
			log.Fatal("get getperf path ", err)
		}
		*configPath = filepath.Join(home, "getperf.ini")
	}
	configEnv := NewConfigEnvBase(hostName, cmdName, *configPath)
	home, err := GetParentAbsPath(*configPath, 1)
	if err != nil {
		log.Fatal("get home ", err)
	}
	cfg := NewConfig(home, configEnv)
	cfg.InitAgent()
	cfg.ParseConfigFile(cfg.SvParamFile)

	log.Infof("url : %s", cfg.Schedule.WebServiceUrl)
	if len(*serviceUrl) == 0 {
		*serviceUrl = cfg.Schedule.WebServiceUrl
	}
	log.Infof("url : %s, dir : %s, config : %s", *serviceUrl, *backupDir, *tlsConfig)
	webService := NewWebServiceBase(*serviceUrl, *backupDir, cfg)

	u, err := url.Parse(webService.ServiceUrl)
	if err != nil {
		return errors.Wrap(err, "exporter service, pasing url")
	}

	router := mux.NewRouter()
	router.HandleFunc("/store/{node}/{job}", webService.DownloadDatastoreZip)
	router.HandleFunc("/store/{node}/{job}/{since}", webService.DownloadDatastoreZip)
	router.HandleFunc("/store", webService.GetLatestDatastoreSets)

	router.PathPrefix("/zip/").Handler(http.StripPrefix("/zip/",
		http.FileServer(http.Dir(*backupDir))))

	server := &http.Server{
		Addr:              u.Host,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	if u.Scheme == "https" {
		server.TLSConfig, err = ConfigToTLSConfig(cfg)
		if err != nil {
			return errors.Wrap(err, "load tls config")
		}
		server.TLSConfig.GetConfigForClient = func(*tls.ClientHelloInfo) (*tls.Config, error) {
			return ConfigToTLSConfig(cfg)
		}
		if err := server.ListenAndServeTLS("", ""); err != nil {
			return errors.Wrap(err, "exporter listen and serve")
		}

	} else {
		if err := server.ListenAndServe(); err != nil {
			return errors.Wrap(err, "exporter listen and serve")
		}
	}

	return nil
}
