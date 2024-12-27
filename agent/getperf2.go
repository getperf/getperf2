package agent

import (
	"context"
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/google/gops/agent"
	log "github.com/sirupsen/logrus"
)

const cmdName = "getperf"

func gops() error {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	return nil
}

// Run the getperf2
func Run(ctx context.Context, argv []string, stdout, stderr io.Writer) error {
	gops()
	var (
		c = flag.String("config", "", "It performs by the specified directory.")
		s = flag.String("statid", "", "Agent run the specified category once.")
		b = flag.Bool("background", true, "Agent run as background service.")
		f = flag.Bool("foreground", false, "Agent run as foreground service.")
	)
	flag.StringVar(c, "c", "", "")
	flag.StringVar(s, "s", "", "")
	flag.BoolVar(b, "b", true, "")
	flag.BoolVar(f, "f", false, "")
	flag.Usage = func() {
		PrintUsage()
	}

	flag.Parse()
	hostName, err := GetHostname()
	if err != nil {
		log.Fatal("get hostname ", err)
	}
	if *c == "" {
		home, err := GetParentAbsPath(os.Args[0], 2)
		if err != nil {
			log.Fatal("get getperf path ", err)
		}
		*c = filepath.Join(home, "getperf.ini")
	}
	configEnv := NewConfigEnvBase(hostName, cmdName, *c)
	home, err := GetParentAbsPath(*c, 1)
	if err != nil {
		log.Fatal("get home ", err)
	}
	config := NewConfig(home, configEnv)
	// 現行版の実行オプションと合わせるため、既定はバックグラウンドモード
	// で実行する。[-f|--foreground] オプション
	if *f == true {
		*b = false
	}
	if *b {
		if err := SetLog(home); err != nil {
			log.Fatal("set log ", err)
		}
	}
	log.Info("agent start ", VersionMessage())
	config.ParseConfigFile(config.ParameterFile)

	// Webサービス設定ファイル読込み
	config.ParseConfigFile(config.SvParamFile)

	if err := config.CheckConfig(); err != nil {
		log.Fatal("check parameter ", err)
	}
	if err := SetLogLevel(config.Schedule.LogLevel); err != nil {
		log.Fatal("set log level ", err)
	}
	if *s != "" {
		log.Info("run single schedule")
		// TODO : conding single scheduler
	} else {
		return config.RunWithContext(ctx)
	}
	return nil
}
