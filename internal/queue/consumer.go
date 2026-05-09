package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	MainQueue = "mail:queue"
	DLQ       = "mail:dlq"
)

type RedisQueuer struct {
	client *redis.Client
}

func NewRedisQueuer(client *redis.Client) *RedisQueuer {
	return &RedisQueuer{client: client}
}

func (r *RedisQueuer) Pop(ctx context.Context) (*MailJob, error) {
	result, err := r.client.BLPop(ctx, 0, MainQueue).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("errore lettura da Redis: %w", err)
	}

	var job MailJob
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, fmt.Errorf("errore deserializzazione job: %w", err)
	}

	return &job, nil
}

func (r *RedisQueuer) Requeue(job MailJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("errore serializzazione job: %w", err)
	}
	return r.client.RPush(context.Background(), MainQueue, data).Err()
}

func (r *RedisQueuer) MoveToDLQ(job MailJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("errore serializzazione job: %w", err)
	}
	return r.client.RPush(context.Background(), DLQ, data).Err()
}
