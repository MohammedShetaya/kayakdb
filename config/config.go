package config

type Configuration struct {
	KayakPort      string   `json:"kayak_port" env:"KAYAK_PORT" default:"8080"`
	RaftPort       string   `json:"raft_port" env:"RAFT_PORT" default:"9090"`
	LogLevel       string   `json:"log_level" env:"LOG_LEVEL" default:"info"`
	MaxLogBatch    uint     `json:"max_log_batch" env:"MAX_LOG_BATCH" default:"50"`
	WorkerPoolSize uint     `json:"worker_pool_size" env:"WORKER_POOL_SIZE" default:"4"`
	WaitQueueSize  uint     `json:"wait_queue_size" env:"WAIT_QUEUE_SIZE" default:"1000"`
	PeerDiscovery  bool     `json:"peer_discovery" env:"PEER_DISCOVERY" default:"false"`
	ServiceName    string   `json:"service_name" env:"SERVICE_NAME" default:"kayakdb"`
	SeedPeers      []string `json:"seed_peers"`
}
