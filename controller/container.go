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

	"github.com/citadel/citadel"
	"github.com/gorilla/mux"
	"github.ewr01.tumblr.net/gabe/clementine/model"
)

func destroyContainer(w http.ResponseWriter, r *http.Request) {
	engine_id := mux.Vars(r)["engine"]
	container_name := mux.Vars(r)["name"]
	var engine *model.Engine
	var target *citadel.Container

	engine, err := model.Get(c.redis, engine_id, &c.config)
	if err != nil {
		jsonError(w, err.Error())
		return
	}
	//make sure client is setup before trying to blow up the container
	engine.SetupClient(&c.config)

	containers, err := engine.Engine.ListContainers(false)
	if err != nil {
		jsonErrorStatus(w, fmt.Sprintf("Error querying containers from engine %s: %s", engine.ID, err.Error()), http.StatusInternalServerError)
	}
	for _, container := range containers {
		if container.Name == container_name || container.ID == container_name {
			//lets allow name or id to pick out what container to kill
			target = container
		}
	}

	log.Printf("Killing container %s on %s\n", target.ID, engine.ID)
	if err := engine.Engine.Kill(target, 9); err != nil {
		log.Printf("Error killing container: %s\n", err.Error())
		jsonErrorStatus(w, fmt.Sprintf("Error killing container %s: %s", target.ID, err.Error()), http.StatusInternalServerError)
		return
	}

	log.Printf("Removing container %s on %s\n", target.ID, engine.ID)
	if err := engine.Engine.Remove(target); err != nil {
		log.Printf("Error removing container: %s\n", err.Error())
		jsonErrorStatus(w, fmt.Sprintf("Error removing container %s: %s", target.ID, err.Error()), http.StatusInternalServerError)
		return
	}

	jsonSuccessStatus(w, fmt.Sprintf("Container %s destroyed on engine %s", target, engine), http.StatusAccepted)
}

func runContainer(w http.ResponseWriter, r *http.Request) {
	var image *citadel.Image
	if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
		jsonError(w, fmt.Sprintf("Unable to decode image: %s", err.Error()))
		return
	}

	log.Printf("Constructing cluster to schedule %s on...\n", image.Name)
	cluster := c.Cluster()
	engines, err := model.AllEngines(c.redis, &c.config)
	if err != nil {
		jsonError(w, fmt.Sprintf("Unable retrieve engines: %s", err.Error()))
		return
	}
	for _, e := range engines {
		e.SetupClient(&c.config) // we need to make sure the dockerclient is configured before we schedule
		if err := cluster.AddEngine(e.Engine); err != nil {
			log.Printf("Error adding engine %s to cluster: %s\n", e.ID, err.Error())
		}
	}

	log.Printf("Attempting to schedule %s\n", image.Name)
	container, err := c.cluster.Start(image, false)
	if err != nil {
		log.Printf("Unable to schedule %s: %s\n", image.Name, err.Error())
		jsonError(w, err.Error())
		return
	}
	log.Printf("Scheduled container %s on engine %s\n", container.ID, container.Engine.ID)

	jsonSuccessStatus(w, container, http.StatusCreated)
}
