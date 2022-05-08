package microservice

type RunStatus string

const (
	RunStatusSuspending RunStatus = "suspending"
	RunStatusRunning    RunStatus = "running"
)

type Status struct {
	IsEnabled bool
	RunStatus
}
