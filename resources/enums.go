package resources

type AccessLevel int64

const (
	Guest AccessLevel = iota + 1
	Worker
	Accountant
	Manager
	Admin
)

type WorkerStatus string

const (
	Busy      WorkerStatus = "busy"
	Available WorkerStatus = "available"
	Vacation  WorkerStatus = "vacation"
)
