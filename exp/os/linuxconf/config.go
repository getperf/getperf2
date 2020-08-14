package linuxconf

import (
	"io"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type Linux struct {
	Server    string    `toml:"server"`
	LocalExec bool      `toml:"local_exec"`
	Servers   []*Server `toml:"servers"`
	Metrics   []*Metric `toml:"metrics"`

	Env          *cfg.RunEnv
	errFile      io.Writer
	remoteServer string
	datastore    string
}

type Server struct {
	Server     string `toml:"server"`
	Url        string `toml:"url"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	SshKeyPath string `toml:"ssh_key"`
	Insecure   bool   `toml:"insecure"`
}

var sampleTemplateConfig = `
# Linux server inventory collector configuration
# When collecting the inventory of Linux platform, execute it locally.
# Therefore, no template setting is required

server = "{{ .Server }}"
local_exec = true

# The following parameters are optional

# Enter the information for login account of remote linux server
# 
# example:
# 
# [[servers]]
# 
# server = "centos8"    # server name
# url = "192.168.10.1"  # server address, example: 192.168.0.1 , 192.168.0.1:22
# user = "test_user"
# password = "P@ssword"
# ssh_key = ""          # ssh private key path, ignore if it not set

# Describe the additional command list. Added to the default command list for
# Linux inventory scenarios. The text parameter using escape codes such as
# '\"', '\\', See these example,
# 
# example:
# 
# [[metrics]]
# 
# id = "oracle_module"   # unique key
# level = 0      # command level [0-2]
# text = "ls /home/oracle/"

[[metrics]]

id = "hostname"
name = "ホスト名"
category = "OSリリース"
level = 0
comment = "hostname -s　コマンドで、ホスト名を検索"
text = '''
hostname -s
'''
`

func (e *Linux) Label() string {
	return "Linux"
}

func (e *Linux) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("linuxconf", func() Exporter {
		host, _ := GetHostname()
		return &Linux{
			Server: host,
		}
	})
}
