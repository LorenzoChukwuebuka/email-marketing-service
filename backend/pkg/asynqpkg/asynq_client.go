package asynqpkg

import (
	"github.com/hibiken/asynq"
	"sync"
)

var (
	client *asynq.Client
	once   sync.Once
	// DefaultRedisAddress defines the default Redis server address
	DefaultRedisAddress = "localhost:6379"
)

// Init initializes the Asynq client with the given Redis address.
// If no address is provided, it defaults to `DefaultRedisAddress`.
func Init(redisAddr ...string) {
	once.Do(func() {
		addr := DefaultRedisAddress
		if len(redisAddr) > 0 && redisAddr[0] != "" {
			addr = redisAddr[0]
		}
		client = asynq.NewClient(asynq.RedisClientOpt{Addr: addr})
	})
}

// GetClient returns the initialized Asynq client.
// Panics if the client is not initialized.
func GetClient() *asynq.Client {
	if client == nil {
		panic("asynqclient: Client is not initialized. Call Init() first.")
	}
	return client
}

func GetRedisAddress() string {
	if client == nil {
		return DefaultRedisAddress
	}
	// You might need to add a field to store the current address
	return DefaultRedisAddress
}

func NewInspector() *asynq.Inspector {
	return asynq.NewInspector(asynq.RedisClientOpt{
		Addr: GetRedisAddress(),
	})
}
