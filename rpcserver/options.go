package rpcserver

import (
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/usecase"
)

type Config struct {
	Interactors []usecase.Interactor
	ServerAddr  string
	Info        openapi3.Info
}

type Option func(c *Config)

func defaultConfig() *Config {
	return &Config{
		ServerAddr: ":8011",
	}
}

func applyOptions(cfg *Config, opts ...Option) *Config {
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func WithServerAddr(addr string) Option {
	return func(c *Config) {
		c.ServerAddr = addr
	}
}

func WithInteractors(interactor ...usecase.Interactor) Option {
	return func(c *Config) {
		c.Interactors = interactor
	}
}

func WithInfo(info *openapi3.Info) Option {
	return func(c *Config) {
		c.Info = *info
	}
}
