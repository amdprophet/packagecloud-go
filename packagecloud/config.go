package packagecloud

import (
	"errors"
	"fmt"
	"net/url"
)

type Config struct {
	ServiceURL string `mapstructure:"url"`
	Token      string `mapstructure:"token"`
	Verbose    bool   `mapstructure:"verbose"`
}

func (c Config) Validate() error {
	if c.Token == "" {
		return errors.New("token must not be empty")
	}

	if c.ServiceURL == "" {
		return errors.New("url must not be empty")
	}

	if _, err := url.Parse(c.ServiceURL); err != nil {
		return fmt.Errorf("invalid url: %s", err)
	}

	return nil
}
