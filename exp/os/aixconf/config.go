package aixconf

import (
    "io"

    "github.com/getperf/getperf2/cfg"
    . "github.com/getperf/getperf2/common"
    . "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type AIX struct {
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
# AIX server inventory collector configuration
# When collecting the inventory of AIX platform, execute it locally.
# Therefore, no template setting is required

server = "{{ .Server }}"
local_exec = true

# The following parameters are optional

# Enter the information for login account of remote AIX server
# 
# example:
# 
# [[servers]]
# 
# server = "aix1"    # server name
# url = "192.168.10.1"  # server address, example: 192.168.0.1 , 192.168.0.1:22
# user = "test_user"
# password = "P@ssword"
# ssh_key = ""          # ssh private key path, ignore if it not set

# Describe the additional command list. Added to the default command list for
# AIX inventory scenarios. The text parameter using escape codes such as
# '\"', '\\', See these example,
# 
# example:
# 
# [[metrics]]
# 
# id = "oracle_module"   # unique key
# level = 0      # command level [0-2]
# type = "Cmd"   # "Cmd":single command, "Script":multi line commands
# text = '''
# ls /home/oracle/"
# '''

[[metrics]]

id = "oslevel"
name = "OS情報"
category = "OSリリース"
text = '''
oslevel -s
'''

[[metrics]]

id = "osname"
name = "OS名"
category = "OSリリース"

[[metrics]]

id = "prtconf.System Model"
name = "モデル"
category = "OSリリース"

[[metrics]]

id = "prtconf.Machine Serial Number"
name = "シリアル"
category = "OSリリース"

[[metrics]]

id = "prtconf.Platform Firmware level"
name = "ファームレアレベル"
category = "OS設定"

[[metrics]]

id = "prtconf.Firmware Version"
name = "ファームウェアバージョン"
category = "OS設定"

[[metrics]]

id = "prtconf.Console Login"
name = "コンソールログイン"
category = "OS設定"

[[metrics]]

id = "prtconf.Auto Restart"
name = "自動起動"
category = "OS設定"

[[metrics]]

id = "prtconf.Full Core"
name = "全CPUコアの有効化"
category = "OS設定"

[[metrics]]

id = "prtconf.NX Crypto Acceleration"
name = "NX Crypto Acceleration"
category = "OS設定"

[[metrics]]

id = "prtconf.Kernel Type"
name = "カーネルタイプ"
category = "OS設定"

[[metrics]]

id = "prtconf.LPAR Info"
name = "LPAR構成"
category = "OS設定"

[[metrics]]

id = "prtconf.Processor Type"
name = "プロセッサータイプ"
category = "CPU"

[[metrics]]

id = "prtconf.Processor Implementation Mode"
name = "プロセッサーモード"
category = "CPU"

[[metrics]]

id = "prtconf.Processor Version"
name = "プロセッサーバージョン"
category = "CPU"

[[metrics]]

id = "prtconf.Number Of Processors"
name = "CPU数"
category = "CPU"

[[metrics]]

id = "prtconf.Processor Clock Speed"
name = "CPUクロック"
category = "CPU"

[[metrics]]

id = "prtconf.CPU Type"
name = "CPUタイプ"
category = "CPU"

[[metrics]]

id = "prtconf.Memory Size"
name = "メモリサイズ"
category = "メモリ"

[[metrics]]

id = "prtconf"
name = "システム構成情報"
category = "OSリリース"
text = '''
LANG=c prtconf
'''

[[metrics]]

id = "network"
name = "ネットワーク情報"
category = "ネットワーク"
text = '''
ifconfig -a
'''
`

func (e *AIX) Label() string {
    return "AIX"
}

func (e *AIX) Config() string {
    return sampleTemplateConfig
}

func init() {
    AddExporter("aixconf", func() Exporter {
        host, _ := GetHostname()
        return &AIX{
            Server: host,
        }
    })
}
