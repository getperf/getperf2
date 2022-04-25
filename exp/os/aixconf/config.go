package aixconf

import (
    "io"

    "github.com/getperf/getperf2/cfg"
    . "github.com/getperf/getperf2/common"
    . "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type AIX struct {
    Server     string `toml:"server"`
    Url        string `toml:"url"`
    User       string `toml:"user"`
    Password   string `toml:"password"`
    SshKeyPath string `toml:"ssh_key"`
    Insecure   bool   `toml:"insecure"`

    Metrics []*Metric `toml:"metrics"`

    Env       *cfg.RunEnv
    errFile   io.Writer
    datastore string
}

var sampleTemplateConfig = `
# AIX server inventory collector configuration
# Enter the information for AIX login account
# 
# example:
#
# url = "192.168.10.100"
# user = "test_user"
# password = "P@ssword"
# server = "aix7"

url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
insecure = true
server = "{{ .Server }}"

# The following parameters are optional


# Describe the additional command list. Added to the default command list for
# AIX inventory scenarios. The text parameter using escape codes such as
# '\"', '\\', See these example,
# 
# example:
# 
# [[metrics]]
# 
# id = "oracle_module"   # unique key
# level = 0    # command level [0:Default,1,2]
# type = "Cmd" # "Cmd":single command(Defalut), "Script":multi line commands
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
    return "aix"
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
