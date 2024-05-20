package main

import (
	"strconv"
	"sync"
	"time"
)

type Storage struct {
	data  map[string]Set
	mutex sync.RWMutex
}

type Set struct {
	value      string
	expires    bool
	expireTime time.Time
}

// Returns a new storage.
func NewStore() *Storage {
	storage := &Storage{
		data:  map[string]Set{},
		mutex: sync.RWMutex{},
	}
	return storage
}

// Sets a new value.
func (s *Storage) Set(key string, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = Set{
		value:   value,
		expires: false,
	}
	return nil
}

// Sets a new expiring value.
//
// It will expire after x miliseconds provided by `expiresMillis`.
func (s *Storage) SetWithExpiration(key string, value string, expiresMillis string) error {
	time, err := getExpirationTime(expiresMillis)

	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = Set{
		value:      value,
		expires:    true,
		expireTime: time,
	}

	return nil
}

// Retrieves value from the storage.
func (s *Storage) Get(key string) string {
	s.mutex.RLock()
	data, ok := s.data[key]
	s.mutex.RUnlock()

	if !ok {
		return ""
	}

	if data.expires && data.expireTime.Before(time.Now()) {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		// Double-check if the key still exists and if it's still expired
		data, ok = s.data[key]
		if ok && data.expires && data.expireTime.Before(time.Now()) {
			delete(s.data, key)
			return ""
		}
	}

	return data.value
}

// Returns expiration time
func getExpirationTime(expiresMillis string) (time.Time, error) {
	milis, err := strconv.Atoi(expiresMillis)
	if err != nil {
		return time.Time{}, err
	}

	expirationTime := time.Now().Add(time.Duration(milis) * time.Millisecond)
	return expirationTime, nil
}
