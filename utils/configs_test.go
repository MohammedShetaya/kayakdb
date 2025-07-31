package utils

import (
	"os"
	"testing"

	"github.com/MohammedShetaya/kayakdb/config"
)

func TestLoadConfigurations(t *testing.T) {
	t.Run("load with defaults only (no json file)", func(t *testing.T) {
		// Ensure no JSON file exists
		os.Remove("raft.json")

		cfg := &config.Configuration{}
		result, err := LoadConfigurations(cfg)
		if err != nil {
			t.Fatalf("Failed to load configurations: %v", err)
		}

		config := result.(*config.Configuration)

		// Verify default values are set
		if config.KayakPort != "8080" {
			t.Errorf("Expected default KayakPort to be 8080, got %s", config.KayakPort)
		}
		if config.RaftPort != "9090" {
			t.Errorf("Expected default RaftPort to be 9090, got %s", config.RaftPort)
		}
		if config.LogLevel != "info" {
			t.Errorf("Expected default LogLevel to be info, got %s", config.LogLevel)
		}
		if config.MaxLogBatch != 50 {
			t.Errorf("Expected default MaxLogBatch to be 50, got %d", config.MaxLogBatch)
		}
		if config.WorkerPoolSize != 4 {
			t.Errorf("Expected default WorkerPoolSize to be 4, got %d", config.WorkerPoolSize)
		}
		if config.WaitQueueSize != 1000 {
			t.Errorf("Expected default WaitQueueSize to be 1000, got %d", config.WaitQueueSize)
		}
		if config.PeerDiscovery != false {
			t.Errorf("Expected default PeerDiscovery to be false, got %v", config.PeerDiscovery)
		}
		if config.ServiceName != "kayakdb" {
			t.Errorf("Expected default ServiceName to be kayakdb, got %s", config.ServiceName)
		}
		// SeedPeers doesn't have a default tag, so it should remain nil/empty
		if config.SeedPeers != nil && len(config.SeedPeers) > 0 {
			t.Errorf("Expected SeedPeers to be nil or empty (no default), got %v", config.SeedPeers)
		}
	})

	t.Run("load from json only", func(t *testing.T) {
		// Create a temporary config file for testing
		const testConfigJSON = `{
			"kayak_port": "8080",
			"raft_port": "9000",
			"log_level": "info",
			"max_log_batch": 100,
			"worker_pool_size": 4,
			"wait_queue_size": 1000,
			"seed_peers": ["localhost:9001", "localhost:9002"]
		}`

		err := os.WriteFile("raft.json", []byte(testConfigJSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}
		defer os.Remove("raft.json")

		cfg := &config.Configuration{}
		result, err := LoadConfigurations(cfg)
		if err != nil {
			t.Fatalf("Failed to load configurations: %v", err)
		}

		config := result.(*config.Configuration)

		// Verify JSON values
		if config.KayakPort != "8080" {
			t.Errorf("Expected KayakPort to be 8080, got %s", config.KayakPort)
		}
		if config.RaftPort != "9000" {
			t.Errorf("Expected RaftPort to be 9000, got %s", config.RaftPort)
		}
		if config.LogLevel != "info" {
			t.Errorf("Expected LogLevel to be info, got %s", config.LogLevel)
		}
		if config.MaxLogBatch != 100 {
			t.Errorf("Expected MaxLogBatch to be 100, got %d", config.MaxLogBatch)
		}
		if len(config.SeedPeers) != 2 {
			t.Errorf("Expected 2 seed peers, got %d", len(config.SeedPeers))
		}
	})

	t.Run("override with environment variables", func(t *testing.T) {
		// Create a temporary config file for testing
		const testConfigJSON = `{
			"kayak_port": "8080",
			"raft_port": "9000",
			"log_level": "info",
			"max_log_batch": 100,
			"worker_pool_size": 4,
			"wait_queue_size": 1000,
			"seed_peers": ["localhost:9001", "localhost:9002"]
		}`

		err := os.WriteFile("raft.json", []byte(testConfigJSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}
		defer os.Remove("raft.json")

		// Set environment variables
		os.Setenv("KAYAK_PORT", "8081")
		os.Setenv("RAFT_PORT", "9001")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("MAX_LOG_BATCH", "200")
		os.Setenv("WORKER_POOL_SIZE", "8")
		os.Setenv("WAIT_QUEUE_SIZE", "2000")
		defer func() {
			os.Unsetenv("KAYAK_PORT")
			os.Unsetenv("RAFT_PORT")
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("MAX_LOG_BATCH")
			os.Unsetenv("WORKER_POOL_SIZE")
			os.Unsetenv("WAIT_QUEUE_SIZE")
		}()

		cfg := &config.Configuration{}
		result, err := LoadConfigurations(cfg)
		if err != nil {
			t.Fatalf("Failed to load configurations: %v", err)
		}

		config := result.(*config.Configuration)

		// Verify environment variables override JSON values
		if config.KayakPort != "8081" {
			t.Errorf("Expected KayakPort to be 8081, got %s", config.KayakPort)
		}
		if config.RaftPort != "9001" {
			t.Errorf("Expected RaftPort to be 9001, got %s", config.RaftPort)
		}
		if config.LogLevel != "debug" {
			t.Errorf("Expected LogLevel to be debug, got %s", config.LogLevel)
		}
		if config.MaxLogBatch != 200 {
			t.Errorf("Expected MaxLogBatch to be 200, got %d", config.MaxLogBatch)
		}
		if config.WorkerPoolSize != 8 {
			t.Errorf("Expected WorkerPoolSize to be 8, got %d", config.WorkerPoolSize)
		}
		if config.WaitQueueSize != 2000 {
			t.Errorf("Expected WaitQueueSize to be 2000, got %d", config.WaitQueueSize)
		}
		// SeedPeers should remain unchanged as it's not environment configurable
		if len(config.SeedPeers) != 2 {
			t.Errorf("Expected 2 seed peers, got %d", len(config.SeedPeers))
		}
	})

	t.Run("invalid json file", func(t *testing.T) {
		// Create invalid JSON
		err := os.WriteFile("raft.json", []byte("invalid json"), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid test config file: %v", err)
		}
		defer os.Remove("raft.json")

		cfg := &config.Configuration{}
		_, err = LoadConfigurations(cfg)
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})

	t.Run("invalid environment variable type", func(t *testing.T) {
		// Create a temporary config file for testing
		const testConfigJSON = `{
			"kayak_port": "8080",
			"raft_port": "9000",
			"log_level": "info",
			"max_log_batch": 100,
			"worker_pool_size": 4,
			"wait_queue_size": 1000,
			"seed_peers": ["localhost:9001", "localhost:9002"]
		}`

		err := os.WriteFile("raft.json", []byte(testConfigJSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}
		defer os.Remove("raft.json")

		// Set invalid environment variable
		os.Setenv("MAX_LOG_BATCH", "not_a_number")
		defer os.Unsetenv("MAX_LOG_BATCH")

		cfg := &config.Configuration{}
		_, err = LoadConfigurations(cfg)
		if err == nil {
			t.Error("Expected error for invalid environment variable type, got nil")
		}
	})

	t.Run("non-pointer config object", func(t *testing.T) {
		cfg := config.Configuration{}
		_, err := LoadConfigurations(cfg)
		if err == nil {
			t.Error("Expected error for non-pointer config object, got nil")
		}
	})

	t.Run("priority test: default < json < env", func(t *testing.T) {
		// Create a JSON file with some values
		const testConfigJSON = `{
			"kayak_port": "8888",
			"log_level": "debug",
			"max_log_batch": 200
		}`

		err := os.WriteFile("raft.json", []byte(testConfigJSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}
		defer os.Remove("raft.json")

		// Set one environment variable
		os.Setenv("KAYAK_PORT", "9999")
		defer os.Unsetenv("KAYAK_PORT")

		cfg := &config.Configuration{}
		result, err := LoadConfigurations(cfg)
		if err != nil {
			t.Fatalf("Failed to load configurations: %v", err)
		}

		config := result.(*config.Configuration)

		// Verify priority:
		// - KayakPort: has env var (9999), should override json (8888) and default (8080)
		if config.KayakPort != "9999" {
			t.Errorf("Expected KayakPort to be 9999 (from env), got %s", config.KayakPort)
		}

		// - LogLevel: has json (debug) but no env, should override default (info)
		if config.LogLevel != "debug" {
			t.Errorf("Expected LogLevel to be debug (from json), got %s", config.LogLevel)
		}

		// - MaxLogBatch: has json (200) but no env, should override default (50)
		if config.MaxLogBatch != 200 {
			t.Errorf("Expected MaxLogBatch to be 200 (from json), got %d", config.MaxLogBatch)
		}

		// - RaftPort: no json, no env, should use default (9090)
		if config.RaftPort != "9090" {
			t.Errorf("Expected RaftPort to be 9090 (from default), got %s", config.RaftPort)
		}

		// - WorkerPoolSize: no json, no env, should use default (4)
		if config.WorkerPoolSize != 4 {
			t.Errorf("Expected WorkerPoolSize to be 4 (from default), got %d", config.WorkerPoolSize)
		}
	})
}
