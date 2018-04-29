package core

import (
	"log"
	"os"
	"strconv"
)

// Config specifies configuration properties
type Config interface {
	GetDbFile() string
	GetDbBucket() string
	GetBlockReward() int
	GetGenesisData() string
	GetWalletStoreFile() string
}

// EnvConfig implements Config via environment
type EnvConfig struct{}

// GetDbFile gets BOLT_DB_FILE
func (env *EnvConfig) GetDbFile() string {
	return env.Get("BOLT_DB_FILE")
}

// GetDbBucket gets BOLT_DB_BUCKET
func (env *EnvConfig) GetDbBucket() string {
	return env.Get("BOLT_DB_BUCKET")
}

// GetBlockReward gets BLOCK_REWARD
func (env *EnvConfig) GetBlockReward() int {
	return env.GetInt("BLOCK_REWARD")
}

// GetGenesisData gets GENESIS_DATA
func (env *EnvConfig) GetGenesisData() string {
	return env.Get("GENESIS_DATA")
}

// GetWalletStoreFile gets WALLET_STORE_FILE
func (env *EnvConfig) GetWalletStoreFile() string {
	return env.Get("WALLET_STORE_FILE")
}

// Get gets string value from config
func (env *EnvConfig) Get(key string) string {
	return os.Getenv(key)
}

// GetInt gets intiger value from config
func (env *EnvConfig) GetInt(key string) int {
	reward, err := strconv.ParseInt(os.Getenv(key), 10, 64)
	if err != nil {
		log.Fatalf("error getting %s: %s", key, err)
	}
	return int(reward)
}
