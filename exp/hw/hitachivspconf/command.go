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
	NewMetric(0, "overview", "",
		"/ConfigurationManager/v1/objects/storages"),
}
