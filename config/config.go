package config

type Configuration struct {
	KayakPort      string   `json:"kayak_port" env:"KAYAK_PORT"`
	RaftPort       string   `json:"raft_port" env:"RAFT_PORT"`
	LogLevel       string   `json:"log_level" env:"LOG_LEVEL"`
	MaxLogBatch    uint     `json:"max_log_batch" env:"MAX_LOG_BATCH"`
	WorkerPoolSize uint     `json:"worker_pool_size" env:"WORKER_POOL_SIZE"`
	WaitQueueSize  uint     `json:"wait_queue_size" env:"WAIT_QUEUE_SIZE"`
	SeedPeers      []string `json:"seed_peers"`
}
