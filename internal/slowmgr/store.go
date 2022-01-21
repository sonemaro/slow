package slowmgr

import (
	"errors"
	"sort"
	"sync"
)

const MinEntries = 10

var ErrMaxEntries = errors.New("maximum entries, cannot add more")

type Store interface {
	// Add adds a slow query to the store
	Add(s *Slow) error

	// Filter returns slow queries with pagination
	Filter(verb string, pageNo int, entries int) []*Slow

	// Sort sorts slow queries with pagination and order.
	Sort(pageNo, entries int) []*Slow
}

type DefaultStoreOptions struct {
	MaxEntries int
}

// StoreDefault is an in-memory implementation of Store
// Note: the performance is so bad. It's for demo purposes.
type StoreDefault struct {
	Options DefaultStoreOptions
	Data    []*Slow
	mu      sync.Mutex
}

func NewStoreDefault(o DefaultStoreOptions) *StoreDefault {
	if o.MaxEntries < MinEntries {
		o.MaxEntries = MinEntries
	}
	return &StoreDefault{
		Options: o,
		Data:    make([]*Slow, 0, o.MaxEntries),
	}
}

func (s *StoreDefault) Add(sq *Slow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.Data) >= s.Options.MaxEntries {
		return ErrMaxEntries
	}
	s.Data = append(s.Data, sq)
	return nil
}

func (s *StoreDefault) Filter(operation string, pageNo int, entries int) []*Slow {
	var ret []*Slow
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, q := range Paginate(s.Data, pageNo, entries) {
		if q.Operation == operation {
			ret = append(ret, q)
		}
	}

	return ret
}

func (s *StoreDefault) Sort(pageNo int, entries int) []*Slow {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := Paginate(s.Data, pageNo, entries)
	sort.Slice(p, func(i int, j int) bool {
		return p[i].Duration < p[j].Duration
	})

	return p
}

func Paginate(x []*Slow, page int, items int) []*Slow {
	start := (page - 1) * items
	stop := start + items

	if start > len(x) {
		return nil
	}
	if stop > len(x) {
		stop = len(x)
	}

	return x[start:stop]
}
