package scheduler

import (
	"fmt"
	. "github.com/byxorna/moroccron/job"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	"time"
)

func loadJobs() (*JobQueue, error) {
	//TODO hook me up with some data store
	//or at least read in from some flat files
	//jobs := []Job{}
	return &JobQueue{}, nil
}

func getLaunchableJobs() ([]Job, error) {
	//TODO FIXME this should query for all jobs ready to launch
	return []Job{}, nil
}

// given a set of offers, and jobs to run, try to pack the jobs into the offers
// returning a list of packed jobs with their offers
func packJobsInOffers(jobs []Job, offers []*mesos.Offer) ([]OfferTasksPair, error) {
	//TODO
	return []OfferTasksPair{}, nil
}

func createTask(job *Job, offer *mesos.Offer) mesos.TaskInfo {
	taskId := &mesos.TaskID{
		Value: proto.String(fmt.Sprintf("moroccron-task-%d-%s", time.Now().Unix(), job.Id)),
	}

	command_info := job.CreateCommandInfo()
	task := mesos.TaskInfo{
		Name:    proto.String(taskId.GetValue()),
		TaskId:  taskId,
		SlaveId: offer.SlaveId,
		Container: &mesos.ContainerInfo{
			Type:     mesos.ContainerInfo_DOCKER.Enum(),
			Volumes:  nil,
			Hostname: nil,
			Docker: &mesos.ContainerInfo_DockerInfo{
				Image:   &DOCKER_IMAGE_DEFAULT,
				Network: mesos.ContainerInfo_DockerInfo_BRIDGE.Enum(),
			},
		},
		Command:  &command_info,
		Executor: nil,
		Resources: []*mesos.Resource{
			util.NewScalarResource("cpus", job.CpuResources),
			util.NewScalarResource("mem", job.MemResources),
		},
		//Data: job_json,
	}
	return task
}
