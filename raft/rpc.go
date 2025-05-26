package raft

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/types"
	"go.uber.org/zap"
)

type VoteRequest struct {
	Term         uint
	CandidateId  string
	LastLogIndex uint
	LastLogTerm  uint
}

type AppendRequest struct {
	Term         uint
	LeaderId     string
	PrevLogIndex uint
	PreLogTerm   uint
	LeaderCommit uint
	Entries      []types.LogEntry
}

type PingRequest struct {
	Term uint
}

type VoteResponse struct {
}

type AppendResponse struct {
	CommitIndex uint
}

// PingResponse no need for extra fields on a ping
type PingResponse struct {
}

type RpcController struct {
	logger *zap.Logger
	raft   *Raft
}

func NewRpcController(raft *Raft, logger *zap.Logger) *RpcController {
	controller := new(RpcController)
	controller.raft = raft
	controller.logger = logger

	return controller
}

func (c *RpcController) Vote(request VoteRequest, response *VoteResponse) error {
	c.logger.Debug("Received a voting request from", zap.String("candidate", request.CandidateId))
	if request.Term < c.raft.State.Persistent.GetCurrentTerm() {
		return fmt.Errorf("request term is less than current term of: %v", c.raft.State.Persistent.GetCurrentTerm())
	}
	if c.raft.State.Persistent.GetVotedFor() == "" &&
		c.raft.State.CommitIndex <= request.LastLogIndex {
		// this is unset when a newer term is received or if the server became a candidate
		err := c.raft.State.Persistent.SetVotedFor(request.CandidateId)
		if err != nil {
			c.logger.Debug(fmt.Sprintf("Failed to set votedFor %v", request.CandidateId), zap.Error(err))
			return fmt.Errorf("server failed to persist the voted for value")
		}
	}
	return nil
}

func (c *RpcController) Append(request AppendRequest, response *AppendResponse) error {
	c.logger.Debug("Received an append request from", zap.String("leader", request.LeaderId))
	// check if there is a log entry on this server at index PrevLogIndex that has the same PreLogTerm
	// also check if the leader is not expired
	if l := c.raft.State.Persistent.GetEntryOfIndex(request.PrevLogIndex); request.Term < c.raft.State.Persistent.GetCurrentTerm() || (l != nil && l.Term != request.PreLogTerm) {
		return fmt.Errorf("missmatched logs or terms")
	}
	// at this point this server should become a follower. become one if not.
	c.raft.compareTerms(request.Term)
	c.raft.resetFollowerTimer()

	// now the incoming logs are checked to be valid. Append them all and override if there are other logs at the same index.
	err := c.raft.State.Persistent.AppendMany(request.PrevLogIndex+1, request.Entries)
	if err != nil {
		c.logger.Debug("Error appending logs", zap.Error(err))
		return fmt.Errorf("error appending logs: %w", err)
	}

	if request.LeaderCommit > c.raft.State.CommitIndex {
		// set the commit index to min(leader commit, index of last received log)
		c.raft.State.CommitIndex = min(request.LeaderCommit, request.PrevLogIndex+uint(len(request.Entries)))
	}

	response.CommitIndex = c.raft.State.CommitIndex

	return nil
}

// Ping should be used by leaders to send empty heartbeats to followers in case the log is synced with the leader
func (c *RpcController) Ping(request PingRequest, response *PingResponse) error {
	c.logger.Debug("Received a ping request from leader")

	// don't accept anything from an expired leader
	if request.Term < c.raft.State.Persistent.GetCurrentTerm() {
		return fmt.Errorf("request term is less than current term of: %v", c.raft.State.Persistent.GetCurrentTerm())
	}
	c.raft.compareTerms(request.Term)
	c.raft.resetFollowerTimer()

	return nil
}
