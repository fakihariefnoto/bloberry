package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

// JobQueue is a simple Redis list-based FIFO queue (domains.md §7 job:queue).
type JobQueue struct {
	c *goredis.Client
}

func NewJobQueue(c *goredis.Client) *JobQueue {
	return &JobQueue{c: c}
}

const jobQueueKey = "job:queue"

func (q *JobQueue) Enqueue(ctx context.Context, jobID string) error {
	return q.c.LPush(ctx, jobQueueKey, jobID).Err()
}

func (q *JobQueue) Dequeue(ctx context.Context) (string, error) {
	return q.c.RPop(ctx, jobQueueKey).Result()
}
