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
$ ./moroccron -master $PLAYA_MESOS_HOST_HERE:5050  --logtostderr=true
```

# TODO

This doesnt actually register with mesos master:

```
-> $ godep go build && ./moroccron -master 10.141.141.10:5050  --logtostderr=true
2015/07/23 15:37:56 Creating executor information
2015/07/23 15:37:56 Creating scheduler
2015/07/23 15:37:56 Creating framework info
2015/07/23 15:37:56 Creating scheduler driver config
2015/07/23 15:37:56 Creating new scheduler driver from config
I0723 15:37:56.357082   50522 scheduler.go:237] Initializing mesos scheduler driver
2015/07/23 15:37:56 Starting scheduler driver
I0723 15:37:56.357427   50522 scheduler.go:645] Starting the scheduler driver...
I0723 15:37:56.357919   50522 http_transporter.go:290] http transport listening on 192.168.1.2:59197
I0723 15:37:57.361785   50522 scheduler.go:664] Mesos scheduler driver started with PID=scheduler(1)@192.168.1.2:59197
I0723 15:37:57.361879   50522 scheduler.go:814] Scheduler driver running.  Waiting to be stopped.
I0723 15:37:57.365989   50522 scheduler.go:277] New master master@127.0.1.1:5050 detected
I0723 15:37:57.366014   50522 scheduler.go:336] No credentials were provided. Attempting to register scheduler without authentication.


W0723 15:39:12.504812   50522 http_transporter.go:117] attempting to recover from error 'Post http://127.0.1.1:5050/master/mesos.internal.RegisterFrameworkMessage: dial tcp 127.0.1.1:5050: operation timed out', waiting before retry: 2s
:q
W0723 15:40:29.649221   50522 http_transporter.go:117] attempting to recover from error 'Post http://127.0.1.1:5050/master/mesos.internal.RegisterFrameworkMessage: dial tcp 127.0.1.1:5050: operation timed out', waiting before retry: 4s
W0723 15:41:48.838690   50522 http_transporter.go:117] attempting to recover from error 'Post http://127.0.1.1:5050/master/mesos.internal.RegisterFrameworkMessage: dial tcp 127.0.1.1:5050: operation timed out', waiting before retry: 8s
W0723 15:43:11.947128   50522 http_transporter.go:117] attempting to recover from error 'Post http://127.0.1.1:5050/master/mesos.internal.RegisterFrameworkMessage: dial tcp 127.0.1.1:5050: operation timed out', waiting before retry: 16s
I0723 15:44:43.519938   50522 scheduler.go:1102] Aborting driver, got error ' Failed to send message mesos.internal.RegisterFrameworkMessage: Post http://127.0.1.1:5050/master/mesos.internal.RegisterFrameworkMessage: dial tcp 127.0.1.1:5050: operation timed out '
I0723 15:44:43.519959   50522 scheduler.go:869] Aborting framework [nil]
I0723 15:44:43.520006   50522 scheduler.go:123] Scheduler received error: Failed to send message mesos.internal.RegisterFrameworkMessage: Post http://127.0.1.1:5050/master/mesos.internal.RegisterFrameworkMessage: dial tcp 127.0.1.1:5050: operation timed out
```
