package job

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/robfig/cron"
	"time"
)

type JobPriority string

const (
	LowPriority    JobPriority = "LOW"
	NormalPriority JobPriority = "NORMAL"
	HighPriority   JobPriority = "HIGH"
)

var JobPriorityValue = map[JobPriority]float64{
	"LOW":    0.5,
	"NORMAL": 1.0,
	"HIGH":   2.0,
}

type Job struct {
	Id string `json:"id"`

	Image   string          `json:"image,omitempty"` // pass into CommandInfo.CommandInfo_ContainerInfo.Image
	Volumes []*mesos.Volume `json:"volumes,omitempty"`
	Shell   *bool           `json:"shell,omitempty"`
	// command is "value" in mesos.CommandInfo
	// if shell == false, command is the binary, arguments are args
	Command *string `json:"command,omitempty"`
	// arguments are only read if shell == false
	Arguments   []string          `json:"arguments,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`

	// for scheduling. TODO: support ISO8601 recurring intervals?
	t_schedule   cron.Schedule `json:"-"`
	t_cron_entry string        `json:"cron_schedule"`
	t_last_run   *time.Time    //`json:"next_run"`

	scheduling_priority JobPriority `json:"priority"`
	// priority ranking for the job queue
	priority float64 `json:"-"`
	index    int     `json:"-"`
}

func New(image string, command string, arguments []string, shell bool, env map[string]string, cronspec string) (*Job, error) {
	//TODO FIXME
	j := Job{}
	return &j, nil
}

func (j *Job) String() string {
	return j.Id
}

// recompute job priority based on last run, time, etc
func (j *Job) ComputePriority() float64 {
	tnow := time.Now()
	next := j.NextScheduledRun()
	diff := tnow.Sub(next)
	if diff.Seconds() < 0 {
		// we havent surpassed our next scheduled run, so just give a low priority
		j.priority = diff.Seconds()
		return j.priority
	}
	j.priority = diff.Seconds() * JobPriorityValue[j.scheduling_priority]
	return j.priority
}

func (j *Job) NextScheduledRun() time.Time {
	if j.t_last_run == nil {
		// job not run yet, so lets compute when it will run next
		return j.t_schedule.Next(time.Now())
	}
	return j.t_schedule.Next(*j.t_last_run)
}
