package job

import (
	"container/heap"
	"testing"
)

func TestPriorityQueueInit(t *testing.T) {
	items := map[string]int{
		"c": 5, "d": 3, "e": 0, "b": 15,
	}

	pq := make(JobQueue, len(items))
	i := 0
	for id, pri := range items {
		pq[i] = &Job{
			Id:       id,
			priority: pri,
			index:    i,
		}
		i++
	}

	heap.Init(&pq)

	//push in a new job with a high and low priority
	heap.Push(&pq, &Job{
		Id:       "a",
		priority: 99,
	})
	heap.Push(&pq, &Job{
		Id:       "z",
		priority: -19,
	})

	// make sure the ordering is correct
	target_order := []string{"a", "b", "c", "d", "e", "z"}
	i = 0
	for pq.Len() > 0 {
		j := heap.Pop(&pq).(*Job)
		t.Logf("Found job:%s pri:%d", j.Id, j.priority)
		if j.Id != target_order[i] {
			t.Errorf("Job id %s expected, but found %s at position %d priority %d", target_order[i], j.Id, i, j.priority)
		}
		i++
	}

}

func TestPriorityQueueUpdate(t *testing.T) {
	items := map[string]int{
		"c": 5, "d": 3, "e": 0, "b": 15,
	}

	pq := make(JobQueue, len(items))
	i := 0
	for id, pri := range items {
		pq[i] = &Job{
			Id:       id,
			priority: pri,
			index:    i,
		}
		i++
	}

	heap.Init(&pq)

	j := &Job{
		Id:       "z",
		priority: -19,
	}
	heap.Push(&pq, j)

	j.priority = 1000
	pq.Update(j)

	top := heap.Pop(&pq).(*Job)
	t.Logf("Found job:%s pri:%d", top.Id, top.priority)
	if top.Id != "z" {
		t.Errorf("z was not the top job of the queue; found %s instead", top.Id)
	}

}
