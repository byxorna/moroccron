package memory

import (
	"github.com/byxorna/moroccron/job"
)

type Driver struct{}

func New() Driver {
	return Driver{}
}

func (d *Driver) GetJobs() []*job.Job {
	return []*job.Job{}
}
func (d *Driver) AddJob(job job.Job) (bool, error) {
	return false, nil
}
func (d *Driver) UpdateJob(job job.Job) (bool, error) {
	return false, nil
}
func (d *Driver) RemoveJob(id string) (bool, error) {
	return false, nil
}
