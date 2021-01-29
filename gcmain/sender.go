package gcmain

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/getperf/getperf2/agent"
	"github.com/getperf/getperf2/cfg"
	"github.com/getperf/getperf2/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Sender struct {
	AgentHome string
	CurrTime  time.Time
	Datastore string
	Platform  string
}

var timeout = 30

func convertInventoryPlatform(platform string) string {
	return platform + "Conf"
}

func zipForAgent(zipPath, outDir, newDir string) error {
	destinationFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	myZip := zip.NewWriter(destinationFile)
	defer myZip.Close()
	outDir = outDir + string(os.PathSeparator)
	err = filepath.Walk(outDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath := strings.TrimPrefix(filePath, outDir)
		targetPath := filepath.Join(newDir, relPath)
		// Windows でディレクトリパスはファイル名として展開されてしまうため"/"に置換
		// Windows\20200109\141630\ProcessorMemory.csv など
		targetPath = strings.Replace(targetPath, "\\", "/", -1)
		log.Debug("zip add ", targetPath)
		header, _ := zip.FileInfoHeader(info)
		header.Name = targetPath

		zipFile, err := myZip.CreateHeader(header)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func NewSender(platform string, env *cfg.RunEnv) *Sender {
	return &Sender{
		AgentHome: env.AgentHome,
		CurrTime:  env.CurrTime,
		Datastore: env.Datastore,
		Platform:  convertInventoryPlatform(platform),
	}
}

func (c *Sender) NewAgentConfig(hostName string) *agent.Config {
	configFile := filepath.Join(c.AgentHome, "getperf.ini")
	agentEnv := agent.NewConfigEnvBase(hostName, "gconf", configFile)
	return agent.NewConfig(c.AgentHome, agentEnv)
}

func (c *Sender) NewAgentDatastore(hostName string) *agent.Datastore {
	return agent.NewDatastoreBase(
		c.Datastore, hostName, c.Platform, c.CurrTime)
}

func (c *Sender) Run() error {
	// ctx, cancel := MakeContext(timeout)
	// defer cancel()
	hostName, err := common.GetHostname()
	if err != nil {
		log.Errorf("get hostname for initialize config %s", err)
		hostName = "UnkownHost"
	}
	log.Info("run sender")

	agentConfig := c.NewAgentConfig(hostName)
	agentDatastore := c.NewAgentDatastore(hostName)
	zipFile := agentDatastore.ZipFile(hostName)
	zipPath := filepath.Join(agentConfig.ArchiveDir, zipFile)
	relPath := agentDatastore.RelDir()
	if err := zipForAgent(zipPath, c.Datastore, relPath); err != nil {
		return errors.Wrap(err, "sending inventory zip")
	}
	if err := agentConfig.SendCollectorData(zipFile); err != nil {
		return errors.Wrap(err, "sending inventory zip")
	}
	return nil
}
