package windowsconf

import (
	"io"

	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type Windows struct {
	Server    string     `toml:"server"`
	LocalExec bool       `toml:"local_exec"`
	Servers   []*Server  `toml:"servers"`
	Commands  []*Command `toml:"commands"`

	errFile    io.Writer
	datastore  string
	ScriptPath string
}

type Server struct {
	Server   string `toml:"server"`
	Url      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Insecure bool   `toml:"insecure"`
}

var sampleTemplateConfig = `
# Windows inventory exporter settings
# When collecting the inventory of Windows platform, execute it locally.
# Therefore, no template setting is required

local_exec = true

## Enter the information for login account of remote linux server
# 
## example:
# 
# [[servers]]
# 
# server = "centos8"    # server name
# url = "192.168.10.1"  # server address, example: 192.168.0.1 , 192.168.0.1:22
# user = "test_user"
# password = "P@ssword"
# ssh_key = ""          # ssh private key path, ignore if it not set

{{if ne .Url "" }}
[[servers]]

server = "{{ .Server }}"
url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
insecure = true
{{end}}

## Describe the additional command list. Added to the default command list for
## Windows Inventory scenarios. The text parameter using escape codes such as
## '\"', '\\', See these example,
#
## example:
# 
# [[commands]]
#
# id = "echo"    # unique key
# type = "Cmd"   # Cmd : cmd.exe -c "..." , Cmdlet : PowerShell -Command {...}
# level = 0      # command level [0-2]
# text = "echo 1"
`

func (e *Windows) Label() string {
	return "Windows"
}

func (e *Windows) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("windowsconf", func() Exporter {
		host, _ := GetHostname()
		return &Windows{
			Server: host,
		}
	})
}
