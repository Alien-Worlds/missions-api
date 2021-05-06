package hooks

import (
	"reflect"

	"github.com/binance-chain/bsc/common"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Assets struct {
	ContractAddress common.Address `fig:"contract_address,required"`
}

var (
	BaseHooks = figure.Hooks{
		"[]hooks.Assets": func(value interface{}) (reflect.Value, error) {
			assets := make([]Assets, 0)
			switch s := value.(type) {
			case []interface{}:
				var err error
				for _, elem := range s {
					mapa, ok := elem.(map[interface{}]interface{})
					if !ok {
						return reflect.Value{}, errors.Wrap(err,
							"failed to cast mapa to interface")
					}

					var data Assets

					contractAddress, ok := mapa["contract_address"].(string)
					if !ok {
						return reflect.Value{}, errors.Wrap(err, "failed to get string smart_contract")
					}

					data.ContractAddress = common.HexToAddress(contractAddress)

					assets = append(assets, data)
				}
				return reflect.ValueOf(assets), nil
			default:
				return reflect.ValueOf(value), errors.New("failed to get case")
			}
		},
	}
)
