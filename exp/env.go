package exp

// type JobStatus string

// const (
// 	JOB_INIT    = JobStatus("JOB_INIT")
// 	JOB_SUCCESS = JobStatus("JOB_SUCCESS")
// 	JOB_WARN    = JobStatus("JOB_WARN")
// 	JOB_ERROR   = JobStatus("JOB_ERROR")
// )

// type Env struct {
// 	Level     int
// 	DryRun    bool
// 	Datastore string
// 	LocalExec bool
// 	Timeout   int
// 	Status    JobStatus
// 	Messages  string
// }

// func (env *Env) NormalEnd(msg string) {
// 	env.Status = JOB_SUCCESS
// 	env.Messages = msg
// }

// func (env *Env) WarningEnd(msg string) {
// 	env.Status = JOB_WARN
// 	env.Messages = msg
// }

// func (env *Env) ErrorEnd(msg string) {
// 	env.Status = JOB_ERROR
// 	env.Messages = msg
// }

// func (env *Env) Result() string {
// 	msg := string(env.Status)
// 	if env.Messages != "" {
// 		msg = msg + ":" + env.Messages
// 	}
// 	return msg
// }
