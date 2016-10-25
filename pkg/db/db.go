package db

import "github.com/llitfkitfk/GoHighPerformance/pkg/model"

type DB interface {

	Save(model.Key, model.Model) error

	Delete(model.Key) error

	Get(model.Key, model.Model) error
}