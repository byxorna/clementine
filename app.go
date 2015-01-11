package main

import (
	"log"
	"net/http"

	conf "github.com/byxorna/clementine/config"
	"github.com/byxorna/clementine/model"
	cc "github.com/citadel/citadel/cluster"
	"github.com/citadel/citadel/scheduler"
	r "github.com/garyburd/redigo/redis"
)

var (
	self *App
)

type App struct {
	config  conf.Config
	cluster *cc.Cluster
	redis   *r.Pool
}

func Setup(config conf.Config) error {
	cluster, _ := cc.New(scheduler.NewResourceManager())

	var (
		labelScheduler  = &scheduler.LabelScheduler{}
		uniqueScheduler = &scheduler.UniqueScheduler{}
		hostScheduler   = &scheduler.HostScheduler{}
		portScheduler   = &scheduler.PortScheduler{}

		multiScheduler = scheduler.NewMultiScheduler(
			labelScheduler,
			uniqueScheduler,
			portScheduler,
		)
	)
	cluster.RegisterScheduler("service", labelScheduler)
	cluster.RegisterScheduler("unique", uniqueScheduler)
	cluster.RegisterScheduler("multi", multiScheduler)
	cluster.RegisterScheduler("host", hostScheduler)
	cluster.RegisterScheduler("port", portScheduler)

	self = &App{
		cluster: cluster,
		config:  config,
		redis:   model.NewRedisPool(config.RedisHost, config.RedisPort),
	}

	// lets just try and connect to redis quickly, and fail if we cannot
	log.Printf("Testing connection to redis at %s:%d\n", self.config.RedisHost, self.config.RedisPort)
	conn := self.redis.Get()
	defer conn.Close()
	if res, err := conn.Do("PING"); err != nil {
		return err
	} else {
		log.Printf("Got response from redis: %s\n", res)
	}
	return nil
}

// return a new cluster with our schedulers already registered
func (self *App) Cluster() cc.Cluster {
	return *self.cluster
}
