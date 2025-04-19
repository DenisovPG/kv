package storage

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"sync"

	"kv/infrastructure/repository/tx_log"
)

type Store struct {
	data map[string]string
	sync.RWMutex
	logger tx_log.TxLogger
}

var ErrorNoSuchKey = errors.New("no such key")

func (s *Store) Put(key string, value string) error {
	s.Lock()
	defer s.Unlock()
	s.logger.LogPut(key, value)
	s.data[key] = value
	return nil
}

func (s *Store) Get(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	value, ok := s.data[key]
	s.logger.LogGet(key)
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}

func (s *Store) Delete(key string) error {
	s.Lock()
	defer s.Unlock()
	s.logger.LogDelete(key)
	delete(s.data, key)
	return nil
}

func (s *Store) Save(filename ...string) error {
	fn := "memory_store"
	if len(filename) > 0 {
		fn = filename[0]
	}
	file, err := os.Create(fn)
	if err != nil {
		fmt.Println("Create file error: ", err)
		return err
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	if err = enc.Encode(s.data); err != nil {
		fmt.Println("Serialisation error: ", err)
		return err
	}
	s.logger.Stop()
	return nil
}

func (s *Store) Load(filename ...string) error {
	fn := "memory_store"
	if len(filename) > 0 {
		fn = filename[0]
	}
	file, err := os.Open(fn)
	if err != nil {
		fmt.Println("Open file error: ", err)
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(&s.data); err != nil {
		fmt.Println("Deserialisation error: ", err)
		return err
	}
	return nil
}

func NewMemoryStore(tx_logger tx_log.TxLogger) DurableOrCrud {
	tx_logger.Run()
	store := &Store{data: make(map[string]string), logger: tx_logger}
	store.Load()
	return store
}
