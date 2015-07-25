package scheduler

import (
	"encoding/json"
	"fmt"
	. "github.com/byxorna/moroccron/job"
	"github.com/gogo/protobuf/proto"
	"time"

	log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
)

var (
	DOCKER_IMAGE_DEFAULT = "debian:latest"
)

type Scheduler struct {
	tasksLaunched int
	tasksFinished int
	totalTasks    int
	Jobs          *JobQueue
}

func NewScheduler() (*Scheduler, error) {
	s := Scheduler{}
	jobs, err := loadJobs()
	if err != nil {
		return nil, err
	}
	s.Jobs = jobs
	return &s, nil
}

func (sched *Scheduler) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Infoln("Scheduler Registered with Master ", masterInfo)
}

func (sched *Scheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Infoln("Scheduler Re-Registered with Master ", masterInfo)
}

func (sched *Scheduler) Disconnected(sched.SchedulerDriver) {
	log.Infoln("Scheduler Disconnected")
}

func (sched *Scheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	logOffers(offers)

	for _, offer := range offers {

		var (
			job Job
			ok  bool
		)
		select {
		case job, ok = <-sched.JobsCh:
			if ok {
				log.Infof("Got job %s\n", job.Id)
			} else {
				//TODO should we abort?
				log.Infoln("Channel closed! FUCK why did this happen")
			}
		default:
			log.Infof("No pending work; declining offer %s", offer.Id)
			driver.DeclineOffer(offer.Id, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
			continue
		}

		taskId := &mesos.TaskID{
			Value: proto.String(fmt.Sprintf("moroccron-task-%d", time.Now().Unix())),
		}

		job_json, err := json.Marshal(job)
		if err != nil {
			log.Errorf("Unable to serialize job %s: %s\n", job.Id, err.Error())
			driver.DeclineOffer(offer.Id, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
			continue
		}

		task := &mesos.TaskInfo{
			//TODO make this the timestamp of invocation, so its unique (with the id of the job name)
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
			Command: &mesos.CommandInfo{
				Shell: proto.Bool(true),
				Value: proto.String("set -x ; /bin/date ; /bin/hostname ; cat /etc/debian_version ; sleep 20 ; echo " + job.String()),
				//Uris: CommandInfo_URI{}
				//Value: string binary
				//Arguments: []string args to value
			},
			Executor: nil,
			Resources: []*mesos.Resource{
				//TODO this is bad. We shouldnt just blindly use up all the offered resources, but... whatever
				util.NewScalarResource("cpus", getOfferCpu(offer)),
				util.NewScalarResource("mem", getOfferMem(offer)),
			},
			Data: job_json,
		}

		log.Infof("Prepared task: %s with offer %s for launch\n", task.GetName(), offer.Id.GetValue())

		var tasks []*mesos.TaskInfo = []*mesos.TaskInfo{task}
		//TODO i dont understand how you can launch multiple tasks for a single offer. Is the up to the framework to slice resources per task?
		log.Infoln("Launching ", len(tasks), " tasks for offer", offer.Id.GetValue())
		driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
		sched.tasksLaunched++
	}
}

func (sched *Scheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	log.Infoln("Status update: task", status.TaskId.GetValue(), " is in state ", status.State.Enum().String())

	if status.GetState() == mesos.TaskState_TASK_FINISHED {
		sched.tasksFinished++
		log.Infoln("%v of %v tasks finished.", sched.tasksFinished, sched.totalTasks)
	}

	/*
		  //never shut down framework!
			if sched.tasksFinished >= sched.totalTasks {
				log.Infoln("Total tasks completed, stopping framework.")
				driver.Stop(false)
			}
	*/

	/*
		if status.GetState() == mesos.TaskState_TASK_LOST ||
			status.GetState() == mesos.TaskState_TASK_KILLED ||
			status.GetState() == mesos.TaskState_TASK_FAILED {
			log.Infoln(
				"Aborting because task", status.TaskId.GetValue(),
				"is in unexpected state", status.State.String(),
				"with message", status.GetMessage(),
			)
			driver.Abort()
		}
	*/
}

func (sched *Scheduler) OfferRescinded(s sched.SchedulerDriver, id *mesos.OfferID) {
	log.Infof("Offer '%v' rescinded.\n", *id)
}

func (sched *Scheduler) FrameworkMessage(s sched.SchedulerDriver, exId *mesos.ExecutorID, slvId *mesos.SlaveID, msg string) {
	log.Infof("Received framework message from executor '%v' on slave '%v': %s.\n", *exId, *slvId, msg)
}

func (sched *Scheduler) SlaveLost(s sched.SchedulerDriver, id *mesos.SlaveID) {
	log.Infof("Slave '%v' lost.\n", *id)
}

func (sched *Scheduler) ExecutorLost(s sched.SchedulerDriver, exId *mesos.ExecutorID, slvId *mesos.SlaveID, i int) {
	log.Infof("Executor '%v' lost on slave '%v' with exit code: %v.\n", *exId, *slvId, i)
}

func (sched *Scheduler) Error(driver sched.SchedulerDriver, err string) {
	log.Infoln("Scheduler received error:", err)
}

func getExecutor(data string) *mesos.ExecutorInfo {
	return &mesos.ExecutorInfo{
		ExecutorId: util.NewExecutorID("default"),
		Name:       proto.String("Moroccron Executor"),
		Source:     proto.String("moroccron"),
		Container: &mesos.ContainerInfo{
			Type:     mesos.ContainerInfo_DOCKER.Enum(),
			Volumes:  nil,
			Hostname: nil,
			Docker: &mesos.ContainerInfo_DockerInfo{
				Image: &DOCKER_IMAGE_DEFAULT,
			},
		},
		Command: &mesos.CommandInfo{
			Shell: proto.Bool(true),
			Value: proto.String("set -x ; /bin/date ; /bin/hostname ; cat /etc/debian_version ; sleep 20 ; echo " + data),
			//Uris: CommandInfo_URI{}
			//Value: string binary
			//Arguments: []string args to value
		},
	}
}
