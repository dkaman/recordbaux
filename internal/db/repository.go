package db

import "github.com/google/uuid"

type ID = uuid.UUID

type Repository[T any] interface {
	All() ([]T, error)
	Get(ID) (T, error)
	Save(entity T) error
	Delete(ID) error
}
