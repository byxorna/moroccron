package job

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
)

type JobPriority int32

const (
	LowPriority    JobPriority = -1
	NormalPriority JobPriority = 0
	HighPriority   JobPriority = 1
)

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

	// for scheduling. TODO: support cron syntax parsing, and ISO8601 recurring intervals?
	//NextRun time.Time     `json:"next_run"`
	//Every   time.Duration `json:"every"`

	scheduling_priority JobPriority `json:"priority"`
	// priority ranking for the job queue
	priority int `json:"-"`
	index    int `json:"-"`
}

func (j *Job) String() string {
	return j.Id
}

// recompute job priority based on last run, time, etc
func (j *Job) ComputePriority() int {
	//TODO FIXME how do you use const types to do math? j.scheduling_priority*10
	return j.priority
}
