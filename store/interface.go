package store

type Crud interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Durable interface {
	Save(filename ...string) error
	Load(filename ...string) error
}


