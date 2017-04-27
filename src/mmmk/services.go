package mmmk

import (
	"math"
	"mm1k"
	"sort"
)

type FIFO struct {
	arr []mm1k.Customer
}

func NewFIFO() (q *FIFO) {
	return &FIFO{make([]mm1k.Customer, 0)}
}

func (q *FIFO) Enqueue(customer mm1k.Customer) {
	q.arr = append(q.arr, customer)
}

func (q *FIFO) Dequeue() (customer mm1k.Customer) {
	customer, q.arr = q.arr[0], q.arr[1:]
	return
}

func (q FIFO) Empty() bool {
	return len(q.arr) == 0
}

type SJF struct {
	arr []mm1k.Customer
}

func NewSJF() (q *SJF) {
	return &SJF{make([]mm1k.Customer, 0)}
}

func (q *SJF) Enqueue(customer mm1k.Customer) {
	q.arr = append(q.arr, customer)
	sort.Sort(mm1k.ByService(q.arr))
	return
}

func (q *SJF) Dequeue() (customer mm1k.Customer) {
	customer, q.arr = q.arr[0], q.arr[1:]
	return
}

func (q SJF) Empty() bool {
	return len(q.arr) == 0
}

// MinService will send next job to the next available server
type MinService struct {
	current []mm1k.Customer
	gen     mm1k.Distribution
}

func (s *MinService) Get() float64 {
	return s.gen.Get()
}

func (s *MinService) Dequeue() (customer mm1k.Customer) {
	customer = s.current[s.NextAvailable()]
	s.current[s.NextAvailable()] = mm1k.Customer{Departure: math.Inf(1)}
	return customer
}

func (s *MinService) Add(customer mm1k.Customer) {
	// fmt.Printf("adding to %d\n", s.NextAvailable())
	s.current[s.NextAvailable()] = customer
}

func (s *MinService) NextAvailable() (i int) {
	// fmt.Printf("%v\n", s.current)
	i = 0
	next := s.current[i]
	for index, c := range s.current {
		if math.IsInf(c.Departure, 1) {
			return index
		}
		if c.Departure < next.Departure { // if server is not occupied, then Departure is infinity
			next = c
			i = index
		}
	}
	return
}

func (s *MinService) isFull() bool {
	for _, c := range s.current {
		if math.IsInf(c.Departure, 1) {
			return false
		}
	}
	return true
}

func (s *MinService) NextCompletion() mm1k.Customer {
	return s.current[s.NextAvailable()]
}

// Make a MinService of m servers
func MakeMinService(m int, serviceDistribution mm1k.Distribution) (s *MinService) {
	s = &MinService{make([]mm1k.Customer, m), serviceDistribution}
	for i := 0; i < m; i++ {
		s.current[i] = mm1k.Customer{Departure: math.Inf(1)}
	}
	return
}
