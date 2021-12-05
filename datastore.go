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

// Read is the R in C.R.U.D and is responsible for finding and returning the
// data associated to an object identifier in a specific repository. If the
// repository or object does not exist, the method returns an error.
func (ds *DataStore) Read(repo RepositoryID, oid ObjectID) ([]byte, error) {
	ds.RLock()
	defer ds.RUnlock()

	if _, exists := ds.objects[repo]; !exists {
		return nil, errRepositoryNotFound
	}

	if _, exists := ds.objects[repo][oid]; !exists {
		return nil, errObjectNotFound
	}

	return ds.objects[repo][oid], nil
}

// Update is the U in C.R.U.D and is responsible for overriding the data linked
// to an object identifier. However, because the Create method already offers
// overrides, this method is not necessary.
func (ds *DataStore) Update(repo RepositoryID, data []byte) error {
	return errNotImplemented
}

// Delete is the D in C.R.U.D and is responsible for finding and deleting the
// data associated to an object identifier for a specific repository. Objects
// with the same identifier in a different repository are left unmodified. If
// the repository ends up empty after this object deletion, the repository is
// also deleted.
func (ds *DataStore) Delete(repo RepositoryID, oid ObjectID) error {
	ds.RLock()
	defer ds.RUnlock()

	if _, exists := ds.objects[repo]; !exists {
		return errRepositoryNotFound
	}

	if _, exists := ds.objects[repo][oid]; !exists {
		return errObjectNotFound
	}

	delete(ds.objects[repo], oid)

	// Delete the repository too, if empty.
	if len(ds.objects[repo]) == 0 {
		delete(ds.objects, repo)
	}

	return nil
}

// generateOID calculates the SHA256 sum of an arbitrary list of bytes. The
// function returns an object identifier with exactly sixty-four characters
// that is supposed to be unique.
func (ds *DataStore) generateOID(data []byte) ObjectID {
	hash := sha256.Sum256(data)
	oid := fmt.Sprintf("%x", hash[:])
	return ObjectID(oid)
}
