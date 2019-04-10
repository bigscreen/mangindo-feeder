package appcontext

import (
	"fmt"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/worker/adapter"
	"github.com/go-redis/redis"
	redigo "github.com/gomodule/redigo/redis"
)

type appContext struct {
	redisClient   *redis.Client
	workerAdapter adapter.Worker
}

var context *appContext

func Initiate() {
	context = &appContext{
		redisClient:   initRedisClient(config.RedisHost(), config.RedisPort(), config.RedisPool()),
		workerAdapter: initWorker(),
	}
}

func initRedisClient(host string, port int, poolSize int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
		PoolSize: poolSize,
	})
}

func initWorker() *adapter.Adapter {
	return adapter.NewAdapter(adapter.Options{
		Pool:           GetWorkerRedisPool(),
		Name:           constants.WorkerName,
		MaxConcurrency: 10,
	})
}

func GetWorkerRedisPool() *redigo.Pool {
	return &redigo.Pool{
		MaxActive: 25,
		MaxIdle:   25,
		Wait:      true,
		Dial: func() (redigo.Conn, error) {
			conn, err := redigo.Dial("tcp", config.WorkerRedisAddress())
			return conn, err
		},
	}
}

func GetRedisClient() *redis.Client {
	return context.redisClient
}

func GetWorkerAdapter() adapter.Worker {
	return context.workerAdapter
}
