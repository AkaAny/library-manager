package es

import "github.com/elastic/go-elasticsearch/v8"

type ESConfig struct {
	Addresses []string `toml:"addresses"`
}

func (esc ESConfig) ToClientConfig() elasticsearch.Config {
	var result = elasticsearch.Config{
		Addresses: esc.Addresses,
	}
	return result
}
