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
	state map[types.Type]types.Type
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
	return s.state[key], nil
}

func (s *State) GetMajority() int {
	return s.peersCount()/2 + 1
}

func (s *State) peersCount() int {
	return len(s.peers)
}

func (s *State) GetLogsRange(start uint, end uint) []types.LogEntry {
	return nil
}
