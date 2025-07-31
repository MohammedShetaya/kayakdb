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
	Append(entry LogEntry) uint
	AppendMany(startIndex uint, entries []LogEntry) error
	GetEntryOfIndex(index uint) *LogEntry
	//FindLastMatchingIndex(startIndex uint, entries []LogEntry) (uint, error)
	ConstructMappingFromLog() (map[string]types.Type, error)
}

type LogEntry struct {
	Term uint
	Pair types.KeyValue
}
