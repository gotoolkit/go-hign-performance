package db

import (
	"github.com/llitfkitfk/GoHighPerformance/pkg/model"
	"gopkg.in/redis.v5"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) Save(model.Key, model.Model) error {
	return nil
}

func (r *Redis) Delete(model.Key) error {
	return nil
}

func (r *Redis) Get(model.Key, model.Model) error {
	return nil
}
