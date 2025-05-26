package utils

import (
	"os"
	"testing"

	"github.com/MohammedShetaya/kayakdb/config"
)

func TestLoadConfigurations(t *testing.T) {
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

	t.Run("load from json only", func(t *testing.T) {
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

		cfg := &config.Configuration{}
		_, err = LoadConfigurations(cfg)
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})

	t.Run("invalid environment variable type", func(t *testing.T) {
		// Restore valid JSON
		err := os.WriteFile("raft.json", []byte(testConfigJSON), 0644)
		if err != nil {
			t.Fatalf("Failed to restore test config file: %v", err)
		}

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
}
