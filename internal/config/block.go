package config

import (
	"fmt"
	"os"
	"strconv"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Block struct {
	FromBlockNum uint32
}

func (c *config) Block() Block {

	envResult, blockIsSet := os.LookupEnv("FROM_BLOCK_NUM")
	if(blockIsSet) {
		res, err := strconv.ParseInt(envResult,10,32)
		if(err != nil) {
			panic(errors.Wrap(err, "failed to get Block Num"))
		}
		fmt.Printf("Starting from Block Number: %+v\n", res)

		return Block{ uint32(res)}
	} else {
		panic("FROM_BLOCK_NUM no set")
	}
}
