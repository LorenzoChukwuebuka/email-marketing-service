package tasks

// asynq lib uses this to know which queues and process belong together through the use of
// of payloads....
const (
	TaskSendWelcomeEmail      = "email:send_welcome"
	TaskSendUserNotification  = "usernotification:send"
	TaskSendAdminNotification = "adminnotifcation:send"
	TaskSendEmail             = "email:send"
)
