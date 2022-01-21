package slowmgr

// Slow holds a slow query in form of a struct so,
// it would be easier to deal with data and filtering
type Slow struct {
	// Query is the slow query which has recently executed
	Query string

	// Verb sql operation. example: SELECT, UPDATE, etc...
	Operation string

	// Duration holds how many seconds has been spent for execution of this slow query
	Duration float64
}
