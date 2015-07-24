package scheduler

import (
	"github.com/gogo/protobuf/proto"
	"strconv"

	log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
)

type Scheduler struct {
	executor      *mesos.ExecutorInfo
	tasksLaunched int
	tasksFinished int
	totalTasks    int
	JobsCh        chan string
}

func NewScheduler(exec *mesos.ExecutorInfo, ch chan string) *Scheduler {
	return &Scheduler{
		executor: exec,
		JobsCh:   ch,
	}
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

		//if we dont have any work to do, just driver.DeclineOffer(offerId *mesos.OfferID, filters *mesos.Filters)
		// see if we have any jobs waiting to run. for now, just use a channel full of jobs
		var data string
		select {
		case data, ok := <-sched.JobsCh:
			if ok {
				log.Infof("Got work %s\n", data)
			} else {
				//TODO should we abort?
				log.Infoln("Channel closed! FUCK why did this happen")
			}
		default:
			log.Infof("No pending work; declining offer %s: %+v", offer.Id, offer)
			driver.DeclineOffer(offer.Id, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
			continue
		}

		taskId := &mesos.TaskID{
			Value: proto.String(strconv.Itoa(sched.tasksLaunched)),
		}

		task := &mesos.TaskInfo{
			Name:      proto.String("moroccron-task-" + taskId.GetValue()),
			TaskId:    taskId,
			SlaveId:   offer.SlaveId,
			Executor:  sched.executor,
			Resources: []*mesos.Resource{
			//TODO stuff in constraints
			//util.NewScalarResource("cpus", sched.cpuPerTask),
			//util.NewScalarResource("mem", sched.memPerTask),
			},
			Data: []byte(data),
		}
		log.Infof("Prepared task: %s with offer %s for launch\n", task.GetName(), offer.Id.GetValue())

		var tasks []*mesos.TaskInfo = []*mesos.TaskInfo{task}
		//TODO i dont understand how you can launch multiple tasks for a single offer
		log.Infoln("Launching ", len(tasks), " tasks for offer", offer.Id.GetValue())
		driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
	}
}

func (sched *Scheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	log.Infoln("Status update: task", status.TaskId.GetValue(), " is in state ", status.State.Enum().String())

	if status.GetState() == mesos.TaskState_TASK_FINISHED {
		sched.tasksFinished++
		log.Infoln("%v of %v tasks finished.", sched.tasksFinished, sched.totalTasks)
	}

	if sched.tasksFinished >= sched.totalTasks {
		log.Infoln("Total tasks completed, stopping framework.")
		driver.Stop(false)
	}

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
