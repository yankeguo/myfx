package myfx

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedsync(r *redis.Client) *redsync.Redsync {
	return redsync.New(goredis.NewPool(r))
}
