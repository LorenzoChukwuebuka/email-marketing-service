package cronjobs

import (
	"context"
	db "email-marketing-service/internal/db/sqlc"
)

type Job interface {
	Run() error
	Schedule() string
	Name() string
	Type() string
	Description() string
	SetStore(store db.Store)
}

// JobResult contains execution results
type JobResult struct {
	Success          bool                   `json:"success"`
	Message          string                 `json:"message,omitempty"`
	Data             map[string]interface{} `json:"data,omitempty"`
	RecordsProcessed int                    `json:"records_processed,omitempty"`
}

// BaseJob provides common functionality for all jobs
type BaseJob struct {
	store       db.Store
	ctx         context.Context
	jobName     string
	jobType     string
	description string
}

func NewBaseJob(store db.Store, ctx context.Context, name, jobType, desc string) *BaseJob {
	return &BaseJob{
		store:       store,
		ctx:         ctx,
		jobName:     name,
		jobType:     jobType,
		description: desc,
	}
}

func (b *BaseJob) Name() string            { return b.jobName }
func (b *BaseJob) Type() string            { return b.jobType }
func (b *BaseJob) Description() string     { return b.description }
func (b *BaseJob) SetStore(store db.Store) { b.store = store }
