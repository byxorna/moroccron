package job

import (
	"container/heap"
)

// A JobQueue implements heap.Interface and holds Jobs.
type JobQueue []*Job

func (pq JobQueue) Len() int { return len(pq) }

func (pq JobQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].ComputePriority() > pq[j].ComputePriority()
}

func (pq JobQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *JobQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Job)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *JobQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *JobQueue) Peek() interface{} {
	n := len(*pq)
	item := (*pq)[n-1]
	return item
}

// update recomputes the priority of a Job in the queue and fixes the PQ up.
func (pq *JobQueue) Update(j *Job) {
	j.ComputePriority()
	//item.priority = priority
	heap.Fix(pq, j.index)
}
