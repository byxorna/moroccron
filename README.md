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

# Devving in vagrant

Vagrant is used to bring up a full dev environment with go, godep, mesos, zookeeper, etc. The code for moroccron is in `~/code`, and can be built with `godep go build`.

```
$ vagrant up
$ vagrant ssh
$ cd code && godep go build && ./morrocron -logtostderr
```

You can access the mesos master at http://10.10.0.5:5050 in your host browser.


#TODO

* represent jobs in a model
  * constraints, image, args/command, resources
* Make http api
  * create job
  * query job/deps
  * delete job
* Make ticker figure out if there are jobs to do
  * ticker could make resource requests when there is work to be done instead of waiting for an offer
* make scheduler smarter about scheduling more than 1 task per offer (keep track of resource limitations per job)
