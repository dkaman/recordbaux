package db

type Repository[T any] interface {
	All() ([]T, error)
	Get(uint) (T, error)
	Save(entity T) error
	Delete(uint) error
}
