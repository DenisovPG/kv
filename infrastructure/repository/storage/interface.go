package storage

type Crud interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type DurableCrud interface {
	Crud
	Save(filename ...string) error
	Load(filename ...string) error
}



type DurableOrCrud interface {
	DurableCrud
	Crud
}