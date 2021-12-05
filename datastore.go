package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

// errRepositoryNotFound is an error when the repository does not exist.
var errRepositoryNotFound = errors.New("repository not found")

// errObjectNotFound is an error when the object does not exist.
var errObjectNotFound = errors.New("object not found")

// errNotImplemented is an error when a function is not implemented.
var errNotImplemented = errors.New("not implemented")

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

// Create is the C in C.R.U.D. and is responsible for adding new objects into
// the in-memory data storage. Because the object identifiers are based on the
// content of the object, overrides are possible but not destructive.
func (ds *DataStore) Create(repo RepositoryID, data []byte) ObjectID {
	ds.Lock()
	defer ds.Unlock()

	// Record the new repository, if necessary.
	if _, exists := ds.objects[repo]; !exists {
		ds.objects[repo] = map[ObjectID][]byte{}
	}

	oid := ds.generateOID(data)

	ds.objects[repo][oid] = data

	return oid
}

// generateOID calculates the SHA256 sum of an arbitrary list of bytes. The
// function returns an object identifier with exactly sixty-four characters
// that is supposed to be unique.
func (ds *DataStore) generateOID(data []byte) ObjectID {
	hash := sha256.Sum256(data)
	oid := fmt.Sprintf("%x", hash[:])
	return ObjectID(oid)
}
