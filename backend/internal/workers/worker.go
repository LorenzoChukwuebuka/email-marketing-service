package worker

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type TaskType string

const (
	TaskSendWelcomeEmail      TaskType = "email:send_welcome"
	TaskSendUserNotification  TaskType = "usernotification:send"
	TaskSendAdminNotification TaskType = "adminnotifcation:send"
	TaskSendEmail             TaskType = "email:send"
	TaskUserDetails           TaskType = "user:store_details"
	TaskSMTPSendEmail         TaskType = "smtp:sendemail"
	TaskAuditLogCreate        TaskType = "audit:log_create"
	TaskAuditLogUpdate        TaskType = "audit:log_update"
	TaskAuditLogDelete        TaskType = "audit:log_delete"
	TaskAuditLogLogin         TaskType = "audit:log_login"
	TaskAuditLogFailedLogin   TaskType = "audit:log_failed_login"
	TaskSendCampaignEmail     TaskType = "email:campaign"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

type Worker struct {
	Store           db.Store
	shutdownChan    chan struct{}
	wg              sync.WaitGroup
	maxRetries      int
	retryDelay      time.Duration
	pollInterval    time.Duration
	staleTaskWindow time.Duration
	mu              sync.Mutex
	isRunning       bool
}

type WorkerConfig struct {
	MaxRetries      int
	RetryDelay      time.Duration
	PollInterval    time.Duration
	StaleTaskWindow time.Duration
}

// NewWorker creates a new database-backed worker
func NewWorker(store db.Store, cfg WorkerConfig) *Worker {
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}
	if cfg.RetryDelay == 0 {
		cfg.RetryDelay = 5 * time.Second
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = 1 * time.Second
	}
	if cfg.StaleTaskWindow == 0 {
		cfg.StaleTaskWindow = 10 * time.Minute
	}

	return &Worker{
		Store:           store,
		shutdownChan:    make(chan struct{}),
		maxRetries:      cfg.MaxRetries,
		retryDelay:      cfg.RetryDelay,
		pollInterval:    cfg.PollInterval,
		staleTaskWindow: cfg.StaleTaskWindow,
	}
}

// Start begins the worker pool with the specified number of workers
func (w *Worker) Start(ctx context.Context, numWorkers int) {
	w.mu.Lock()
	if w.isRunning {
		w.mu.Unlock()
		log.Println("Worker already running")
		return
	}
	w.isRunning = true
	w.mu.Unlock()

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		w.wg.Add(1)
		go w.run(ctx, i)
	}

	// Start stale task recovery goroutine
	w.wg.Add(1)
	go w.recoverStaleTasks(ctx)

	log.Printf("Worker pool started with %d workers", numWorkers)
}

func (w *Worker) run(ctx context.Context, workerID int) {
	defer w.wg.Done()
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	log.Printf("Worker #%d started", workerID)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker #%d shutting down (context cancelled)", workerID)
			return
		case <-w.shutdownChan:
			log.Printf("Worker #%d shutting down", workerID)
			return
		case <-ticker.C:
			if err := w.processNextTask(ctx); err != nil {
				if err != sql.ErrNoRows {
					log.Printf("Worker #%d error: %v", workerID, err)
				}
			}
		}
	}
}

func (w *Worker) processNextTask(ctx context.Context) error {
	// Claim a task atomically using FOR UPDATE SKIP LOCKED
	task, err := w.Store.ClaimNextTask(ctx)
	if err != nil {
		return err
	}

	log.Printf("Processing task ID: %d, Type: %s, Retry: %d/%d",
		task.ID, task.TaskType, task.RetryCount, task.MaxRetries)

	// Process the task
	if err := w.ProcessTask(ctx, task); err != nil {
		// Calculate exponential backoff
		backoffMultiplier := 1 << uint(task.RetryCount)
		retryDelay := int(w.retryDelay.Seconds()) * backoffMultiplier

		// Mark as failed with retry scheduling
		if err := w.Store.MarkTaskFailed(ctx, db.MarkTaskFailedParams{
			ID:           task.ID,
			ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
			Column3:      sql.NullString{String: fmt.Sprintf("%d", retryDelay), Valid: true},
		}); err != nil {
			log.Printf("Failed to mark task as failed: %v", err)
		}

		log.Printf("Task %d failed (attempt %d/%d): %v",
			task.ID, task.RetryCount+1, task.MaxRetries, err)
		return err
	}

	// Mark as completed
	if err := w.Store.MarkTaskCompleted(ctx, task.ID); err != nil {
		log.Printf("Failed to mark task as completed: %v", err)
		return err
	}

	log.Printf("Task %d completed successfully", task.ID)
	return nil
}

func (w *Worker) recoverStaleTasks(ctx context.Context) {
	defer w.wg.Done()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.shutdownChan:
			return
		case <-ticker.C:
			staleMinutes := fmt.Sprintf("%d", int(w.staleTaskWindow.Minutes()))
			if err := w.Store.ResetStaleTasks(ctx, sql.NullString{String: staleMinutes, Valid: true}); err != nil {
				log.Printf("Failed to reset stale tasks: %v", err)
			} else {
				log.Println("Checked for stale tasks")
			}
		}
	}
}

// EnqueueTask adds a new task to the queue
func (w *Worker) EnqueueTask(ctx context.Context, taskType TaskType, payload interface{}) (int64, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal payload: %w", err)
	}

	task, err := w.Store.CreateTask(ctx, db.CreateTaskParams{
		TaskType:    string(taskType),
		Payload:     payloadJSON,
		Status:      string(TaskStatusPending),
		MaxRetries:  int32(w.maxRetries),
		ScheduledAt: time.Now(),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	log.Printf("Task enqueued: ID=%d, Type=%s", task.ID, task.TaskType)
	return task.ID, nil
}

// EnqueueTaskDelayed adds a task that will be processed after a delay
func (w *Worker) EnqueueTaskDelayed(ctx context.Context, taskType TaskType, payload interface{}, delay time.Duration) (int64, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal payload: %w", err)
	}

	task, err := w.Store.CreateTask(ctx, db.CreateTaskParams{
		TaskType:    string(taskType),
		Payload:     payloadJSON,
		Status:      string(TaskStatusPending),
		MaxRetries:  int32(w.maxRetries),
		ScheduledAt: time.Now().Add(delay),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	log.Printf("Delayed task enqueued: ID=%d, Type=%s, ScheduledAt=%s",
		task.ID, task.TaskType, task.ScheduledAt)
	return task.ID, nil
}

// Shutdown gracefully shuts down the worker
func (w *Worker) Shutdown(ctx context.Context) error {
	w.mu.Lock()
	if !w.isRunning {
		w.mu.Unlock()
		return nil
	}
	w.mu.Unlock()

	close(w.shutdownChan)

	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		w.mu.Lock()
		w.isRunning = false
		w.mu.Unlock()
		log.Println("Worker shutdown complete")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// GetStats returns worker statistics
func (w *Worker) GetStats(ctx context.Context) (map[string]interface{}, error) {
	pendingCount, err := w.Store.GetPendingTasksCount(ctx)
	if err != nil {
		return nil, err
	}

	processingCount, err := w.Store.GetProcessingTasksCount(ctx)
	if err != nil {
		return nil, err
	}

	failedCount, err := w.Store.GetFailedTasksCount(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"pending_tasks":    pendingCount,
		"processing_tasks": processingCount,
		"failed_tasks":     failedCount,
		"max_retries":      w.maxRetries,
		"retry_delay":      w.retryDelay.String(),
		"poll_interval":    w.pollInterval.String(),
	}, nil
}

// CleanupOldTasks removes completed tasks older than specified days
func (w *Worker) CleanupOldTasks(ctx context.Context, daysOld int) error {
	return w.Store.CleanupOldTasks(ctx, sql.NullString{String: fmt.Sprintf("%d", daysOld), Valid: true})
}
