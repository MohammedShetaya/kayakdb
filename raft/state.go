package raft

import (
	"github.com/MohammedShetaya/kayakdb/raft/storage"
	"github.com/MohammedShetaya/kayakdb/types"
	"net/rpc"
	"sync"
	"time"
)

type Peer struct {
	addr   string
	client *rpc.Client
	mutex  sync.Mutex

	// leader specific state
	nextIndex  uint
	matchIndex uint
}

type State struct {
	Persistent storage.Driver
	// volatile state
	CommitIndex uint // last committed log entry. initialized to 0 (0 is not considered a log index)
	LastApplied uint // last applied to the state map

	peers []Peer

	ServerId       string
	IsLeader       bool
	cancelElection chan struct{}
	FollowerTimer  *time.Timer

	// constructed key-value map from the log
	// TODO: use swap and disk (lru based)
	state map[string]types.Type
}

func NewState(driver storage.Driver) *State {
	s := &State{
		Persistent: driver,
	}
	constructedState, err := s.Persistent.ConstructMappingFromLog()
	if err != nil {
		return nil
	}
	s.state = constructedState
	return s
}

func (s *State) Get(key types.Type) (types.Type, error) {
	// TODO: after implementing swapping make sure to retrieve cold values
	return s.state[string(key.Bytes())], nil
}

//func (s *State) Put(key types.Type, val types.Type) error {
//
//}

func (s *State) GetMajority() int {
	return s.peersCount()/2 + 1
}

func (s *State) peersCount() int {
	return len(s.peers)
}

func (s *State) GetLogsRange(start uint, end uint) []storage.LogEntry {
	// Return a slice containing log entries in the inclusive range [start,end].
	// If end < start an empty slice is returned.
	if end < start {
		return nil
	}
	var result []storage.LogEntry
	for idx := start; idx <= end; idx++ {
		if entry := s.Persistent.GetEntryOfIndex(idx); entry != nil {
			result = append(result, *entry)
		}
	}
	return result
}

// ApplyNewEntries applies all log entries that have been committed but not yet applied
// to the in-memory state map. After execution LastApplied will equal CommitIndex.
func (s *State) ApplyNewEntries() {
	for idx := s.LastApplied + 1; idx <= s.CommitIndex; idx++ {
		entry := s.Persistent.GetEntryOfIndex(idx)
		if entry == nil {
			continue
		}
		s.state[string(entry.Pair.Key.Bytes())] = entry.Pair.Value
		s.LastApplied = idx
	}
}
