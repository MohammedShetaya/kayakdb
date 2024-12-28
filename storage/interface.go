package storage

import (
	"github.com/MohammedShetaya/kayakdb/types"
)

/*
 This is the interface that will be used by the raft lib and the api server
 to deal with the underlying storage.
*/

type Driver interface {
	// SetCurrentTerm persist the value of the current term
	SetCurrentTerm(term int) bool
	// SetVotedFor persists the value of the last performed vote
	SetVotedFor(term int) bool
	// Append appends a log entry and return the log index
	Append(entry types.KeyValue) int
	// ConstructMappingFromLog construct the map of the key value
	ConstructMappingFromLog() map[types.Type]types.Type
}
