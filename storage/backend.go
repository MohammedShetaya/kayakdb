package storage

import "github.com/MohammedShetaya/kayakdb/types"

// TODO: replace this with a desk implmentation

// SliceDriver is an in-memory implementation of the Driver interface
// using slices to store the log entries.
type SliceDriver struct {
	currentTerm int
	votedFor    int
	log         []types.KeyValue
}

func NewSliceDriver() *SliceDriver {
	return &SliceDriver{
		log: make([]types.KeyValue, 0),
	}
}

// SetCurrentTerm persist the value of the current term.
func (d *SliceDriver) SetCurrentTerm(term int) bool {
	d.currentTerm = term
	return true
}

// SetVotedFor persists the value of the last performed vote.
func (d *SliceDriver) SetVotedFor(term int) bool {
	d.votedFor = term
	return true
}

// Append appends a log entry and returns the log index.
func (d *SliceDriver) Append(entry types.KeyValue) int {
	d.log = append(d.log, entry)
	return len(d.log) - 1
}

// ConstructMappingFromLog constructs the map of the key-value pairs
// from the log entries.
func (d *SliceDriver) ConstructMappingFromLog() map[types.Type]types.Type {
	mapping := make(map[types.Type]types.Type)
	for _, entry := range d.log {
		mapping[entry.Key] = entry.Value
	}
	return mapping
}
