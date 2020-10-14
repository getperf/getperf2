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
	"github.com/iancoleman/strcase"
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
// 	"config", "summary", "capability", "datastore",
// 	"environmentBrowser", "guest", "guestHeartbeatStatus",
// 	"layoutEx", "network", "resourceConfig", "resourcePool",
// 	"snapshot", "storage", "runtime",
// }

func (e *VMWare) saveJson(ioErr io.Writer, outfile, query string) {
	value := gjson.Get(e.json, query).String()
	if value == "" {
		fmt.Fprintf(ioErr, "not found json query[%s:%s] : '%s'\n", e.vmName, outfile, query)
		return
	}
	outPath := filepath.Join(e.datastore, outfile)
	if err := ioutil.WriteFile(outPath, []byte(value), 0666); err != nil {
		fmt.Fprintf(ioErr, err.Error())
	}
}

func (e *VMWare) retrieveInventory(env *cfg.RunEnv, ioErr io.Writer, vm string) {
	for _, metric := range e.Metrics {
		objectId := metric.getObjectId()
		if metric.Level == -1 || metric.Level > env.Level || objectId == "" {
			continue
		}
		query := strcase.ToCamel(objectId)
		e.saveJson(ioErr, objectId, query)
	}
	// Default extraction
	e.saveJson(ioErr, "base_hardware",
		"Config.Hardware")
	e.saveJson(ioErr, "base_cpu_memory_resource",
		"{ResourceConfig.CpuAllocation,ResourceConfig.MemoryAllocation}")
	e.saveJson(ioErr, "base_datastore",
		"Storage.PerDatastoreUsage")
	e.saveJson(ioErr, "base_disk",
		"Guest.Disk")
	e.saveJson(ioErr, "base_boot",
		"Config.BootOptions")
	e.saveJson(ioErr, "base_ipstack",
		"Guest.IpStack")
	e.saveJson(ioErr, "base_net",
		"Guest.Net")
	e.saveJson(ioErr, "base_ext_config",
		"Config.ExtraConfig")
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
	log.Info("search vm : ", e.Server)
	refVms, err := finder.VirtualMachineList(ctx, e.Server)
	if err != nil {
		HandleError(errFile, err, "get local vm defined in 'server' parameter")
	}
	if len(e.Servers) > 0 {
		for _, addedServer := range e.Servers {
			log.Info("search added vm : ", addedServer)
			refAddedVm, err := finder.VirtualMachine(ctx, addedServer)
			if err == nil {
				refVms = append(refVms, refAddedVm)
			} else {
				HandleError(errFile, err, "get remote vm defined in 'servers' parameter")
			}
		}
	}

	// 参照リストに変換する
	var refs []types.ManagedObjectReference
	for _, refVm := range refVms {
		refs = append(refs, refVm.Reference())
	}

	// 指定したVMのプロパティの取得
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	pc := property.DefaultCollector(session.Client)
	var vms []mo.VirtualMachine
	// if len(e.Metrics) > 0 {
	// 	log.Info("add metrics : ", e.Metrics)
	// 	vmMetrics = append(vmMetrics, e.Metrics...)
	// }
	for _, metric := range e.Metrics {
		if metric.Level == -1 || metric.Level > env.Level {
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
		return HandleError(errFile, err, "get vm info")
	}

	// 仮想インスタンス情報を出力
	if len(e.Servers) == 0 && len(vms) != 1 {
		return HandleError(errFile, err, fmt.Sprintf("vm not found '%s'", e.Server))
	}
	for _, vm := range vms {
		e.vmName = vm.Config.Name
		bytes, err := json.Marshal(vm)
		log.Info("json size : ", len(bytes))
		if err != nil {
			return HandleError(errFile, err, "convert json from vm info")
		}
		e.datastore = filepath.Join(env.Datastore, e.vmName)
		if err := RemoveAndCreateDir(e.datastore); err != nil {
			return HandleError(errFile, err, "create log directory")
		}
		alldataPath := filepath.Join(e.datastore, "all")
		if err := ioutil.WriteFile(alldataPath, bytes, 0666); err != nil {
			return HandleError(errFile, err, "write vm all info json")
		}
		e.json = string(bytes)
		e.retrieveInventory(env, errFile, e.vmName)
	}
	log.Infof("retrieve vm : %d, elapse %s", len(vms), time.Since(startTime))

	return nil
}
