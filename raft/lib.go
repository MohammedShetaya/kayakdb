package raft

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/config"
	"github.com/MohammedShetaya/kayakdb/raft/storage"
	"github.com/MohammedShetaya/kayakdb/types"
	"github.com/MohammedShetaya/kayakdb/utils"
	guuid "github.com/google/uuid"
	"go.uber.org/zap"
	"math/rand"
	"net"
	"net/rpc"
	"reflect"
	"sync/atomic"
	"time"
)

type Raft struct {
	logger     *zap.Logger
	config     *config.Configuration
	State      *State
	workerPool utils.WorkerPool
}

func NewRaft(config *config.Configuration, logger *zap.Logger) *Raft {
	raft := Raft{
		logger:     logger,
		config:     config,
		workerPool: utils.NewWorkerPool(config.WorkerPoolSize, config.WaitQueueSize),
	}

	raft.State = NewState(storage.NewInMemoryDriver())

	var p []string
	if raft.config.PeerDiscovery {
		p = utils.PeerServiceDiscovery(raft.config.ServiceName, raft.config.RaftPort)
	} else {
		p = raft.config.SeedPeers
	}
	raft.State.peers = make([]Peer, len(p))
	for i, addr := range p {
		raft.State.peers[i] = Peer{
			addr: addr,
		}
	}

	raft.State.ServerId = guuid.NewString()
	raft.State.FollowerTimer = time.NewTimer(time.Duration(rand.Intn(150)+150) * time.Millisecond)
	return &raft
}

func (r *Raft) Start() {
	listener, err := net.Listen("tcp", ":"+r.config.RaftPort)
	defer listener.Close()

	if err != nil {
		r.logger.Fatal("Failed to start background server", zap.Error(err))
	}

	r.logger.Info(fmt.Sprintf("Raft Server is listening on port: %v", r.config.RaftPort))

	r.workerPool.Start()

	r.logger.Info("Worker Pool has started")

	rpcController := NewRpcController(r, r.logger)
	err = rpc.Register(rpcController)
	if err != nil {
		r.logger.Info("Failed to register rpc controller", zap.Error(err))
		return
	}

	go r.registerNode()

	for {
		conn, err := listener.Accept()
		if err != nil {
			r.logger.Warn("Error accept connection from remote server", zap.Error(err))
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func (r *Raft) registerNode() {

	r.State.peers = make([]Peer, len(r.config.SeedPeers))
	for i, addr := range r.config.SeedPeers {
		r.State.peers[i] = Peer{
			addr: addr,
		}
	}

	// if the server just started try to start an election instead of looking who is the current leader
	// TODO: change this after implementing cluster membership change
	r.startElection()

	nodeTimeout := time.Duration(rand.Intn(150)+150) * time.Millisecond // 150-300 ms
	r.State.FollowerTimer.Reset(nodeTimeout)

	for {
		if r.State.IsLeader {
			// this periodically updates the followers to make sure they are in sync with the leader.
			// it shouldn't be frequently sending updates since new entries are populated upon receiving them by the leader
			// it will mostly be effective when leadership change happens or a new server joins the cluster.
			for i := 0; i < len(r.State.peers); i++ {
				go func(peerIdx int) {
					// check if the follower log is not up to date.
					if r.State.peers[peerIdx].matchIndex != r.State.CommitIndex {
						// then send log updates.
						diff := r.State.CommitIndex - r.State.peers[peerIdx].nextIndex

						// this will be nil if diff>0. In that case no logs are expected to be sent, just
						var logsToSend []storage.LogEntry
						if diff > 0 {
							batchSize := min(diff, r.config.MaxLogBatch)
							logsToSend = r.State.GetLogsRange(r.State.peers[peerIdx].nextIndex, r.State.peers[peerIdx].nextIndex+batchSize)
						}

						request := AppendRequest{
							Term:         r.State.Persistent.GetCurrentTerm(),
							LeaderId:     r.State.ServerId,
							PrevLogIndex: r.State.peers[peerIdx].nextIndex - 1,
							LeaderCommit: r.State.CommitIndex,
							Entries:      logsToSend,
						}

						response := new(AppendResponse)
						job, err := utils.NewJob(
							// the main job to execute
							sendRPC,
							[]any{
								&r.State.peers[peerIdx],
								"RpcController.Append",
								request,
								response,
							},
							// post exec function
							func(
								peer uint,
								raft *Raft,
								response *AppendResponse,
								// the return values of the main job that will be appended upon job execution
								jobReturns ...any,
							) {
								// update the matchIndex of that peer after a successful replication.
								if jobReturns[0] == nil {
									raft.logger.Debug(fmt.Sprintf("Append message has been sent to follower: %v", raft.State.peers[peer].addr))
									raft.State.peers[peer].matchIndex = response.CommitIndex
									raft.State.peers[peer].nextIndex = response.CommitIndex + 1
								} else {
									raft.logger.Error(fmt.Sprintf("Append RPC to followr: %v has failed", raft.State.peers[peer].addr), zap.Error(jobReturns[0].(error)))
									raft.State.peers[peer].nextIndex = raft.State.peers[peer].nextIndex - 1 // try older logs to find a matching point
								}
							}, []any{peerIdx, r, response})

						if err != nil {
							r.logger.Error("Failed to create append job", zap.Error(err))
						}

						err = r.workerPool.Enqueue(job)
						if err != nil {
							r.logger.Debug("Failed to enqueue append job", zap.Error(err))
						}

					} else { // otherwise, just ping the follower
						request := PingRequest{Term: r.State.Persistent.GetCurrentTerm()}
						response := new(PingResponse)

						job, err := utils.NewJob(
							sendRPC,
							[]any{
								&r.State.peers[peerIdx],
								"RpcController.Ping",
								request,
								response,
							},
							nil, nil)

						if err != nil {
							r.logger.Error("Failed to create ping job", zap.Error(err))
						}

						err = r.workerPool.Enqueue(job)
						if err != nil {
							r.logger.Debug("Failed to enqueue ping job", zap.Error(err))
						}
					}
				}(i)
			}

		} else {
			select {
			// if the leader didn't send a message for too long, start an election.
			case <-r.State.FollowerTimer.C:
				r.startElection()
			}
		}

	}
}

func (r *Raft) startElection() {

	r.logger.Info("Starting a new Election")
	err := r.State.Persistent.SetCurrentTerm(r.State.Persistent.GetCurrentTerm() + 1)
	if err != nil {
		r.logger.Debug("Failed to set term", zap.Error(err))
		return
	}

	err = r.State.Persistent.SetVotedFor(r.State.ServerId)
	if err != nil {
		r.logger.Debug(fmt.Sprintf("Server failed to vote for itself ID: %v", r.State.ServerId), zap.Error(err))
		return
	}

	var votes = atomic.Uint64{}
	// vote for itself
	votes.Add(1)

	// the channel to signal a vote upon follower response
	signal := make(chan struct{})
	// this makes sure that the current election is not infected by previous elections terminations
	r.State.cancelElection = make(chan struct{})

	request := VoteRequest{
		Term:         r.State.Persistent.GetCurrentTerm(),
		LastLogIndex: r.State.CommitIndex,
		CandidateId:  r.State.ServerId,
	}

	for i := 0; i < len(r.State.peers); i++ {
		go func(peerIdx uint) {
			response := new(VoteResponse)

			job, err := utils.NewJob(
				// the main job to execute
				sendRPC,
				[]any{
					&r.State.peers[peerIdx],
					"RpcController.Vote",
					request,
					response,
				},
				// post job
				func(
					peer uint,
					raft *Raft,
					vote *chan struct{},
					jobReturns ...any,
				) {
					if jobReturns[0] != nil {
						raft.logger.Error(fmt.Sprintf("Vote RPC to follower: %v has failed", raft.State.peers[peer].addr), zap.Error(jobReturns[0].(error)))
					}
					// at this point, vote was granted since there are no errors
					*vote <- struct{}{}
				}, []any{peerIdx, r, &signal},
			)

			if err != nil {
				r.logger.Error("Failed to create vote job", zap.Error(err))
			}

			err = r.workerPool.Enqueue(job)
			if err != nil {
				r.logger.Debug("Failed to enqueue vote job", zap.Error(err))
			}
		}(uint(i))

	}

	electionTimeout := time.Duration(rand.Intn(150)+150) * time.Millisecond // 150-300 ms
	r.logger.Info("Election will timeout after", zap.Duration("timeout_ms", electionTimeout))

	timer := time.NewTimer(electionTimeout)
	for {
		select {
		case <-timer.C:
			r.logger.Info("Election timed out!")
			return
		case <-signal:
			votes.Add(1)
		case <-r.State.cancelElection:
			r.logger.Info("Received an election cancellation signal, canceling election! until next time.")
			return
		default:
			if int(votes.Load()) >= r.State.GetMajority() {
				r.logger.Info("Became leader with majority of votes", zap.Int("votes", int(votes.Load())))
				r.State.IsLeader = true
				// initialized to leaders last log + 1
				for i := 0; i < len(r.State.peers); i++ {
					r.State.peers[i].matchIndex = 0
					r.State.peers[i].nextIndex = r.State.CommitIndex + 1
				}
				return
			}
		}
	}

}

func (r *Raft) compareTerms(term uint) {
	if term > r.State.Persistent.GetCurrentTerm() {
		err := r.State.Persistent.SetCurrentTerm(term)
		if err != nil {
			r.logger.Debug("Failed to set term", zap.Error(err))
		}
		err = r.State.Persistent.SetVotedFor("")
		if err != nil {
			r.logger.Debug("Failed to unset votedFor", zap.Error(err))
		}

		// if leader/candidate, then become a follower
		r.State.IsLeader = false
		if r.State.cancelElection != nil {
			r.State.cancelElection <- struct{}{}
		}
	}
}

func (r *Raft) resetFollowerTimer() {
	// reset the follower timer to avoid starting a new election
	// This will be invoked when term of a sender = term of the current server, either from prev cycles or after setting a new term
	r.State.FollowerTimer.Reset(time.Duration(rand.Intn(150)+150) * time.Millisecond) // 150-300 ms
}

func sendRPC(peer *Peer, remoteFunc string, request any, response any) error {
	if reflect.ValueOf(response).Kind() != reflect.Ptr {
		return fmt.Errorf("response must be a pointer")
	}

	peer.mutex.Lock()
	defer peer.mutex.Unlock()

	// Reuse the client if it's already connected, otherwise create new one
	if peer.client == nil {
		client, err := rpc.Dial("tcp", peer.addr)
		if err != nil {
			return fmt.Errorf("unable to connect to peer %s: %w", peer.addr, err)
		}
		peer.client = client
	}

	err := peer.client.Call(remoteFunc, request, response)
	if err != nil {
		return fmt.Errorf("error calling rpc: %w", err)
	}

	return nil
}

// Put handles the logic for putting a new entry on the leader or redirecting the command to the current leader on followers.
// it cannot be async since the user will be waiting for a response.
func (r *Raft) Put(data []types.Type) []storage.LogEntry {
	// TODO: check if this server is a leader, if not forward to the current leader for now assume that put will be only called on leaders

	// create a log entry of the new values
	var entries []storage.LogEntry
	var lastIndex uint // will hold the index of the last appended log entry

	for _, kv := range data {
		if pair, ok := kv.(types.KeyValue); ok {
			entry := storage.LogEntry{
				Term: r.State.Persistent.GetCurrentTerm(),
				Pair: pair,
			}
			idx := r.State.Persistent.Append(entry)
			r.State.CommitIndex = r.State.CommitIndex + 1
			lastIndex = idx
			entries = append(entries, entry)
		} else {
			r.logger.Error("unable to assert to KeyValue")
		}
	}

	commits := atomic.Uint64{}
	// add this leader server
	commits.Add(1)

	signal := make(chan struct{})

	// populate the new entry to followers
	for i := 0; i < len(r.State.peers); i++ {
		go func(peerIdx int, raft *Raft) {

			request := AppendRequest{
				Term:         r.State.Persistent.GetCurrentTerm(),
				LeaderId:     r.State.ServerId,
				PrevLogIndex: r.State.peers[peerIdx].nextIndex - 1,
				LeaderCommit: r.State.CommitIndex,
				Entries:      entries,
			}

			response := new(AppendResponse)

			err := sendRPC(&r.State.peers[peerIdx],
				"RpcController.Append",
				request,
				response)
			if err != nil {
				raft.logger.Debug(fmt.Sprintf("error appending new entry to follower: %v", peerIdx), zap.Error(err))
				return
			}

			signal <- struct{}{}

		}(i, r)
	}

	// TODO: rethink this one since it's dangerous in the senario when leader is not able to get a majority to commit the entry
	for {
		select {
		case <-signal:
			commits.Add(1)
		default:
			if int(commits.Load()) >= r.State.GetMajority() {
				// majority has acknowledged, apply entries
				r.State.CommitIndex = lastIndex
				r.State.ApplyNewEntries()
				return entries
			}
		}
	}

}

func (r *Raft) Get(key types.Type) (types.Type, error) {
	// If this node is not the leader, in a fully-fledged implementation we would
	// forward the request to the current leader. For the time being – until
	// redirection logic is implemented – we simply try to satisfy the request
	// from the local constructed state map. This will work correctly when the
	// request is sent to the current leader and during single-node deployments.
	return r.State.Get(key)
}
