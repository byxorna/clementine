package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	SSLCertificate string     `json:"ssl-cert,omitempty"`
	SSLKey         string     `json:"ssl-key,omitempty"`
	CACertificate  string     `json:"ca-cert,omitempty"`
	SkipCAVerify   bool       `json:"skip-verify"`
	ListenAddr     string     `json:"listen-addr,omitempty"`
	TlsConfig      tls.Config `json:"-"`
	RedisHost      string     `json:"redis-host,omitempty"`
	RedisPort      int        `json:"redis-port"`
}

func LoadConfig(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c := Config{}
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, err
	}
	if err := c.loadTlsCrap(); err != nil {
		return nil, err
	}
	//TODO i dont know how to make a struct have default values. That would be nice...
	if c.RedisPort == 0 {
		c.RedisPort = 6379
	}
	return &c, nil
}

func (c *Config) loadTlsCrap() error {
	cert, err := tls.LoadX509KeyPair(c.SSLCertificate, c.SSLKey)
	if err != nil {
		return err
	}
	// read in our certs for the TLS connection
	certs := make([]tls.Certificate, 1)
	certs = append(certs, cert)
	// make sure we provide the CA cert as well
	rootCAPool := x509.NewCertPool()
	caf, err := os.Open(c.CACertificate)
	if err != nil {
		return err
	}
	cafInfo, _ := caf.Stat()
	caData := make([]byte, cafInfo.Size())
	if _, err := caf.Read(caData); err != nil {
		log.Fatal(err)
	}
	success := rootCAPool.AppendCertsFromPEM(caData)
	if !success {
		return fmt.Errorf("Unable to load Root CA cert %s", c.CACertificate)
	}

	c.TlsConfig = tls.Config{
		InsecureSkipVerify: c.SkipCAVerify,
		Certificates:       certs,
		RootCAs:            rootCAPool,
	}

	return nil

}
