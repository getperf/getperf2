package agent

import (
	"fmt"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type ReportTaskTime struct {
	Start time.Time
	End   time.Time
}

type ReportTaskJob struct {
	Id     int
	Out    string
	Cmd    string
	Start  time.Time
	End    time.Time
	Pid    int
	Rc     int
	Error  string
	Status string
}

type ReportTask struct {
	Schedule ReportTaskTime
	Jobs     []ReportTaskJob
}

func (t *Task) MakeReportTaskJobs() []ReportTaskJob {
	reports := []ReportTaskJob{}

	for seq, taskJob := range t.TaskJobs {
		// コマンドラインに空白があると80桁で改行してしまうため、空白を除く
		// https://github.com/go-yaml/yaml/issues/348
		cmdFix := strings.Replace(taskJob.CmdLine, " ", "", -1)

		report := ReportTaskJob{
			Id:    seq + 1,
			Out:   taskJob.Job.Ofile,
			Cmd:   cmdFix,
			Start: taskJob.StartTime,
			End:   taskJob.EndTime,
		}
		cInfo := taskJob.CommandInfo
		if cInfo != nil {
			report.Pid = cInfo.Pid
			report.Rc = cInfo.ExitCode
			report.Error = cInfo.OutBuf
			report.Status = cInfo.Status
		}
		reports = append(reports, report)
	}
	return reports
}

func (t *Task) MakeReport() (string, error) {
	taskTime := ReportTaskTime{Start: t.StartTime, End: t.EndTime}
	taskJobs := t.MakeReportTaskJobs()
	report := ReportTask{Schedule: taskTime, Jobs: taskJobs}
	data, err := yaml.Marshal(report)
	if err != nil {
		return "", fmt.Errorf("make report %s", err)
	}
	return string(data), nil
}
