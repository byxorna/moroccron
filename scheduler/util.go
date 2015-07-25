package scheduler

import (
	. "github.com/byxorna/moroccron/job"
	log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
)

//TODO this is probably all garbage and should get thrown away. Just including in the first pass to understand mesos better
func getOfferScalar(offer *mesos.Offer, name string) float64 {
	resources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
		return res.GetName() == name
	})

	value := 0.0
	for _, res := range resources {
		value += res.GetScalar().GetValue()
	}

	return value
}

func getOfferCpu(offer *mesos.Offer) float64 {
	return getOfferScalar(offer, "cpus")
}

func getOfferMem(offer *mesos.Offer) float64 {
	return getOfferScalar(offer, "mem")
}

func logOffers(offers []*mesos.Offer) {
	for _, offer := range offers {
		log.Infof("Received Offer <%v> with cpus=%v mem=%v", offer.Id.GetValue(), getOfferCpu(offer), getOfferMem(offer))
	}
}

func loadJobs() (*JobQueue, error) {
	//TODO hook me up with some data store
	//or at least read in from some flat files
	//jobs := []Job{}
	return &JobQueue{}, nil
}
