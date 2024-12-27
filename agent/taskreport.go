package agent

import (
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	// stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05"))
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))
	return []byte(stamp), nil
}

type ReportTaskTime struct {
	Start time.Time
	End   time.Time
	// Start JSONTime
	// End   JSONTime
}

type ReportTaskJob struct {
	Id    int
	Out   string
	Cmd   string
	Start time.Time
	End   time.Time
	// Start  JSONTime
	// End    JSONTime
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
		// cmdFix := strings.Replace(taskJob.CmdLine, " ", "", -1)

		report := ReportTaskJob{
			Id:  seq + 1,
			Out: taskJob.Job.Ofile,
			// Cmd:   cmdFix,
			Cmd:   taskJob.CmdLine,
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
