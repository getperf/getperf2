package agent

import "time"

type Collector struct {
	Id            int    /**< primary key(sequence) */
	StatName      string /**< Metric */
	Status        string
	StatEnable    bool /**< Enabled */
	Build         int  /**< Build version */
	StatStdoutLog bool /**< Standard output flag */
	StatInterval  int  /**< Interval(sec) */
	StatTimeout   int  /**< Timeout(sec) */
	NextTimestamp time.Time
	StatMode      string /**< Stataus mode */

	Jobs []*Job /**< First job */
}

func NewCollector(statName string) *Collector {
	collector := Collector{
		StatName:      statName,
		NextTimestamp: time.Now(),
	}
	return &collector
}

func (collector *Collector) AddJob(job *Job) {
	collector.Jobs = append(collector.Jobs, job)
}

// func (collector *Collector)FindOrCreateJob(cmd string) *Job {
//     if job, ok := collector.Jobs[cmd]; ok {
//         log.Info("DBG found : ", *job)
//         return job
//     } else {
//         id := len(collector.Jobs) + 1
//         newJob := NewJob(id, cmd)
//         collector.AddJob(newJob)
//         log.Info("DBG new : ", newJob)
//         return newJob
//     }
// }
