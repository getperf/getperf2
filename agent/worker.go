package agent

// "time"

type Job struct {
	Id int /**< Primary key(sequence) */
	// Pid    int    /**< Worker process id */
	// Status int    /**< Process status */
	Cmd   string /**< Execute command */
	Ofile string /**< Output file */
	Cycle int    /**< Interval(sec) */
	Step  int    /**< Execute count */
}

func NewJob(id int, cmd string) *Job {
	var job Job
	job.Id = id
	job.Cmd = cmd
	return &job
}
