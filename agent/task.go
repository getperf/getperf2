package agent

// Task は Collecter クラスで渡されたコマンド実行セットを TaskJob
// クラスに渡して複数のコマンド実行をスケジュールします。実行モードが
// concurrent の場合は並列に、serial の場合はシーケンシャルにTaskJob
//  を実行します。
// 全ての TaskJob 終了後、各 TaskJob 構造体から実行結果を収集してレ
// ポートを作成します。

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Task struct {
	Collector   *Collector
	Datastore   *Datastore
	Pid         int      /**< Collector process id */
	Status      ProcMode /**< Process status */
	StatTimeout int

	StartTime time.Time /**< Start time(UTC) */
	EndTime   time.Time /**< End time(UTC) */

	// DateDir   string /**< Date directory */
	// TimeDir   string /**< Time directory */
	Odir      string /**< Output directory */
	ScriptDir string /**< Script directory */
	TaskJobs  []*TaskJob
}

func NewTask(collector *Collector, odir, scriptDir string) *Task {
	task := &Task{}
	task.Collector = collector
	task.StartTime = time.Now()
	task.StatTimeout = collector.StatTimeout
	task.Odir = odir
	task.ScriptDir = scriptDir
	for i, job := range collector.Jobs {
		taskJob := NewTaskJob(i, job, odir, scriptDir)
		taskJob.Timeout = collector.StatTimeout
		task.TaskJobs = append(task.TaskJobs, taskJob)
	}
	return task
}

func (task *Task) Run() error {
	ctx, cancel := MakeContext(task.StatTimeout)
	defer cancel()
	return task.RunWithContext(ctx)
}

func (task *Task) RunWithContext(ctx context.Context) error {
	var err error
	collector := task.Collector
	prev, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("run task %s", err)
	}
	defer os.Chdir(prev)
	os.Chdir(task.ScriptDir)
	log.Info("chdir ", task.ScriptDir)
	if _, err := os.Stat(task.Odir); !os.IsNotExist(err) {
		if err := os.RemoveAll(task.Odir); err != nil {
			return fmt.Errorf("run task %s", err)
		}
	}
	if err := os.MkdirAll(task.Odir, 0777); err != nil {
		return fmt.Errorf("run task %s", err)
	}
	log.Debug("RunWithContext ", collector.StatMode)
	switch collector.StatMode {
	case "serial":
		for _, taskJob := range task.TaskJobs {
			_, err = taskJob.RunWithContext(ctx)
		}
	case "concurrent":
		begin := make(chan interface{})
		var wg sync.WaitGroup
		for i, taskJob := range task.TaskJobs {
			wg.Add(1)
			// 何れも id2の taskJobが実行されてしまう問題あり
			// serial では正常動作、 go func の使い方調査 ※ v1.4で修正
			// go test -v -run  TestConcurrentTask
			go func(i int, job *TaskJob) {
				defer wg.Done()
				<-begin
				_, err = job.RunWithContext(ctx)
			}(i, taskJob)
		}
		close(begin)
		wg.Wait()
	}
	task.EndTime = time.Now()
	return err
}
