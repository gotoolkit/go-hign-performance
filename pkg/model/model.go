package model

import "encoding"

type Model interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	Set(Model) error
}
