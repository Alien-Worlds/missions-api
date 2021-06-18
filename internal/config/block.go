package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

type Block struct {
	FromBlockNum uint32 `fig:"from_block_num,required"`
}

func (c *config) Block() Block {
	var result Block

	c.once.Do(func() interface{} {
		err := figure.
			Out(&result).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "blockchain_info")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out last block info"))
		}

		return nil
	})

	return result
}
