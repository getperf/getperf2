package agent

import (
	"context"
	"flag"
	"io"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type program struct{}

func RunAdmin(ctx context.Context, argv []string, stdout, stderr io.Writer) error {
	var (
		flags   = flag.NewFlagSet("admin", flag.ExitOnError)
		fConfig = flags.String("config", "", "It performs by the specified directory.")
		fKey    = flags.String("key", "", "Site key.")
		fPass   = flags.String("pass", "", "Site administrator user password.")
		fUrl    = flags.String("url", "", "URL of Getperf admin web service.")
	)
	flags.StringVar(fConfig, "c", "", "")
	flags.StringVar(fKey, "k", "", "")
	flags.StringVar(fPass, "p", "", "")
	flags.StringVar(fUrl, "u", "", "")
	flags.Usage = func() {
		PrintUsage()
	}
	InitCommandMessages()
	// ./getperf2 [start|setup|..] -c aaa -k bbb ... の3列目から実行オプションを解析する
	flags.Parse(argv[2:])
	hostName, err := GetHostname()
	if err != nil {
		log.Fatal("get hostname ", err)
	}
	if *fConfig == "" {
		home, err := GetParentAbsPath(argv[0], 2)
		if err != nil {
			log.Fatal("get getperf path ", err)
		}
		*fConfig = filepath.Join(home, "getperf.ini")
	}
	configEnv := NewConfigEnvBase(hostName, cmdName, *fConfig)
	home, err := GetParentAbsPath(*fConfig, 1)
	if err != nil {
		log.Fatal("get home ", err)
	}
	config := NewConfig(home, configEnv)

	if config == nil {
		log.Fatal("initialize config error")
	}
	// if len(flags.Args()) != 1 {
	// 	flags.Usage()
	// 	return errors.Errorf("sub command [start|stop|setup] not found")
	// }
	command := argv[1]
	if runtime.GOOS != "windows" {
		switch command {
		case "start":
			log.Info("run start ", VersionMessage())
			config.RunDaemon()
		case "stop":
			err = config.StopDaemon()
		case "setup":
			err = config.RunSetup(*fKey, *fPass, *fUrl)
		default:
			flags.Usage()
			return errors.Errorf("unkown sub command %v", command)
		}
	} else {
		switch command {
		case "start":
			log.Info("run start ", VersionMessage())
			config.WindowsServiceControl(command)
		case "stop":
			err = config.WindowsServiceControl(command)
		case "install":
			err = config.WindowsServiceControl(command)
		case "remove":
			err = config.WindowsServiceControl("uninstall")
		case "setup":
			err = config.RunSetup(*fKey, *fPass, *fUrl)
		default:
			flags.Usage()
			return errors.Errorf("unkown sub command %v", command)
		}
	}
	return err
}
