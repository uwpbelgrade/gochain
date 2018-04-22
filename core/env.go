package core

import (
	"log"
	"os"
	"strconv"
)

// DbFile env config
func DbFile() string {
	return os.Getenv("BOLT_DB_FILE")
}

// DbBucket env config
func DbBucket() string {
	return os.Getenv("BOLT_DB_BUCKET")
}

// BlockReward env config
func BlockReward() int {
	reward, err := strconv.ParseInt(os.Getenv("BLOCK_REWARD"), 10, 2)
	if err != nil {
		log.Fatal(err)
	}
	return int(reward)
}

// GenesisData env config
func GenesisData() string {
	return os.Getenv("GENESIS_DATA")
}
