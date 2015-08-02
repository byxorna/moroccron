package memory

import (
	"github.com/byxorna/moroccron/state"
)

type Driver struct{}

func New() *state.Storage {
	return (&Driver{}).(*state.Storage)
}

func (d *Driver) GetJobs() []*job.Job {
}
func (d *Driver) AddJob(job job.Job) (bool, error) {
}
func (d *Driver) UpdateJob(job job.Job) (bool, error) {
}
func (d *Driver) RemoveJob(id string) (bool, error) {
}
