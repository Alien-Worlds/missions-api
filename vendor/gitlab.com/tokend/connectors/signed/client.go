package signed

import (
	"net/http"
	"net/url"
	"path"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	depkeypair "gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/signcontrol"
	"gitlab.com/tokend/keypair"
)

func throttle() chan time.Time {
	burst := 2 << 10
	ch := make(chan time.Time, burst)

	go func() {
		tick := time.Tick(1 * time.Second)
		// prefill buffer
		for i := 0; i < burst; i++ {
			ch <- time.Now()
		}
		for {
			select {
			case ch <- <-tick:
			default:
			}
		}
	}()
	return ch
}

type Client struct {
	base     *url.URL
	signer   keypair.Full
	throttle chan time.Time
	source   keypair.Address
	// Client must only be called from Client.Do() methods only
	client *http.Client
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {
	<-c.throttle

	// ensure content-type just in case
	request.Header.Set("content-type", "application/json")

	if c.signer != nil {
		seed := depkeypair.MustParse(c.signer.Seed())

		request.Header.Set("account-id", c.signer.Address())
		if c.source != nil {
			request.Header.Set("account-id", c.source.Address())
		}

		err := signcontrol.SignRequest(request, seed)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to sign request")
		}
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to perform http request")
	}

	return response, nil
}

func NewClient(client *http.Client, base *url.URL) *Client {
	return &Client{
		base, nil, throttle(), nil, client,
	}
}

func (c *Client) WithSigner(kp keypair.Full) *Client {
	return &Client{
		c.base,
		kp,
		c.throttle,
		c.source,
		http.DefaultClient,
	}
}

func (c *Client) WithSource(source keypair.Address) *Client {
	return &Client{
		c.base,
		c.signer,
		c.throttle,
		source,
		http.DefaultClient,
	}
}

func (c *Client) Resolve(endpoint *url.URL) (string, error) {
	u := *c.base
	basePath := u.Path
	prevPath := endpoint.Path

	if basePath != "" {
		endpoint.Path = path.Join(basePath, endpoint.Path)
		u.Path = ""
	}

	resolved := u.ResolveReference(endpoint)
	endpoint.Path = prevPath
	return resolved.String(), nil
}
