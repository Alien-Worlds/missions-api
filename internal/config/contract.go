package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

type ContractAddress struct {
	Address string `fig:"address,required"`
}

func (c *config) Contract() ContractAddress {
	c.once.Do(func() interface{} {
		var result ContractAddress

		err := figure.
			Out(&result).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "contract")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out contract"))
		}

		c.contract = result
		return nil
	})

	return c.contract
}
