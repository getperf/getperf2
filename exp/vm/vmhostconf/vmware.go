package vmwareconf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	Insecure = true // 自己証明書の許可 `true`
)

// VMWare Managed Object Description
// Reference : http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html

var vmMetrics = []string{}

// var vmMetrics = []string{
// 	"summary", "capability", "datastore",
// 	"configManager", "hardware", "licensableResource",
// 	"network",
// }

func (e *VMWare) saveJson(ioErr io.Writer, outfile, query string) {
	value := gjson.Get(e.json, query).String()
	if value == "" {
		fmt.Fprintf(ioErr, "not found json query[%s:%s] : '%s'\n", e.vmName, outfile, query)
		return
	}
	outPath := filepath.Join(e.datastore, outfile+".json")
	if err := ioutil.WriteFile(outPath, []byte(value), 0666); err != nil {
		fmt.Fprintf(ioErr, err.Error())
	}
}

func (e *VMWare) retrieveInventory(ioErr io.Writer, vm string) {
	e.saveJson(ioErr, "hardware", "Summary.Hardware")
	e.saveJson(ioErr, "config", "Summary.Config")
	e.saveJson(ioErr, "product", "Summary.Config.Product")
}

// func HandleError(w io.Writer, inErr error, message string) error {
// 	if inErr != nil {
// 		_, err := fmt.Fprintf(w, "%s : %s\n", message, inErr)
// 		if err != nil {
// 			log.Errorf("write log error : %s", err)
// 		}
// 		return errors.Wrap(inErr, message)
// 	} else {
// 		return nil
// 	}
// }

func (e *VMWare) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()

	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare windows inventory error")
	}
	defer errFile.Close()

	urlConfig, err := url.Parse(e.Url)
	if err != nil {
		return HandleError(errFile, err, "parse vcenter url")
	}
	urlConfig.User = url.UserPassword(e.User, e.Password)

	session, err := govmomi.NewClient(ctx, urlConfig, Insecure)
	if err != nil {
		return HandleError(errFile, err, "connect vcenter url")
	}
	finder := find.NewFinder(session.Client, true)

	dc, err := finder.DefaultDatacenter(ctx)
	if err != nil {
		return HandleError(errFile, err, "attach data center")
	}
	finder.SetDatacenter(dc)

	// データセンター内の仮装インスタンスリストを取得
	log.Info("search host : ", e.Server)
	// refVms, err := finder.HostSystemList(ctx, e.Server)
	refVms, err := finder.HostSystemList(ctx, "*")
	if err != nil {
		HandleError(errFile, err, "get local host defined in 'server' parameter")
	}
	if len(e.Servers) > 0 {
		for _, addedServer := range e.Servers {
			log.Info("search added host : ", addedServer)
			refAddedVm, err := finder.HostSystem(ctx, addedServer)
			if err == nil {
				refVms = append(refVms, refAddedVm)
			} else {
				HandleError(errFile, err, "get remote host defined in 'servers' parameter")
			}
		}
	}

	// 参照リストに変換する
	var refs []types.ManagedObjectReference
	for _, refVm := range refVms {
		refs = append(refs, refVm.Reference())
	}

	// 指定したhostのプロパティの取得
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	pc := property.DefaultCollector(session.Client)
	var vms []mo.HostSystem
	// if len(e.Metrics) > 0 {
	// 	log.Info("add metrics : ", e.Metrics)
	// 	vmMetrics = append(vmMetrics, e.Metrics...)
	// }
	for _, metric := range e.Metrics {
		if metric.Level > env.Level {
			continue
		}
		objectId := metric.getObjectId()
		if objectId == "" {
			continue
		}
		log.Info("add metrics : ", objectId)
		vmMetrics = append(vmMetrics, objectId)
	}
	err = pc.Retrieve(ctx, refs, vmMetrics, &vms)

	if err != nil {
		return HandleError(errFile, err, "get host info")
	}

	// 仮想インスタンス情報を出力
	if len(e.Servers) == 0 && len(vms) != 1 {
		return HandleError(errFile, err, fmt.Sprintf("host not found '%s'", e.Server))
	}
	for _, vm := range vms {
		e.vmName = vm.Summary.Config.Name
		bytes, err := json.Marshal(vm)
		log.Info("json size : ", len(bytes))
		if err != nil {
			return HandleError(errFile, err, "convert json from host info")
		}
		e.datastore = filepath.Join(env.Datastore, e.vmName)
		if err := RemoveAndCreateDir(e.datastore); err != nil {
			return HandleError(errFile, err, "create log directory")
		}
		alldataPath := filepath.Join(e.datastore, "all.json")
		if err := ioutil.WriteFile(alldataPath, bytes, 0666); err != nil {
			return HandleError(errFile, err, "write host all info json")
		}
		e.json = string(bytes)
		e.retrieveInventory(errFile, e.vmName)
	}
	log.Infof("retrieve host : %d, elapse %s", len(vms), time.Since(startTime))

	return nil
}
