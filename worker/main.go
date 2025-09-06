package main

import (
	"os"

	"github.com/hibiken/asynq"
)

func main() {
	logger := NewLogger()

	// 1. Create a new database connection pool.
	dbPool, dbPoolErr := newDBConnectionPool(logger)
	if dbPoolErr != nil {
		logger.Error("Failed to connect to database pool.", "err", dbPoolErr)
		return
	}
	defer dbPool.Close()

	// 2. Create a new Redis connection client.
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		logger.Error("FATAL: REDIS_URL environment variable is not set!")
		return
	}

	redisConnOpt, redisErr := asynq.ParseRedisURI(redisURL)
	if redisErr != nil {
		logger.Error("Failed to parse Redis URI.", "err", redisErr)
		return
	}

	redisClientOpt, ok := redisConnOpt.(asynq.RedisClientOpt)
	if !ok {
		logger.Error("Failed to convert Redis connection options to RedisClientOpt type.")
		return
	}

	// 3. Create a new TaskProcessor.
	processor := NewTaskProcessor(redisClientOpt, dbPool, logger)

	// 4. Start the processor.
	logger.Info("Starting worker service...")
	if err := processor.Start(); err != nil {
		logger.Error("Could not run worker server.", "err", err)
		return
	}
}
