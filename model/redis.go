package model

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
	"github.ewr01.tumblr.net/gabe/clementine/config"
)

func NewRedisPool(host string, port int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 1000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host+":"+strconv.Itoa(port))
			if err != nil {
				return nil, fmt.Errorf("Error dialing redis at %s:%d: %s", host, port, err.Error())
			}
			return c, err
		},
	}
}

func (e *Engine) New(p *redis.Pool) error {
	e.ID = uuid.New()
	log.Printf("Assigned ID %s to engine %s\n", e.ID, e.Engine.ID)
	return e.Save(p)
}

func (e *Engine) Save(p *redis.Pool) error {
	if e.ID == "" {
		return fmt.Errorf("No ID set! You should call save the engine first!")
	}
	conn := p.Get()
	defer conn.Close()
	js, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("Unable to marshal engine to JSON: %s", err.Error())
	}
	log.Printf("Saving engine %s\n", e.ID)
	if _, err := conn.Do("SET", fmt.Sprintf("engine:%s", e.ID), js); err != nil {
		log.Printf("Error when saving %s to redis: %s\n", e.ID, err.Error())
		return err
	}
	//TODO this is ghetto, but... make sure each save sets the ID into the set of engines
	_, err = conn.Do("SADD", "engines", e.ID)
	return err
}

func Get(p *redis.Pool, id string, c *config.Config) (*Engine, error) {
	conn := p.Get()
	defer conn.Close()
	res, err := conn.Do("SISMEMBER", "engines", id)
	if err != nil {
		return nil, fmt.Errorf("Error getting engine %s: %s", id, err.Error())
	}
	if res.(int64) == 0 {
		log.Printf("No engine found with ID %s", id)
		return nil, fmt.Errorf("No engine found with ID %s", id)
	}

	log.Printf("Fetching engine %s\n", id)
	if b, err := conn.Do("GET", fmt.Sprintf("engine:%s", id)); err != nil {
		log.Printf("Error when fetching engine %s from redis: %s\n", id, err.Error())
		return nil, err
	} else {
		if b == nil {
			log.Printf("No engine found with ID %s, despite engine being listed in engines set", id)
			return nil, fmt.Errorf("Found engine with ID %s in list, but cannot find data", id)
		}
		var engine Engine
		if err := json.Unmarshal(b.([]byte), &engine); err != nil {
			log.Printf("Error unmarshalling engine %s: %s %s\n", id, err.Error(), b)
			return nil, fmt.Errorf("Error unmarshalling engine %s: %s", id, err.Error())
		}
		return &engine, nil
	}
}

func AllEngines(p *redis.Pool, c *config.Config) ([]*Engine, error) {
	conn := p.Get()
	defer conn.Close()
	res, err := conn.Do("SMEMBERS", "engines")
	engineIds, err := redis.Strings(res, err)
	if err != nil {
		return nil, err
	}
	log.Printf("Found %d engines in index; querying for models\n", len(engineIds))
	engines := []*Engine{}
	for _, id := range engineIds {
		e, err := Get(p, id, c)
		if err != nil {
			log.Printf("Error getting engine %s: %s\n", id, err.Error())
		}
		engines = append(engines, e)
	}
	return engines, nil
}
