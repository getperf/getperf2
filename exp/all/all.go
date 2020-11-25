package all

import (
	_ "github.com/getperf/getperf2/exp/hw/hpiloconf"
	_ "github.com/getperf/getperf2/exp/hw/primergyconf"
	_ "github.com/getperf/getperf2/exp/monitor/zabbixconf"
	_ "github.com/getperf/getperf2/exp/os/linuxconf"
	_ "github.com/getperf/getperf2/exp/os/windowsconf"
	_ "github.com/getperf/getperf2/exp/vm/vmhostconf"
	_ "github.com/getperf/getperf2/exp/vm/vmwareconf"
)
