package vmwareconf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
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

func parseUrl(uri string) (string, error) {
	if !strings.HasPrefix(uri, "http://") &&
		!strings.HasPrefix(uri, "https://") {
		uri = "https://" + uri + "/sdk"
		log.Info("convert url : ", uri)
	}
	_, err := url.Parse(uri)
	if err != nil {
		return uri, errors.Wrapf(err, "parse url %s", uri)
	}
	return uri, nil
}

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
	for _, metric := range metrics {
		objectId := metric.getObjectId()
		if metric.Level == -1 || metric.Level > env.Level || objectId == "" {
			continue
		}
		query := strcase.ToCamel(objectId)
		e.saveJson(ioErr, metric.Id, query)
	}
}

func (e *VMWare) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()

	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare windows inventory error")
	}
	defer errFile.Close()

	e.Url, err = parseUrl(e.Url)
	if err != nil {
		return HandleError(errFile, err, "prepare rest url")
	}
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
			// refAddedVms, err := finder.VirtualMachine(ctx, addedServer)
			refAddedVms, err := finder.VirtualMachineList(ctx, addedServer)
			if err == nil {
				refVms = append(refVms, refAddedVms[0])
			} else {
				return HandleError(errFile, err, "get remote vm defined in 'servers' parameter")
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
	metrics = append(metrics, e.Metrics...)
	for _, metric := range metrics {
		if metric.Level == -1 || metric.Level > env.Level {
			continue
		}
		objectId := metric.getObjectId()
		if objectId == "" {
			continue
		}
		log.Info("add metrics : ", metric.Id)
		vmMetrics = append(vmMetrics, objectId)
	}
	err = pc.Retrieve(ctx, refs, vmMetrics, &vms)
	if err != nil {
		return HandleError(errFile, err, "get vm info")
	}

	cmd := info{
		General:   true,
		Resources: true,
	}
	res := infoResult{
		VirtualMachines: vms,
		cmd:             &cmd,
	}
	if err = res.collectReferences(pc, ctx); err != nil {
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

		hostName := "N/A"
		if href := vm.Summary.Runtime.Host; href != nil {
			if name, ok := res.entities[*href]; ok {
				hostName = name
			}
		}
		hostPath := filepath.Join(e.datastore, "host.txt")
		if err := ioutil.WriteFile(hostPath, []byte(hostName), 0666); err != nil {
			return HandleError(errFile, err, "write vm host info")
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
