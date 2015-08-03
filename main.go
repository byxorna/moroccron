package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gogo/protobuf/proto"

	//log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"

	. "github.com/byxorna/moroccron/scheduler"
	"github.com/byxorna/moroccron/state"
	"github.com/byxorna/moroccron/state/memory"
	"github.com/byxorna/moroccron/web"
)

const (
	VERSION = "0.0.0"
)

var (
	storageEngine = flag.String("storage", "memory", "Storage backend (memory, ...)")
	master        = flag.String("master", "127.0.0.1:5050", "Master address <ip:port>")
	webPort       = flag.Int("web-port", 8000, "Port to serve http on (default: 8000)")
)

func init() {
	flag.Parse()
}

func main() {
	// create a storage backend
	//TODO switch based on storageEngine
	log.Printf("Creating %s storage driver\n", storageEngine)
	mem := memory.New()
	backend := state.Storage(&mem)

	// create our scheduler
	log.Println("Creating scheduler")
	scheduler, err := NewScheduler(&backend)
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
	go func() {
		if stat, err := driver.Run(); err != nil {
			log.Fatalf("Framework stopped with status %s and error: %s\n", stat.String(), err.Error())
			os.Exit(2)
		}
	}()

	log.Printf("Bringing up web interface at :%d\n", *webPort)
	router := web.New()
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", *webPort), router)
		if err != nil {
			log.Fatalf("Error launching web interface: %s\n", err.Error())
			os.Exit(2)
		}
	}()

	// chill out, squire! our web interface and scheduler are running in goroutines
	select {}
}
