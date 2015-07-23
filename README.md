# moroccron
mesos + cron = moroccron

This is horrible and broken and probably wont even compile. Don't use it.

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
