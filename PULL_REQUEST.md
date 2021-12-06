## Problem or Feature

Implement three (3) API endpoints to Create, Read (Download), and Delete arbitrary data from a virtual data store.

## Solution

1. Each API endpoint was implemented in a separate file for accessbility.
2. All the data storage logic was implemented in the same file (datastore.go)
3. Mutexes are used to coordinate read/write access to the in-memory data store
4. Only Go standard packages were used, even for the HTTP router

## Potential Improvements

1. I would make use of an external library like [httprouter](https://github.com/julienschmidt/httprouter) or [middleware](https://github.com/cixtor/middleware) to replace the rudimentary URL parsing operation that is happening in every API endpoint. This will also make the code easier to read and augment.
2. Instead of storing duplicate objects in separate repositories, I would create a centralized object store using a key-value data structure (HashMap) where the key is an ObjectID and the value is `[]byte`. Then, `DataStore.repositories` will remain as a key-value data structure (HashMap) where the key is a RepositoryID, but the value is now a list of ObjectIDs. This way, if a client uploads the same data to one-million different repositories, the server will only need to store one object because the SHA256 is the same.
3. Following the same idea, I would modify the Delete endpoint to be idempotent. Instead of returning a “404 Not Found” when a repository or object identifier has already been deleted (or never existed), I will return a “200 OK” to avoid unnecessry client retries.

## Verification

All existing unit tests pass.
