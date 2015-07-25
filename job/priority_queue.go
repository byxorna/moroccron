package job

import (
	"container/heap"
)

// A JobQueue implements heap.Interface and holds Jobs.
type JobQueue []*Job

func (pq JobQueue) Len() int { return len(pq) }

func (pq JobQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
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

// update modifies the priority and value of a Job in the queue.
func (pq *JobQueue) update(item *Job, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}

/*
func main() {
  // Some items and their priorities.
  items := map[string]int{
    "banana": 3, "apple": 2, "pear": 4,
  }

  // Create a priority queue, put the items in it, and
  // establish the priority queue (heap) invariants.
  pq := make(JobQueue, len(items))
  i := 0
  for value, priority := range items {
    pq[i] = &Item{
      value:    value,
      priority: priority,
      index:    i,
    }
    i++
  }
  heap.Init(&pq)

  // Insert a new item and then modify its priority.
  item := &Item{
    value:    "orange",
    priority: 1,
  }
  heap.Push(&pq, item)
  pq.update(item, item.value, 5)

  // Take the items out; they arrive in decreasing priority order.
  for pq.Len() > 0 {
    item := heap.Pop(&pq).(*Item)
    fmt.Printf("%.2d:%s ", item.priority, item.value)
  }
}
*/
