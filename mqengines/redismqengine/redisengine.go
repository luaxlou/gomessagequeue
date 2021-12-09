package redismqengine

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/luaxlou/gomessagequeue/mqengines"
	"log"
	"strings"
	"time"
)

var ctx = context.Background()

type RedisEngine struct {
	client *redis.Client
}

func New(addr, password string, db int) (mqengines.MqEngine, error) {

	opt := &redis.Options{
		Addr:        addr,
		Password:    password, // no password set
		DB:          db,       // use default DB
		DialTimeout: time.Minute,
		ReadTimeout: time.Minute,
		IdleTimeout: time.Minute,
	}

	c := redis.NewClient(opt)

	if c == nil {
		return nil, errors.New("Redis connect failed:" + addr)
	}

	return &RedisEngine{client: c}, nil
}
func NewWithClient(c *redis.Client) (mqengines.MqEngine, error) {

	if c == nil {
		return nil, errors.New("Redis client must be not nil.")
	}

	return &RedisEngine{client: c}, nil
}

func (r *RedisEngine) Add(key string, data interface{}) error {

	args := &redis.XAddArgs{
		MaxLen: 100000,
		Stream: key,
		Values: map[string]interface{}{
			"ts":   time.Now().Unix(),
			"data": data,
		},
	}
	_, err := r.client.XAdd(ctx, args).Result()

	return err
}

func (r *RedisEngine) Read(key string, count int64, onRead func(ds []interface{}) error) {

	for {
		streams, err := r.client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{key, "0"},
			Count:   count,
			Block:   time.Minute,
		}).Result()

		if err != nil || len(streams) == 0 {

			if strings.Contains(err.Error(), "redis: nil") {
				continue
			}
			log.Println("read", err.Error())
			time.Sleep(time.Second * 3)
			continue
		}

		ids := make([]string, 0)
		ds := make([]interface{}, 0)

		for _, s := range streams {

			for _, msg := range s.Messages {

				data, ok := msg.Values["data"]

				if !ok {
					continue
				}

				ds = append(ds, data)

				ids = append(ids, msg.ID)

			}

		}

		if err := onRead(ds); err != nil {

			log.Println(err)
			continue

		}

		r.client.XDel(ctx, key, ids...)

	}

}
