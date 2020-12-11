package hitachivspconf

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
	NewMetric(0, "storage", "",
		"/ConfigurationManager/v1/objects/storages/{id}"),
	NewMetric(0, "host-groups", "",
		"/ConfigurationManager/v1/objects/storages/{id}/host-groups"),
	NewMetric(0, "ports", "",
		"/ConfigurationManager/v1/objects/storages/{id}/ports"),
	NewMetric(0, "parity-groups", "",
		"/ConfigurationManager/v1/objects/storages/{id}/parity-groups"),
	NewMetric(0, "ldevs", "",
		"/ConfigurationManager/v1/objects/storages/{id}/ldevs?headLdevId=0&count=100"),
	NewMetric(0, "users", "",
		"/ConfigurationManager/v1/objects/storages/{id}/users"),
	NewMetric(0, "ambient", "",
		"/ConfigurationManager/v1/objects/storages/{id}/components/instance"),
	NewMetric(0, "snmp", "",
		"/ConfigurationManager/v1/objects/storages/{id}/snmp-settings/instance"),
}
