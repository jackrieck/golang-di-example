package storage

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

// Client abstracts the datastore
type Client interface {
	Get(key string) (string, error)
}

type client struct {
	memcacheClient *memcache.Client
}

// New initializes a new storage client
func New() Client {
	return &client{
		memcacheClient: memcache.New("localhost:11211"),
	}
}

// Get takes a key and returns a string value or an error
func (c *client) Get(key string) (string, error) {

	// query memcache for item matching given key
	item, err := c.memcacheClient.Get(key)
	if err != nil {
		return "", err
	}

	// check if item is nil, return error means it wasn't found
	if item == nil {
		return "", errors.New("item not found")
	}

	return string(item.Value), nil
}
