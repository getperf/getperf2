package linuxconf

import (
	"testing"
)

func TestLinuxUri(t *testing.T) {
	var tests = []struct {
		uri  string
		ip   string
		port string
	}{
		{"ssh://192.168.10.1:22", "192.168.10.1", "22"},
		{"192.168.10.1:22", "192.168.10.1", "22"},
		{"192.168.10.1", "192.168.10.1", "22"},
		{"hoge:", "hoge", "22"},
	}
	for _, test := range tests {
		ip, port, err := parseSshUrl(test.uri)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%s:Hostname()\t=>\t%s\n", test.uri, ip)
		t.Logf("%s:Port()\t=>\t%s\n", test.uri, port)
		if test.ip != ip || test.port != port {
			t.Errorf("parseShhUrl(%q) != %v, %v", test.uri, ip, port)
		}
	}
}
