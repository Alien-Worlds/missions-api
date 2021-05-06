package config

import (
	"github.com/binance-chain/bsc/common"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

func (c *config) Contract() common.Address {
	c.once.Do(func() interface{} {
		var result string

		err := figure.
			Out(&result).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "contract")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out contract"))
		}

		c.contract = common.HexToAddress(result)
		return nil
	})

	return c.contract
}
