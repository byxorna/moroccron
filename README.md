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

# Devving in vagrant

Vagrant is used to bring up a full dev environment with go, godep, mesos, zookeeper, etc. The code for moroccron is in `~/code`, and can be built with `godep go build`.

```
$ vagrant up
$ vagrant ssh
$ cd code && godep go build && ./morrocron -logtostderr
```

You can access the mesos master at http://10.10.0.5:5050 in your host browser.


#TODO

## General

* create job packing function that packs jobs into offers
  * make scheduler smarter about scheduling more than 1 task per offer (keep track of resource limitations per job)
* Track running jobs in scheduler
  * add jobs when launching
  * update/remove jobs when statusUpdate

* Make http api
  * create job
  * query job/deps
  * delete job
* represent jobs in a model
  * constraints, image, args/command, resources
* should we request resources when we have a job to do instead of waiting for an offer?

## Metrics

Record start time, completion time, and skew for each job id

## Stuff to implement for HA framework
From Tan:
A few things to improve reliability and facilitate recovery once you have the basic functionality working:
1. For recovery of the framework, most people will persist the FrameworkID returned by the master to the framework via the Registered callback in an external HA store like ZK.  The mesos master will only allow one framework with a given FrameworkID to be registered at any time (when a new one tries to register, it kicks off the old one registered with the same ID and the old one).  By re-registering with the previous FrameworkID, it is possible for the failed-over new framework instance to interact with and receive status updates from tasks that were started by a previous instance of the framework.

2. When a framework receives the Registered callback (or at any time, but this and maybe after the Reregistered callback is received are the only places you probably want to do this) you can call the driver's ReconcileTasks method with an empty slice of TaskStatus's as the argument, which will cause the master to iterate through the last known task updates for a task and it will send them to the framework.  This is useful so that a failed-over framework can discover which tasks are currently running.  The updates cached by the master for reconciliation purposes don't have the "data" field that may have been sent by the executor (the original update will be received with it by the framework if the framework is alive, but the master deletes this field for caching as it may contain a lot of data, and we don't want the master to OOM).  These reconciled task updates will be received asynchronously and trigger the same StatusUpdate callback that non-reconciled updates trigger.

3. It's important to enable a relatively high FailoverTimeout in your FrameworkInfo that you register using.  This is the number of seconds that the master will wait before transitioning your framework to the Completed state after your framework process becomes unavailable.  When your framework goes into this state, all tasks are killed.  For long-running frameworks, some people set this to the number of seconds in a week, to give you a lot of time to try to fix whatever problem may be preventing your framework from running.

4. It's a good idea to set the Checkpoint field of FrameworkInfo to true.  This causes the slave to persist state about the tasks in this framework to the local filesystem, so that when the slave is restarted due to a problem or upgrade the new process can recover the task info, executor info, and status updates.  This allows you to upgrade your slaves without losing the tasks running on those machines.
