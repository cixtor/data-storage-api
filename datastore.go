package main

import (
	"sync"
)

// RepositoryID represents a unique repository identifier.
type RepositoryID string

// ObjectID represents a unique object identifier.
type ObjectID string

type ObjectData struct {
	// data is a binary representation of some object.
	data []byte

	// links represents the number of repositories to which the object has been
	// associated. When the Delete API endpoint identifies an object in memory,
	// it decreases this value by one, and when the value reaches zero, it
	// deletes the object entirely from the data storage.
	links int
}

// DataStore is the core of the web service.
type DataStore struct {
	// Since the web service is storing the objects in memory and the server is
	// handling the HTTP requests concurrently, it should utilize mutexes to
	// avoid data races where one endpoint is reading/writing some data while
	// another endpoint is also trying to access or modify that data.
	sync.RWMutex

	// objects is a key-value data structure where the key represents a unique
	// object identifier and the value is an array of bytes. The map guarantees
	// de-duplication of objects across all repositories.
	objects map[RepositoryID]map[ObjectID][]byte
}

// NewDataStore creates a new instance of DataStore and initializes all its
// internal attributes, including, but not limited to a map of data objects.
func NewDataStore() *DataStore {
	ds := new(DataStore)
	ds.objects = map[RepositoryID]map[ObjectID][]byte{}
	return ds
}
