package primergyconf

type Metric struct {
	Level int    `toml:"level"`
	Id    string `toml:"id"`
	Batch string `toml:"batch"`
	Text  string `toml:"text"`
}

type Metrics struct {
	Metrics []*Metric
}

func NewMetric(level int, id string, batch string, text string) *Metric {
	metric := &Metric{
		Level: level,
		Id:    id,
		Batch: batch,
		Text:  text,
	}
	return metric
}

var metrics = []*Metric{
	NewMetric(0, "overview", "",
		"/redfish/v1/"),
	NewMetric(0, "firmware", "",
		"/redfish/v1/Systems/0/Oem/ts_fujitsu/FirmwareInventory"),
	NewMetric(0, "nic", "",
		"/redfish/v1/Systems/0/Oem/ts_fujitsu/FirmwareInventory/NIC"),
	NewMetric(0, "ntp0", "",
		"/redfish/v1/Managers/iRMC/Oem/ts_fujitsu/iRMCConfiguration/Time/NtpServers/0"),
	NewMetric(0, "ntp1", "",
		"/redfish/v1/Managers/iRMC/Oem/ts_fujitsu/iRMCConfiguration/Time/NtpServers/1"),
	NewMetric(0, "network", "",
		"/redfish/v1/Managers/iRMC/EthernetInterfaces/0"),
	NewMetric(0, "disk",
		"/rest/v1/Oem/eLCM/ProfileManagement/RAIDAdapter",
		"/rest/v1/Oem/eLCM/ProfileManagement/get?PARAM_PATH=Server/HWConfigurationIrmc/Adapters/RAIDAdapter"),
	NewMetric(0, "snmp",
		"/rest/v1/Oem/eLCM/ProfileManagement/NetworkServices",
		"/rest/v1/Oem/eLCM/ProfileManagement/get?PARAM_PATH=Server/SystemConfig/IrmcConfig/NetworkServices"),
}
