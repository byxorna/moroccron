package main

import (
	"flag"
	"log"
	"os"

	"github.com/gogo/protobuf/proto"

	//log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"

	. "github.com/byxorna/moroccron/scheduler"
)

const (
	VERSION = "0.0.0"
)

var (
	master = flag.String("master", "127.0.0.1:5050", "Master address <ip:port>")
)

func init() {
	flag.Parse()
}

func main() {

	// create our scheduler
	log.Println("Creating scheduler")
	scheduler, err := NewScheduler()
	if err != nil {
		log.Fatalf("Unable to create scheduler: %s\n", err.Error())
		os.Exit(1)
	}
	log.Printf("Created scheduler %+v\n", scheduler)

	// Framework
	log.Println("Creating framework info")
	fwinfo := &mesos.FrameworkInfo{
		User: proto.String(""), // Mesos-go will fill in user.
		Name: proto.String("moroccron-" + VERSION),
	}
	log.Printf("Created fwinfo %+v\n", fwinfo)

	// Scheduler Driver
	log.Println("Creating scheduler driver config")
	config := sched.DriverConfig{
		Scheduler:  scheduler,
		Framework:  fwinfo,
		Master:     *master,
		Credential: (*mesos.Credential)(nil),
	}
	log.Printf("Created driver config %+v\n", config)

	log.Println("Creating new scheduler driver from config")
	driver, err := sched.NewMesosSchedulerDriver(config)

	if err != nil {
		log.Fatalf("Unable to create a SchedulerDriver: %v\n", err.Error())
		os.Exit(3)
	}
	log.Printf("Created scheduler driver %+v\n", driver)

	log.Println("Starting scheduler driver")
	if stat, err := driver.Run(); err != nil {
		log.Fatalf("Framework stopped with status %s and error: %s\n", stat.String(), err.Error())
		os.Exit(4)
	}
}
