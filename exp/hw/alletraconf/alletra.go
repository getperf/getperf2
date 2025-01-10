package alletraconf

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	// . "github.com/go-resty/resty/v2"
    "github.com/hpe-storage/nimble-golang-sdk/pkg/service"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (e *Alletra) saveJson(ioErr io.Writer, outfile string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
        fmt.Fprintf(ioErr, err.Error())
		return
	}
	outPath := filepath.Join(e.datastore, outfile)
    if err := ioutil.WriteFile(outPath, []byte(jsonData), 0666); err != nil {
        fmt.Fprintf(ioErr, err.Error())
    }
}

func (e *Alletra) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()

	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare error.log")
	}
	defer errFile.Close()

	e.datastore = filepath.Join(env.Datastore, e.Server)
	if err := os.MkdirAll(e.datastore, 0755); err != nil {
		return HandleError(errFile, err, "create log directory")
	}

    groupService, err := service.NewNsGroupService(
        e.Url,
        e.User,
        e.Password,
        "v1",
        true)
	if err != nil {
		return HandleError(errFile, err, "create nimble sdk client service")
	}
    defer groupService.LogoutService()

	for _, metric := range metrics {
		if metric.Level > env.Level {
			continue
		}
		if metric.Id == "" {
			continue
		}
		if metric.Id == "arrays" {
			arrays, err := groupService.GetArrayService().GetArrays(nil)
			if err != nil {
				HandleError(errFile, err, metric.Text)
				continue
			}
			e.saveJson(errFile, "arrays", arrays)

		} else if metric.Id == "disks" {
			disks, err := groupService.GetDiskService().GetDisks(nil)
			if err != nil {
				HandleError(errFile, err, metric.Text)
				continue
			} 
			e.saveJson(errFile, "disks", disks)

		} else if metric.Id == "netconfig" {
			netconfig, err := groupService.GetNetworkConfigService().GetNetworkConfigs(nil)
			if err != nil {
				HandleError(errFile, err, metric.Text)
				continue
			} 
			e.saveJson(errFile, "netconfig", netconfig)

		} else if metric.Id == "networks" {
			networks, err := groupService.GetNetworkInterfaceService().GetNetworkInterfaces(nil)
			if err != nil {
				HandleError(errFile, err, metric.Text)
				continue
			} 
			e.saveJson(errFile, "networks", networks)
		}
	}
	log.Infof("run %s:elapse %s", e.Server, time.Since(startTime))

	return nil
}

