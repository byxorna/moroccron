package job

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
	"time"
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
	Arguments   []string           `json:"arguments,omitempty"`
	Environment *mesos.Environment `json:"environment,omitempty"`

	// for scheduling. TODO: support cron syntax parsing, and ISO8601 recurring intervals?
	NextRun time.Time     `json:"next_run"`
	Every   time.Duration `json:"every"`

	// for the job queue
	priority int
	index    int
}

func (j *Job) String() string {
	return j.Id
}
