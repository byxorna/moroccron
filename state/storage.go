package state

import (
	"github.com/byxorna/moroccron/job"
)

// interface that should be implemented by any job storage backend
type Storage interface {
	New() *Storage
	GetJobs() []*job.Job
	AddJob(job.Job) (bool, error)
	UpdateJob(job.Job) (bool, error)
	RemoveJob(string) (bool, error)
}
