package models

type Model interface {
	GetId() int64
	GetStructMap() *map[string]interface{}
}
