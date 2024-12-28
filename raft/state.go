package raft

import "github.com/MohammedShetaya/kayakdb/storage"

type State struct {
	// persistent state
	persistentStorage storage.Driver
	// volatile state
	commitIndex int
	lastApplied int
}

type LeaderState struct {
	State
	nextIndex  int
	matchIndex int
}
