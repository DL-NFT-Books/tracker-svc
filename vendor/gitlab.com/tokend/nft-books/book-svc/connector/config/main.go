package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/nft-books/book-svc/connector/api"
)

const bookerYamlKey = "booker"

type Booker interface {
	Connector() *api.Connector
}

type booker struct {
	once   comfig.Once
	getter kv.Getter
}

type bookerCfg struct {
	URL   string `fig:"url,required"`
	Token string `fig:"token,required"`
}

func NewBooker(getter kv.Getter) Booker {
	return &booker{getter: getter}
}

func (c *booker) Connector() *api.Connector {
	return c.once.Do(func() interface{} {
		var cfg bookerCfg

		raw := kv.MustGetStringMap(c.getter, bookerYamlKey)

		if err := figure.
			Out(&cfg).
			From(raw).
			Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		return api.NewConnector(cfg.Token, cfg.URL)
	}).(*api.Connector)
}
