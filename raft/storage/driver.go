package storage

import (
	"github.com/MohammedShetaya/kayakdb/types"
)

// Driver This is the interface that will be used by the raft lib to deal with the underlying storage.
type Driver interface {
	GetCurrentTerm() uint
	SetCurrentTerm(term uint) error
	GetVotedFor() string
	SetVotedFor(candidate string) error
	Append(entry types.LogEntry) uint
	AppendMany(startIndex uint, entries []types.LogEntry) error
	GetEntryOfIndex(index uint) *types.LogEntry
	//FindLastMatchingIndex(startIndex uint, entries []types.LogEntry) (uint, error)
	ConstructMappingFromLog() (map[types.Type]types.Type, error)
}
