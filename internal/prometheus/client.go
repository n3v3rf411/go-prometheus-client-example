package prometheus

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/prometheus/client_golang/api"
)

// NewClient returns a new Client.
//
// It is safe to use the returned Client from multiple goroutines.
func NewClient(cfg api.Config, beforeDo func(req *http.Request)) (api.Client, error) {
	u, err := url.Parse(cfg.Address)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimRight(u.Path, "/")

	return &httpClient{
		endpoint: u,
		client:   http.Client{Transport: cfg.RoundTripper},
		beforeDo: beforeDo,
	}, nil
}

type httpClient struct {
	endpoint *url.URL
	client   http.Client
	beforeDo func(req *http.Request)
}

func (c *httpClient) URL(ep string, args map[string]string) *url.URL {
	p := path.Join(c.endpoint.Path, ep)

	for arg, val := range args {
		arg = ":" + arg
		p = strings.Replace(p, arg, val, -1)
	}

	u := *c.endpoint
	u.Path = p

	return &u
}

func (c *httpClient) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	if c.beforeDo != nil {
		c.beforeDo(req)
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	resp, err := c.client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return nil, nil, err
	}

	var body []byte
	done := make(chan struct{})
	go func() {
		body, err = ioutil.ReadAll(resp.Body)
		close(done)
	}()

	select {
	case <-ctx.Done():
		<-done
		err = resp.Body.Close()
		if err == nil {
			err = ctx.Err()
		}
	case <-done:
	}

	return resp, body, err
}
