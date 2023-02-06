package client

type PipelineRun struct {
	Id     string
	Date   string
	Type   string
	Status string
}

type DBInterface interface {
	GetRuns() ([]PipelineRun, error)
	// GetRun(id string) string
}
