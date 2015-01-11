/*
All the actions associated with cluster management
i.e. adding and removing engines, listing engines, etc
*/
package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/byxorna/clementine/model"
	"github.com/citadel/citadel"
	"github.com/gorilla/mux"
)

/*
  Engine actions
*/

func engines(w http.ResponseWriter, r *http.Request) {
	engines, err := model.AllEngines(c.redis, &c.config)
	if err != nil {
		jsonError(w, fmt.Sprintf("Error fetching engines: %s", err.Error()))
		return
	}
	jsonSuccess(w, engines)
}

func (c *Controller) CreateEngine(e *citadel.Engine) (*model.Engine, error) {
	engine := model.Engine{Engine: e}
	if err := engine.SetupClient(&c.config); err != nil {
		return nil, err
	}
	if !engine.Engine.IsConnected() {
		return nil, fmt.Errorf("Unable to connect to engine %s", e.ID)
	}
	if err := engine.New(c.redis); err != nil {
		return nil, err
	}
	return &engine, nil
}

func createEngine(w http.ResponseWriter, r *http.Request) {
	var engine *citadel.Engine
	if err := json.NewDecoder(r.Body).Decode(&engine); err != nil {
		jsonError(w, err.Error())
		return
	}
	if e, err := c.CreateEngine(engine); err != nil {
		jsonError(w, err.Error())
	} else {
		log.Printf("Engine %s created: %s\n", e.ID, e.Engine)
		jsonSuccessStatus(w, e, http.StatusCreated)
	}
}

func getEngine(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	e, err := model.Get(c.redis, id, &c.config)
	if err != nil {
		jsonError(w, err.Error())
		return
	}
	jsonSuccess(w, e)
}

/*
  Cluster actions
*/
func clusterInfo(w http.ResponseWriter, r *http.Request) {
	//TODO fixme: create a new cluster every query, and add all engines to the cluster
	info := *c.cluster.ClusterInfo()
	//FIXME when there are no engines, the clusterinfo json is just {}
	log.Printf("ClusterInfo accessed\n")
	jsonSuccess(w, info)
}

/*
  Container actions
*/

// get containers on a specific engine
func getContainers(w http.ResponseWriter, r *http.Request) {
	all := mux.Vars(r)["all"] == "true"
	id := mux.Vars(r)["id"]
	if containers, err := c.GetContainers(id, all); err != nil {
		jsonError(w, fmt.Sprintf("Unable to list containers for %s: %s", id, err.Error()))
	} else {
		log.Printf("%s has containers: %s\n", id, containers)
		jsonSuccess(w, containers)
	}
}

func (c *Controller) GetContainers(id string, all bool) ([]*citadel.Container, error) {
	engine, err := model.Get(c.redis, id, &c.config)
	if err != nil {
		return nil, fmt.Errorf("No engine found with ID %s: %s", id, err.Error())
	}
	if err := engine.SetupClient(&c.config); err != nil {
		return nil, fmt.Errorf("Unable to setup docker client for %s: %s", id, err.Error())
	}
	return engine.Engine.ListContainers(all)
}
