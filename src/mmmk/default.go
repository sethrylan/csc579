package mmmk

import (
	"mm1k"
	"sort"
)

type InfiniteQueue struct {
	qType int // 0 - FIFO; 1 - SJF
	arr   []float64
	back  int
}

func NewInfiniteQueue(queueType int) (q InfiniteQueue) {
	q.arr = make([]float64, 10000000)
	q.qType = queueType
	return
}

func (q InfiniteQueue) Enqueue(arriveAt, serverAvailableAt float64) (t1 float64, sid int) {
	switch q.qType {
	case 0: // FIFO
		if arriveAt < serverAvailableAt { // If arrival is before next available service time
			t1 = serverAvailableAt // then queue job for next service time
		} else { // otherwise
			t1 = arriveAt // queue the job for the arrival time
		}
		q.arr[q.back] = t1
		sid = q.back
		q.back++
	case 1: // SJF
	default:
		panic("queue type not implementated")
	}
	return
}

func (q InfiniteQueue) Next() float64 {
	return q.arr[q.back]
}

type server struct {
	id  int
	now float64
}

// ByNow implements sort.Interface for []*server
type ByNow []*server

func (a ByNow) Len() int           { return len(a) }
func (a ByNow) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNow) Less(i, j int) bool { return a[i].now < a[j].now }

// MinService will send next job to the next available server
type MinService struct {
	a   [](*server)
	gen mm1k.Distribution
}

func (s MinService) Get() float64 {
	return s.gen.Get()
}

// Return departure time and serverID
func (s MinService) Serve(now float64) (depTime float64, sid int) {
	sid = s.a[0].id
	depTime = now + s.gen.Get()
	s.a[0].now = depTime
	sort.Sort(ByNow(s.a))
	return
}

// Eariliest available time is returned
func (s MinService) Next() float64 {
	return s.a[0].now
}

// Make a MinService of m servers
func MakeMinService(m int, serviceDistribution mm1k.Distribution) (s MinService) {
	s.gen = serviceDistribution
	s.a = make([]*server, m)
	p := make([]server, m) // pointer to array with actual objects
	for i := 0; i < m; i++ {
		p[i].id = i
		s.a[i] = &p[i]
	}
	return
}
