package config

import (
	"fmt"
	"os"

	bscClient "github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *config) BSC() bscClient.Client {
	envResult, isSet := os.LookupEnv("RPC_END_POINT")
	if(isSet) {
		dial, err := bscClient.Dial(envResult)
		if err != nil {
			panic(errors.Wrap(err, "failed to dial bscClient"))
		}
		c.bsc = *dial
		fmt.Printf("Connecting to BSC endpoint: %v\n", dial)

	}
	return c.bsc
}
