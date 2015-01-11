package model

import (
	"github.com/byxorna/clementine/config"
	"github.com/citadel/citadel"
	"github.com/samalba/dockerclient"
)

// engine statuses
// these should be stored separately from the engines
const (
	Pending EngineStatus = iota
	OK
	Maintenance
	Down
)

type (
	EngineStatus int
	Engine       struct {
		ID     string          `json:"id,omitempty"`
		Engine *citadel.Engine `json:"engine,omitempty"`
	}
)

func (e *Engine) SetupClient(c *config.Config) error {
	client, err := dockerclient.NewDockerClient(e.Engine.Addr, &c.TlsConfig)
	if err != nil {
		return err
	}
	e.Engine.SetClient(client)
	return nil
}
