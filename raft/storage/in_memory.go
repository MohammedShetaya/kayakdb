package storage

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/types"
)

// TODO: replace this with a desk implmentation

// InMemoryDriver is an in-memory implementation of the Driver interface
// using slices to store the log entries.
type InMemoryDriver struct {
	currentTerm uint
	votedFor    string
	log         []LogEntry
}

func NewInMemoryDriver() *InMemoryDriver {
	return &InMemoryDriver{
		log: make([]LogEntry, 0),
	}
}

func (d *InMemoryDriver) GetCurrentTerm() uint {
	return d.currentTerm
}

// SetCurrentTerm persist the value of the current term.
func (d *InMemoryDriver) SetCurrentTerm(term uint) error {
	d.currentTerm = term
	return nil
}

func (d *InMemoryDriver) GetVotedFor() string {
	return d.votedFor
}

// SetVotedFor persists the value of the last performed vote.
func (d *InMemoryDriver) SetVotedFor(candidate string) error {
	d.votedFor = candidate
	return nil
}

// Append appends a log entry and returns the log index.
func (d *InMemoryDriver) Append(entry LogEntry) uint {
	fmt.Println("appending")
	d.log = append(d.log, entry)
	return uint(len(d.log))
}

func (d *InMemoryDriver) AppendInIndex(index uint, entry LogEntry) uint {
	// overwrite if it is possible
	if index < uint(len(d.log)) {
		d.log[index] = entry
		return index
	}
	// otherwise append last
	return d.Append(entry)
}

func (d *InMemoryDriver) GetEntryOfIndex(index uint) *LogEntry {
	if index-1 > uint(len(d.log)) {
		return nil
	}
	return &d.log[index-1]
}

func (d *InMemoryDriver) AppendMany(startIndex uint, entries []LogEntry) error {
	if startIndex > uint(len(d.log)) {
		return fmt.Errorf("start index is larger than log length")
	}

	for i := startIndex; i < uint(len(entries))+startIndex; i++ {
		d.AppendInIndex(i, entries[i])
	}
	return nil
}

//// FindLastMatchingIndex the function finds the largest index that matches the incoming entries.
//// startIndex is the start index in the log where the log[startIndex+1] on the sender is the same as  entries[0]
//func (d *InMemoryDriver) FindLastMatchingIndex(startIndex uint, entries []LogEntry) (uint, error) {
//	if startIndex > uint(len(d.log)) {
//		return 0, fmt.Errorf("start index is larger than log length")
//	}
//
//	l := startIndex
//	h := startIndex + uint(len(entries))
//	mid := l + (h-l)/2
//
//	for l < h {
//		mid = l + (h-l)/2
//		if d.log[mid].Term != entries[mid-startIndex].Term {
//			h = mid - 1
//		} else {
//			l = mid + 1
//		}
//	}
//
//	return mid, nil
//}

// ConstructMappingFromLog constructs the map of the key-value pairs from the log entries.
func (d *InMemoryDriver) ConstructMappingFromLog() (map[string]types.Type, error) {
	mapping := make(map[string]types.Type)
	for _, entry := range d.log {
		mapping[string(entry.Pair.Key.Bytes())] = entry.Pair.Value
	}
	return mapping, nil
}
