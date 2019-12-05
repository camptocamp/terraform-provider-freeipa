package rancher

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/tehwalris/go-freeipa/freeipa"
)

// Config is the configuration parameters for the FreeIPA API
type Config struct {
	Host               string
	Username           string
	Password           string
	InsecureSkipVerify bool
}

// Client creates a FreeIPA client scoped to the global API
func (c *Config) Client() (*freeipa.Client, error) {
	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.InsecureSkipVerify,
		},
	}

	client, err := freeipa.Connect(c.Host, tspt, c.Username, c.Password)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] FreeIPA Client configured for host: %s", c.Host)

	return client, nil
}
