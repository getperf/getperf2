package zabbixconf

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/pkg/errors"
	"github.com/rday/zabbix"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	// . "github.com/getperf/getperf2/common"
)

const (
	Insecure = true // 自己証明書の許可 `true`
)

func (e *Zabbix) initHostIds() {
	e.HostIds = make(map[string]string, 0)
	for _, host := range append(e.Servers, e.Server) {
		if host != "" {
			e.HostIds[host] = ""
		}
	}
}

func (e *Zabbix) saveJson(host, outfile string, value []byte) error {
	outPath := filepath.Join(e.datastore, host, outfile)
	if err := ioutil.WriteFile(outPath, value, 0666); err != nil {
		return HandleError(e.errFile, err, "save json")
	}
	return nil
}

func (e *Zabbix) GetHosts() error {
	hosts := []string{}
	for host := range e.HostIds {
		hosts = append(hosts, host)
	}
	params := map[string]interface{}{
		"output":                "extend",
		"selectInterfaces":      "extend",
		"selectGroups":          "extend",
		"selectMacros":          "extend",
		"selectParentTemplates": "extend",
		"filter": map[string]interface{}{
			"host": hosts,
		},
	}
	response, err := e.session.ZabbixRequest("host.get", params)
	if err != nil {
		return HandleError(e.errFile, err, "host.get request")
	}
	if response.Error.Code != 0 {
		return HandleError(e.errFile, &response.Error, "host.get request")
	}
	bytes, err := json.Marshal(response.Result)
	if err != nil {
		return HandleError(e.errFile, err, "decode host.get json result")
	}
	for _, host := range hosts {
		query := fmt.Sprintf("#(host == \"%s\")", host)
		value := gjson.Get(string(bytes), query).String()
		e.HostIds[host] = gjson.Get(value, "hostid").String()
		e.saveJson(host, "hosts", []byte(value))
	}
	return nil
}

func (e *Zabbix) GetLogItems() error {
	for host, hostId := range e.HostIds {
		if hostId == "" {
			continue
		}
		params := map[string]interface{}{
			"output": "extend",
			"filter": map[string]interface{}{
				"value_type": "2",
			},
			"hostids": hostId,
		}
		response, err := e.session.ZabbixRequest("Item.get", params)
		if err != nil {
			return HandleError(e.errFile, err, "Item.get request")
		}
		if response.Error.Code != 0 {
			return HandleError(e.errFile, &response.Error, "Item.get request")
		}
		bytes, err := json.Marshal(response.Result)
		if err != nil {
			return HandleError(e.errFile, err, "decode Item.get json result")
		}
		e.saveJson(host, "logItems", bytes)
	}
	return nil
}

func (e *Zabbix) GetTriggers() error {
	for host, hostId := range e.HostIds {
		if hostId == "" {
			continue
		}
		params := map[string]interface{}{
			"output":            "extend",
			"selectFunctions":   "extend",
			"expandData":        true,
			"expandDescription": true,
			"expandExpression":  true,
			"hostids":           hostId,
		}
		response, err := e.session.ZabbixRequest("Trigger.get", params)
		if err != nil {
			return HandleError(e.errFile, err, "Trigger.get request")
		}
		if response.Error.Code != 0 {
			return HandleError(e.errFile, &response.Error, "Trigger.get request")
		}
		bytes, err := json.Marshal(response.Result)
		if err != nil {
			return HandleError(e.errFile, err, "decode Trigger.get json result")
		}
		e.saveJson(host, "triggers", bytes)
	}
	return nil
}

func (e *Zabbix) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()
	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare error log")
	}
	defer errFile.Close()
	e.errFile = errFile
	e.HostIds = make(map[string]string, 0)
	e.datastore = env.Datastore

	e.initHostIds()
	for host := range e.HostIds {
		datastore := filepath.Join(e.datastore, host)
		if err := RemoveAndCreateDir(datastore); err != nil {
			return HandleError(errFile, err, "create log directory")
		}
	}
	url := fmt.Sprintf("%s/api_jsonrpc.php", e.Url)
	e.session, err = zabbix.NewAPI(url, e.User, e.Password)
	if err != nil {
		return HandleError(errFile, err, "prepare zabbix session")
	}
	_, err = e.session.Login()
	if err != nil {
		return HandleError(errFile, err, "login zabbix")
	}
	e.GetHosts()
	e.GetLogItems()
	e.GetTriggers()

	log.Infof("retrieve %s, elapse %s", e.HostIds, time.Since(startTime))

	return nil
}
