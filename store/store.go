package store

import (
	"errors"
	"hash"

	"github.com/euforia/thrap/thrapb"
)

var (
	errIDMissing = errors.New("ID missing")
	ErrRefExists = errors.New("ref exists")
)

//
// type Config struct {
// 	ObjectStorage ObjectStorage
// 	HashFunc      func() hash.Hash
// }
//
// func DefaultConfig() *Config {
// 	return &Config{
// 		HashFunc: sha256.New,
// 	}
// }

// Object is a datastructure that is hashable and protobuf friendly
type Object interface {
	// Computes the hash of the object given the hash function
	Hash(h hash.Hash) []byte
	// Marshals the object to byte slice
	Marshal() ([]byte, error)
	// Unmarshals byte slice to object
	Unmarshal(b []byte) error
}

// ObjectStorage implements a namespaced object storage interface
type ObjectStorage interface {
	// CreateRef creates a new ref under the namespace.  Previous should be the
	// zero hash
	CreateRef(namespace, ref string) ([]byte, *thrapb.ChainHeader, error)
	// Set the given ref under the namespace to the Header
	SetRef(namespace, ref string, robj *thrapb.ChainHeader) ([]byte, error)
	// Returns the chain header for the ref
	GetRef(namespace, ref string) (*thrapb.ChainHeader, []byte, error)
	// Iterate over each reference in the namespace
	IterRefs(namepace string, callback func(string, []byte) error) error
	// Sets an object returning the hash which the object is stored under
	Set(namespace string, obj Object) ([]byte, error)
	// Populates obj under a namespace by the digest
	Get(namepace string, digest []byte, obj Object) error
}
