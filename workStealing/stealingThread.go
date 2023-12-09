package workStealing

import (
	"math/rand"
	"time"
)

/*
Contains data necessary for work stealing implementation
*/
type WorkStealingThread struct {
	ThreadID     int
	Queue        *TaskQueue // Own task queue
	Rand         *rand.Rand
	AllWorkDone  *bool // Is all of the work done across all threads?
	WorkLeft     bool  // Is there work left within this thread?
	TotalThreads int   // Total number of threads total
}

// Initialize the thread data
func CreateThreadData(threadID int, totalThreads int, globalDone *bool) *WorkStealingThread {
	return &WorkStealingThread{
		ThreadID:     threadID,
		TotalThreads: totalThreads,
		Queue:        NewTaskQueue(),
		AllWorkDone:  globalDone,
		Rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		WorkLeft:     false,
	}
}
