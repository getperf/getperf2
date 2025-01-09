package alletraconf

type Metric struct {
	Level int    `toml:"level"`
	Id    string `toml:"id"`
	Text  string `toml:"text"`
}

type Metrics struct {
	Metrics []*Metric
}

func NewMetric(level int, id string, text string) *Metric {
	metric := &Metric{
		Level: level,
		Id:    id,
		Text:  text,
	}
	return metric
}

var metrics = []*Metric{
	NewMetric(0, "arrays", "arrays"),
	NewMetric(0, "disks", "disks"),
	NewMetric(0, "networks", "networks"),
	NewMetric(0, "netconfig", "netconfig"),
}
