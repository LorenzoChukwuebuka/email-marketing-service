package tasks

// import (
// 	"email-marketing-service/internals/workers/payloads"
// 	"encoding/json"

// 	"github.com/google/uuid"
// 	"github.com/hibiken/asynq"
// )

// func NewStoreUserDetailsTask(details map[string]interface{}) (*asynq.Task, error) {
// 	payload, err := json.Marshal(payloads.UserDetailsPayload{Details: details})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return asynq.NewTask(TaskUserDetails, payload, asynq.TaskID(uuid.New().String())), nil
// }
