package signed

import (
	"gitlab.com/tokend/connectors/keyer"
	"gitlab.com/tokend/keypair"
	"net/http"
	"net/url"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/keypair/figurekeypair"
)

type Clienter interface {
	Client() *Client
}

type clienter struct {
	getter kv.Getter
	keyer.Keyer
	once comfig.Once
}

func NewClienter(getter kv.Getter) *clienter {
	return &clienter{
		getter: getter,
		Keyer:  keyer.NewKeyer(getter),
	}
}

func (h *clienter) Client() *Client {
	return h.once.Do(func() interface{} {

		var config struct {
			Endpoint *url.URL        `fig:"endpoint,required"`
			Signer   keypair.Full    `fig:"signer"`
			Source   keypair.Address `fig:"source"`
		}

		err := figure.
			Out(&config).
			With(figure.BaseHooks, figurekeypair.Hooks).
			From(kv.MustGetStringMap(h.getter, "client")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out client"))
		}
		var keys keyer.Keys

		if config.Signer != nil {
			keys = keyer.Keys{
				Signer: config.Signer,
				Source: config.Source,
			}
		} else {
			keys = h.Keyer.Keys()
		}

		cli := NewClient(http.DefaultClient, config.Endpoint)
		if keys.Signer != nil {
			cli = cli.WithSigner(keys.Signer)
		}
		if keys.Source != nil {
			cli = cli.WithSource(keys.Source)
		}

		return cli
	}).(*Client)
}
