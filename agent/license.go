package agent

type License struct {
	Hostname string
	Expired  string
	Code     string
}

func NewLicense() *License {
	license := License{}
	return &license
}
