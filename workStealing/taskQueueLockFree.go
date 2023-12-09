package workStealing

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// TaskNode contains particle number that needs to be worked on. Particle number is the index within the list of particles.
type TaskNode struct {
	next        unsafe.Pointer // points towards the bottom
	prev        unsafe.Pointer // points towards the top
	particleNum int64          // The particle we need to calculate for the task
}

// Create a new Task Node
func newTaskNode(particleNum int64) *TaskNode {
	node := &TaskNode{
		next:        nil,
		prev:        nil,
		particleNum: particleNum,
	}
	return node
}

/*
Unbounded lock free double ended task queue

Pointers from top -> bottom.
*/
type TaskQueue struct {
	bottom unsafe.Pointer
	top    unsafe.Pointer
}

// Initialize a new Task DeQueue
func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		bottom: nil,
		top:    nil,
	}
}

// Get the bottom element
func (q *TaskQueue) GetBottom() *TaskNode {
	return (*TaskNode)(atomic.LoadPointer(&q.bottom))
}

// Get the top element
func (q *TaskQueue) GetTop() *TaskNode {
	return (*TaskNode)(atomic.LoadPointer(&q.top))
}

// Print the task queue from top to bottom
func (q *TaskQueue) PrintTaskQueueFromTop() {
	fmt.Printf("Top -> ")
	cur := (*TaskNode)(atomic.LoadPointer(&q.top))
	bottom := (*TaskNode)(atomic.LoadPointer(&q.bottom))
	for {
		fmt.Printf("%d ", cur.particleNum)
		if cur == bottom {
			break
		}
		cur = (*TaskNode)(cur.next)
	}
	fmt.Printf("-> Bottom\n")
}

// Print the task queue from bottom to top
func (q *TaskQueue) PrintTaskQueueFromBottom() {
	fmt.Printf("Bottom <- ")
	cur := (*TaskNode)(atomic.LoadPointer(&q.bottom))
	top := (*TaskNode)(atomic.LoadPointer(&q.top))
	for {
		fmt.Printf("%d ", cur.particleNum)
		if cur == top {
			break
		}
		cur = (*TaskNode)(cur.prev)
	}
	fmt.Printf("<- Top\n")
}

/*
Push a task node to the bottom of the task queue

Cannot happen concurrently with the popBottom
*/
func (q *TaskQueue) PushBottom(parNum int64) {
	task := newTaskNode(parNum)

	for {
		// If popTop removes the bottom element, then we must reset references and try again

		bottom := (*TaskNode)(atomic.LoadPointer(&q.bottom))
		task.prev = unsafe.Pointer(bottom)
		if bottom != nil {
			next := (*TaskNode)(atomic.LoadPointer(&bottom.next))
			if atomic.CompareAndSwapPointer(&bottom.next, unsafe.Pointer(next), unsafe.Pointer(task)) {
				// Would fail with a concurrent pop of the bottom
				// If it is popped then the next pointer having been changed does not matter?
				if atomic.CompareAndSwapPointer(&q.bottom, unsafe.Pointer(bottom), unsafe.Pointer(task)) {
					break
				}
			}
		} else {
			// If the queue is empty
			if atomic.CompareAndSwapPointer(&q.bottom, unsafe.Pointer(bottom), unsafe.Pointer(task)) {
				q.top = unsafe.Pointer(task)
				break
			}
		}
	}
}

/*
Pop a task from the bottom of the queue. If fails, return -1

Will not occur concurrently with PushBottom
*/
func (q *TaskQueue) PopBottom() int64 {
	bottom := (*TaskNode)(atomic.LoadPointer(&q.bottom))
	top := (*TaskNode)(atomic.LoadPointer(&q.top))
	if bottom != nil && top != nil {
		// If queue not empty
		prev := (*TaskNode)(atomic.LoadPointer(&bottom.prev))
		if top == bottom {
			// One item in the list
			// Fails if popTop removes top == bottom first
			if atomic.CompareAndSwapPointer(&q.top, unsafe.Pointer(top), nil) {
				// Once in here popTop cannot pop since either top is changed from what it was or nil
				q.bottom = nil
				return bottom.particleNum
			}
		} else {
			if atomic.CompareAndSwapPointer(&q.bottom, unsafe.Pointer(bottom), unsafe.Pointer(prev)) {
				// If the top is equal to the popped bottom at this point the list must be empty
				// This works for popBottom Becasue there are no concurrent popBottoms that will take the bottom and try to pop it
				if atomic.CompareAndSwapPointer(&q.top, unsafe.Pointer(bottom), nil) {
					atomic.CompareAndSwapPointer(&q.bottom, unsafe.Pointer(prev), nil)
				}
				return bottom.particleNum
			}
		}
	}

	return -1

}

/*
Pop a task from the top of the queue. Used for stealing work. If fails, return -1.
*/
func (q *TaskQueue) PopTop() int64 {
	top := (*TaskNode)(atomic.LoadPointer(&q.top))
	bottom := (*TaskNode)(atomic.LoadPointer(&q.bottom))
	if top != nil && bottom != nil {
		// Not empty
		next := (*TaskNode)(atomic.LoadPointer(&top.next))
		if top == bottom {
			if atomic.CompareAndSwapPointer(&q.top, unsafe.Pointer(top), nil) {
				// Will fail if concurrent popTop or popBottom pops the top
				// If success no adds so other popTop will fail
				// If a parallel popBottom thinks that it is not at the last element it will pop the bottom before changing
				// the bottom pointer so if the bottom fails then all fails

				if atomic.CompareAndSwapPointer(&q.bottom, unsafe.Pointer(bottom), nil) {
					// When this fails it must be that a concurrent popBottom that started when
					// the queue has more elements popped the bottom off
					// Then we can fail and give it to the popBottom
					return top.particleNum
				}
			}
		} else {
			if atomic.CompareAndSwapPointer(&q.top, unsafe.Pointer(top), unsafe.Pointer(next)) {
				// pop top loads variables here so top == bottom
				// Pop Bottom right after above in a 2 element list
				// Concurrent popTop2 takes the top at this point as the new top which is a popped bottom
				// Before top is changed to nil

				// If the queue is empty after the pop
				// Can happen with concurrent popTop and Bottom
				if atomic.CompareAndSwapPointer(&q.bottom, unsafe.Pointer(top), nil) {
					atomic.CompareAndSwapPointer(&q.top, unsafe.Pointer(next), nil)
				}

				return top.particleNum
			}
		}
	}

	return -1
}
