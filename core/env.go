package core

import (
	"fmt"
	"os"
	"strconv"
)

// Config specifies configuration properties
type Config interface {
	GetDbFile(nodeID string) string
	GetDbBucket() string
	GetDbUtxoBucket() string
	GetBlockReward() int
	GetGenesisData() string
	GetWalletStoreFile(nodeID string) string
}

// EnvConfig implements Config via environment
type EnvConfig struct{}

// GetDbFile gets BOLT_DB_FILE
func (env *EnvConfig) GetDbFile(nodeID string) string {
	return fmt.Sprintf(env.Get("BOLT_DB_FILE"), nodeID)
}

// GetDbBucket gets BOLT_DB_BUCKET
func (env *EnvConfig) GetDbBucket() string {
	return env.Get("BOLT_DB_BUCKET")
}

// GetDbUtxoBucket gets BOLT_DB_UTXO_BUCKET
func (env *EnvConfig) GetDbUtxoBucket() string {
	return env.Get("BOLT_DB_UTXO_BUCKET")
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
func (env *EnvConfig) GetWalletStoreFile(nodeID string) string {
	return fmt.Sprintf(env.Get("WALLET_STORE_FILE"), nodeID)
}

// Get gets string value from config
func (env *EnvConfig) Get(key string) string {
	return os.Getenv(key)
}

// GetInt gets intiger value from config
func (env *EnvConfig) GetInt(key string) int {
	value, _ := strconv.ParseInt(os.Getenv(key), 10, 64)
	return int(value)
}
