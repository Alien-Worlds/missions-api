package config

import (
	bscClient "github.com/binance-chain/bsc/ethclient"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *config) BSC() bscClient.Client {
	c.once.Do(func() interface{} {
		var result struct {
			Endpoint string `fig:"endpoint,required"`
		}

		err := figure.
			Out(&result).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "rpc")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out stellar"))
		}

		dial, err := bscClient.Dial(result.Endpoint)
		if err != nil {
			panic(errors.Wrap(err, "failed to dial bscClient"))
		}

		c.bsc = *dial
		return nil
	})

	return c.bsc
}