package config

import (
	"fmt"
	"net"
	"os"
)

func (c *config) Listener() net.Listener {
	envResult, envIsSet := os.LookupEnv("LISTEN_ADDRESS")
	if(envIsSet) {
		listener, err := net.Listen("tcp", envResult)
		if (err == nil) {
			fmt.Printf("Listening on Address and Port: %+v\n", envResult)

			return listener
		} 
	}
	panic("LISTEN_ADDRESS not set")	
}
