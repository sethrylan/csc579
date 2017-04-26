package mmmk

import (
	"math"
	"mm1k"
	"sort"
)

type job struct {
	serveAt     float64
	serviceTime float64
}
type byServiceTime []job

func (a byServiceTime) Len() int           { return len(a) }
func (a byServiceTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byServiceTime) Less(i, j int) bool { return a[i].serviceTime < a[j].serviceTime }

type FIFO struct {
	arr []job
	// back int
}

func NewFIFO() (q *FIFO) {
	return &FIFO{make([]job, 0)}
}

func (q *FIFO) Enqueue(arriveAt, serviceTime, serverAvailableAt float64) (j job) {
	j.serviceTime = serviceTime
	if arriveAt < serverAvailableAt { // If arrival is before next available service time
		j.serveAt = serverAvailableAt // then queue job for next service time
	} else { // otherwise
		j.serveAt = arriveAt // queue the job for the arrival time
	}
	q.arr = append(q.arr, j)
	return
}

func (q *FIFO) Next() (j job) {
	j = q.arr[0]
	q.arr = q.arr[1:]
	return
}

type SJF struct {
	arr []job
}

func NewSJF() (q SJF) {
	q.arr = make([]job, 0)
	return
}

func (q SJF) Enqueue(arriveAt, serviceTime, serverAvailableAt float64) (j job) {
	j.serviceTime = serviceTime
	if arriveAt < serverAvailableAt { // If arrival is before next available service time
		j.serveAt = serverAvailableAt // then queue job for next service time
	} else { // otherwise
		j.serveAt = arriveAt // queue the job for the arrival time
	}
	q.arr = append(q.arr, j)
	sort.Sort(byServiceTime(q.arr))
	return
}

func (q SJF) Next() (j job) {
	j = q.arr[:1][0]
	q.arr = q.arr[1:]
	return
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
func (s MinService) Serve(now job) (depTime float64, sid int) {
	sid = s.a[0].id
	depTime = math.Max(s.a[0].now, now.serveAt) + now.serviceTime
	s.a[0].now = depTime
	sort.Sort(ByNow(s.a))
	return
}

// Eariliest available time is returned
func (s MinService) Next() job {
	return job{s.a[0].now, 0}
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
