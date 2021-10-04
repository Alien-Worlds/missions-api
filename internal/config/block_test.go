package config

import (
	"os"
	"testing"
)


func TestFetchBlock(t *testing.T) {
	os.Setenv("FROM_BLOCK_NUM","1234")
	config := new(config)
	block := config.Block()

	if block.FromBlockNum == 456456 {
		t.Error("failed to get block")
	}
	
}