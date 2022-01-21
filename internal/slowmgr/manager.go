package slowmgr

type Manager interface {
	// Save persists the slow query
	Save(slow *Slow) error

	// Filter returns slow queries with pagination
	Filter(operation string, pageNo int, items int) []*Slow

	// Sort sorts slow queries with pagination and order.
	Sort(pageNo, items int) []*Slow

	// Start starts the manager
	Start() error

	// Stop stops the manager
	Stop()
}
