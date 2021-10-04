package config

import (
	"fmt"
	"os"
)

type ContractAddress struct {
	Address string
}

func (c *config) Contract() ContractAddress {
	envResult, isSet := os.LookupEnv("CONTRACT_ADDRESS")
	if(isSet) {
		fmt.Printf("Reading from Contract address: %+v\n", envResult)

		return ContractAddress{ envResult }
	} else {
		panic("CONTRACT_ADDRESS no set")
	}
}
