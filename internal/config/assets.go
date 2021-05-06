package config

import (
	"github.com/pkg/errors"
	"github.com/redcuckoo/bsc-checker-events/internal/config/hooks"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

func (c *config) Assets() []hooks.Assets {
	c.once.Do(func() interface{} {
		var config struct {
			D []hooks.Assets `fig:"data"`
		}

		err := figure.
			Out(&config).
			With(figure.BaseHooks, hooks.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "assets")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out assets"))
		}

		c.assets = config.D
		return nil
	})

	return c.assets
}
