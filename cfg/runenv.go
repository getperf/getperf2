package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	. "github.com/getperf/getperf2/common"
	"github.com/pkg/errors"
)

const runEnvTemporaryName = "getconfigout"

type RunStatus string

const (
	INIT    = RunStatus("INIT")
	SUCCESS = RunStatus("SUCCESS")
	WARN    = RunStatus("WARN")
	ERROR   = RunStatus("ERROR")
)

type RunEnv struct {
	Level     int
	DryRun    bool
	Send      bool
	Datastore string
	LocalExec bool
	Timeout   int
	LogLevel  int
	CurrTime  time.Time
	AgentHome string
	Retrieve  *RetrieveConfig
	Filter    *FilterConfig
	Status    RunStatus
	Messages  string
}

func (e *RunEnv) Check() error {
	if e.Datastore == "" {
		datastore := filepath.Join(os.TempDir(), runEnvTemporaryName)
		if err := RemoveAndCreateDir(datastore); err != nil {
			return errors.Wrap(err, "init run env")
		}
		e.Datastore = datastore
	} else {
		if ok, _ := CheckDirectory(e.Datastore); !ok {
			return os.MkdirAll(e.Datastore, 0755)
		}
	}
	if e.Send {
		if e.AgentHome == "" {
			e.AgentHome = GetParentPath(GetBaseDir(), 1)
		}
		agentBackup := filepath.Join(e.AgentHome, "_bk")
		if exit, _ := CheckDirectory(agentBackup); exit != true {
			return fmt.Errorf("not found agent backup directory %s",
				agentBackup)
		}
	}
	return nil
}

func (e *RunEnv) OpenLog(fileName string) (*os.File, error) {
	outPath := filepath.Join(e.Datastore, fileName)
	return os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func (e *RunEnv) OpenServerLog(server, fileName string) (*os.File, error) {
	outPath := filepath.Join(e.Datastore, server, fileName)
	return os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func (e *RunEnv) NormalEnd(msg string) {
	e.Status = SUCCESS
	e.Messages = msg
}

func (e *RunEnv) WarningEnd(msg string) {
	e.Status = WARN
	e.Messages = msg
}

func (e *RunEnv) ErrorEnd(msg string) {
	e.Status = ERROR
	e.Messages = msg
}

func (e *RunEnv) Result() string {
	msg := string(e.Status)
	if e.Messages != "" {
		msg = msg + ":" + e.Messages
	}
	return msg
}
