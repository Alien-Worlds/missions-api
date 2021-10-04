package config

import (
	"fmt"
	"os"
)

type EventsConfig struct {
	MissionJoinedHash string
	MissionCreatedHash string 
	RewardWithdrawnHash string
}

func (c *config) EventsConfig() EventsConfig {
	eventsConfig := new(EventsConfig)
	envResult, isSet := os.LookupEnv("MISSION_CREATED_HASH")
	if(isSet) {
		eventsConfig.MissionCreatedHash = envResult
	} else {
		panic("MISSION_CREATED_HASH no set")
	}
	envResult, isSet = os.LookupEnv("MISSION_JOINED_HASH")
	if(isSet) {
		eventsConfig.MissionJoinedHash = envResult
	} else {
		panic("MISSION_JOINED_HASH no set")
	}
	envResult, isSet = os.LookupEnv("REWARD_WITHDRAWN_HASH")
	if(isSet) {
		eventsConfig.RewardWithdrawnHash = envResult
	} else {
		panic("REWARD_WITHDRAWN_HASH no set")
	}
	fmt.Printf("eventsConfigs: %+v\n", eventsConfig)

	c.eventsConfig = *eventsConfig
	return c.eventsConfig
}
