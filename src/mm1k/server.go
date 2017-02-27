package mm1k

type server struct {
  id  int
  now float64
  gen Distribution
}

type SingleExpService server

// Return a scheduled completion time for a given start time
func (s *SingleExpService) Serve(startTime float64) (completionTime float64, sid int) {
  sid = s.id
  completionTime = startTime + s.gen.Get()
  s.now = completionTime
  return
}

func (s *SingleExpService) Next() float64 {
  return s.gen.Get()
}

// Create a single server with service time r
func MakeSingleExpService(r float64) (s *SingleExpService) {
  return &SingleExpService{1, 0, NewExpDistribution(r, 42)}
}








// A MinheapExpService is servers arranged in min-heap (according to their clock).
type MinheapExpService [](*server) // where heap is built upon

// Given start servicing time, return a scheduled departure time,
// and the corresponding server ID. Properties of the heap is maintained.
func (h MinheapExpService) Serve(now float64) (depTime float64, sid int) {
  sid = h[0].id
  depTime = now + h[0].gen.Get()
  h[0].now = depTime
  return
}

// Guarantees that the eariliest available time of the service module is
// returned.
func (h MinheapExpService) Next() float64 {
  return h[0].now
}

// Make a MinheapExpService of n servers, all of which are specified having service rate r.
func MakeMinheapExpService(r float64, seed int64) (h MinheapExpService) {
  h = make([]*server, 1)
  p := make([]server, 1) // pointer to the underlying array that stores servers
  p[0].id = 1
  p[0].gen = NewExpDistribution(r, seed)
  h[0] = &p[0]
  return
}
