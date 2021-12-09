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

func (r *RedisEngine) Add(key string, content string) error {

	args := &redis.XAddArgs{
		MaxLen: 100000,
		Stream: key,
		Values: map[string]interface{}{
			"content": content,
		},
	}
	_, err := r.client.XAdd(ctx, args).Result()

	return err
}

func (r *RedisEngine) Read(key string, count int64, onRead func(contents []string) error) {

	r.read(key, 0, count, onRead)
}

func (r *RedisEngine) ReadBlock(key string, count int64, onRead func(contents []string) error) {

	for {

		err := r.read(key, time.Minute, count, onRead)

		if err != nil {

			if strings.Contains(err.Error(), "redis: nil") {
				continue
			}
			log.Println("read", err.Error())
			time.Sleep(time.Second * 3)
			continue
		}

	}

}

func (r *RedisEngine) read(key string, blockTime time.Duration, count int64, onRead func(contents []string) error) error {
	streams, err := r.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{key, "$"},
		Count:   count,
		Block:   blockTime,
	}).Result()

	if err != nil {

		return err
	}

	ids := make([]string, 0)
	ds := make([]string, 0)

	for _, s := range streams {

		for _, msg := range s.Messages {

			content, ok := msg.Values["content"]

			if !ok {
				continue
			}

			switch content.(type) {

			case string:
				ds = append(ds, content.(string))

				ids = append(ids, msg.ID)
			default:
				continue
			}

		}

	}

	if err := onRead(ds); err != nil {

		log.Println(err)
		return nil

	}

	r.client.XDel(ctx, key, ids...)

	return nil
}
