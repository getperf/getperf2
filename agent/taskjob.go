package agent

// TaskJob は Job クラスで渡されたコマンド実行情報を Process クラスに
// 渡してコマンドを実行します。定期実行の場合は、指定時間、回数で
// ループ実行します。リダイレクトモードの有無を判別し、非リダイレク
// トモードの場合はワークファイルに出力します。
// コマンド終了後、プロセスID, 開始/終了時刻, 終了コード, 標準出力,
// 標準エラーを構造体にセットします。標準出力, 標準エラーは指定した
// 行数分の先頭行を取得してバッファを返します。

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type TaskJob struct {
	Job       *Job
	Seq       int
	Timeout   int
	MaxRow    int
	Odir      string
	ScriptDir string
	CmdLine   string

	LoopCount   int
	CommandInfo *CommandInfo
	Status      ProcMode
	StartTime   time.Time
	EndTime     time.Time
}

func NewTaskJob(seq int, job *Job, odir, scriptDir string) *TaskJob {
	var taskJob TaskJob
	taskJob.Seq = seq
	taskJob.Job = job
	taskJob.Odir = odir
	taskJob.ScriptDir = scriptDir
	taskJob.MaxRow = DEFAULT_MAX_ERROR_LOG
	taskJob.Timeout = DEFAULT_COMMAND_TIMEOUT

	cmdLine := job.Cmd
	// Windows のコマンドで script.bat > \"log.txt\" の様に "の前の\を削除
	cmdLine = strings.Replace(cmdLine, "\"", "", -1)
	cmdLine = strings.Replace(cmdLine, "_odir_", taskJob.Odir, -1)
	cmdLine = strings.Replace(cmdLine, "_script_", taskJob.ScriptDir, -1)
	log.Debug("new task ", cmdLine)
	taskJob.CmdLine = cmdLine

	return &taskJob
}

func (taskJob *TaskJob) OutPath(ofile string) string {
	return filepath.Join(taskJob.Odir, ofile)
}

func (taskJob *TaskJob) Run() (*CommandInfo, error) {
	ctx, cancel := MakeContext(taskJob.Timeout)
	defer cancel()
	return taskJob.RunWithContext(ctx)
}

func (taskJob *TaskJob) RunWithContext(ctx context.Context) (*CommandInfo, error) {
	job := taskJob.Job
	commandInfo := &CommandInfo{
		CmdLine: taskJob.CmdLine,
		Timeout: taskJob.Timeout,
	}
	taskJob.CommandInfo = commandInfo

	var err error
	taskJob.StartTime = time.Now()
	if job.Ofile != "" {
		commandInfo.OutPath = taskJob.OutPath(job.Ofile)
		if job.Cycle > 0 && job.Step > 0 {
			err = commandInfo.PeriodicLoopCommand(ctx, job.Cycle, job.Step)
		} else {
			err = commandInfo.ExecCommandRedirectWithContext(ctx)
		}
	} else {
		err = commandInfo.ExecCommandNoRedirectWithContext(ctx)
	}
	taskJob.EndTime = time.Now()
	elapse := taskJob.EndTime.Sub(taskJob.StartTime)
	log.Infof("end [%d] %s [%s]", commandInfo.Pid, commandInfo.CmdLine, elapse)
	return commandInfo, err
}
