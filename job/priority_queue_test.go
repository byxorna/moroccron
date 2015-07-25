package job

import (
	"container/heap"
	"github.com/robfig/cron"
	"testing"
	"time"
)

func newJob(id string, pri int) *Job {
	sch, _ := cron.Parse("* * * * *")
	// subtract priority from fixed last run, so higher priority jobs were last run further back in time
	t_last_run := time.Unix(1437856044-int64(pri), 0)
	return &Job{
		Id:                  id,
		t_schedule:          sch,
		t_last_run:          &t_last_run,
		scheduling_priority: NormalPriority,
	}
}

func TestPriorityQueueInit(t *testing.T) {
	items := map[string]int{
		"c": 5, "d": 3, "e": 0, "b": 15,
	}

	pq := JobQueue{}
	heap.Init(&pq)
	for id, pri := range items {
		heap.Push(&pq, newJob(id, pri))
	}

	//push in a new job with a high and low priority
	heap.Push(&pq, newJob("a", 99))
	heap.Push(&pq, newJob("z", -19))

	// make sure the ordering is correct
	target_order := []string{"a", "b", "c", "d", "e", "z"}
	i := 0
	for pq.Len() > 0 {
		j := heap.Pop(&pq).(*Job)
		t.Logf("Found job:%s pri:%f", j.Id, j.priority)
		if j.Id != target_order[i] {
			t.Errorf("Job id %s expected, but found %s at position %d priority %f", target_order[i], j.Id, i, j.priority)
		}
		i++
	}

}

func TestPriorityQueueUpdate(t *testing.T) {
	items := map[string]int{
		"c": 5, "d": 3, "e": 0, "b": 15,
	}

	pq := JobQueue{}
	heap.Init(&pq)
	for id, pri := range items {
		heap.Push(&pq, newJob(id, pri))
	}

	j := newJob("z", -19)
	heap.Push(&pq, j)

	// make j last run further back in time to act like we increased its priority
	newt := time.Unix(1437856044-int64(1000), 0)
	j.t_last_run = &newt
	pq.Update(j)

	top := heap.Pop(&pq).(*Job)
	t.Logf("Found job:%s pri:%f", top.Id, top.priority)
	if top.Id != "z" {
		t.Errorf("z was not the top job of the queue; found %s instead", top.Id)
	}

}

//TODO FUCK this is broken. wtf happened
func TestPriorityQueuePeek(t *testing.T) {
	items := map[string]int{
		"c": 5, "d": 3, "e": 0, "b": 15,
	}

	pq := make(JobQueue, len(items))
	i := 0
	for id, pri := range items {
		pq[i] = newJob(id, pri)
		i++
	}
	heap.Init(&pq)

	for i, v := range pq {
		t.Logf("%d %s:%f:%d", i, v.Id, v.priority, v.index)
	}

	top := heap.Pop(&pq).(*Job)
	t.Logf("Popped job:%s pri:%d", top.Id, top.priority)
	if top.Id != "b" {
		t.Errorf("b was not the top job of the queue; found %s instead", top.Id)
	}
	top = heap.Pop(&pq).(*Job)
	t.Logf("Popped job:%s pri:%d", top.Id, top.priority)
	if top.Id != "c" {
		t.Errorf("c was not the top job of the queue; found %s instead", top.Id)
	}

	top = pq.Peek().(*Job)
	t.Logf("Peeked at job:%s pri:%f", top.Id, top.priority)
	if top.Id != "b" {
		t.Errorf("b was not the top job of the queue; found %s instead", top.Id)
	}

	top = heap.Pop(&pq).(*Job)
	t.Logf("Popped job:%s pri:%d", top.Id, top.priority)
	if top.Id != "b" {
		t.Errorf("b was not the top job of the queue; found %s instead", top.Id)
	}

	top = pq.Peek().(*Job)
	t.Logf("Peeked at job:%s pri:%d", top.Id, top.priority)
	if top.Id != "c" {
		t.Errorf("c was not the top job of the queue; found %s instead", top.Id)
	}

}
