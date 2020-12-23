package vmwareconf

type Metric struct {
	Id    string `toml:"id"`
	Level int    `toml:"level"`
	Text  string `toml:"text"`
}

func (metric *Metric) getObjectId() string {
	return metric.Text
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
	NewMetric(0, "summary", "summary"),
	NewMetric(0, "resourceConfig", "resourceConfig"),
	NewMetric(0, "guestHeartbeatStatus", "guestHeartbeatStatus"),
	NewMetric(0, "config", "config"),
	NewMetric(0, "guest", "guest"),
}
