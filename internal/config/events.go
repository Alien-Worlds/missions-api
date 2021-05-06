package config

import (
	"github.com/stellar/go/support/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

type EventsConfig struct {
	MissionJoinedHash string `fig:"mission_joined_hash,required"`
	MissionCreatedHash string `fig:"mission_created_hash,required"`
	RewardWithdrawnHash string `fig:"reward_withdrawn_hash,required"`
}

func (c *config) EventsConfig() EventsConfig {
	c.once.Do(func() interface{} {
		var config EventsConfig

		err := figure.
			Out(&config).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "events")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out EventConfig"))
		}

		c.eventsConfig = config
		return nil
	})

	return c.eventsConfig
}
