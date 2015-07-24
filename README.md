# moroccron
mesos + cron = moroccron

This is horrible and broken and probably wont even compile. Don't use it.

This is just a playground to see how easy it is to make a mesos framework.

## Features (wishful thinking)

* HA with master election
* API driven (register jobs, check job status)
* Triggers on success/failure
* Support arbitrary notification targets (email, pagerduty)
* Jobs have operational owners, but some may not want failure emails
* Jobs can have dependencies (a failure email could say what subtree of jobs are pending successful completion?)
* Not shitty web UI
* Resource isolation parameters
* Pluggable state backend? (docker/libkv? Do we want to use zk for state?)

## Build

```
$ godep go build
```
## Run
Start up mesos-playa (`vagrant up` in the playa-mesos directory) and get the IP for the mesos host: `vagrant hosts list`.
```
$ ./moroccron -master $PLAYA_MESOS_HOST_HERE:5050  -logtostderr
```


# Architecture

##Scheduler

##Executor

What runs on each slave. This is actually distributed by the framework as a standalone binary. This consumes resources, and will be delegated work (or finds work?).
