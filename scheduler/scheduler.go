package scheduler

import (
	//	. "github.com/byxorna/moroccron/job"
	"github.com/byxorna/moroccron/state"
	"github.com/gogo/protobuf/proto"

	log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
)

var (
	DOCKER_IMAGE_DEFAULT = "debian:latest"
)

type Scheduler struct {
	tasksLaunched int
	tasksFinished int
	totalTasks    int
	State         state.Storage
	//RunningJobs   map[string]*Job
}

func NewScheduler(backend *state.Storage) (*Scheduler, error) {
	s := Scheduler{
		State: *backend,
	}
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
	jobs, err := getLaunchableJobs()
	if err != nil {
		log.Errorf("Unable to get pending jobs! %s\n", err.Error())
		return
	}

	offersAndTasks, err := packJobsInOffers(jobs, offers)
	if err != nil {
		log.Errorf("Unable to pack jobs into offers! %s\n", err.Error())
		return
	}

	for _, ot := range offersAndTasks {
		if len(ot.Tasks) == 0 {
			log.Infof("Declining unused offer %s", ot.Offer.Id.GetValue())
			driver.DeclineOffer(ot.Offer.Id, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
			continue
		} else {
			log.Infof("Launching %d tasks for offer %s\n", len(ot.Tasks), ot.Offer.Id.GetValue())
			driver.LaunchTasks([]*mesos.OfferID{ot.Offer.Id}, ot.Tasks, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
			sched.tasksLaunched = sched.tasksLaunched + len(ot.Tasks)
		}
	}

}

func (sched *Scheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	log.Infoln("Status update: task", status.TaskId.GetValue(), " is in state ", status.State.Enum().String())

	if status.GetState() == mesos.TaskState_TASK_FINISHED {
		sched.tasksFinished++
		log.Infoln("%v of %v tasks finished.", sched.tasksFinished, sched.totalTasks)
	}

	//TODO if a job is finished, failed, error, lost, killed
	// figure out how this impacts dependent jobs and update job graph

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

/*
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
*/
