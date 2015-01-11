package controller

import (
	"net/http"

	"github.com/byxorna/clementine/config"
	"github.com/gorilla/mux"
)

var (
	self *Controller
)

type Controller struct {
	conf   *config.Config
	router *mux.Router
}

func Setup(conf *config.Config) error {
	self = &Controller{
		router: mux.NewRouter(),
		conf:   conf,
	}

	self.router.HandleFunc("/", root).Methods("GET")
	self.router.HandleFunc("/apiv1/config", getConfig).Methods("GET")
	self.router.HandleFunc("/apiv1/cluster", clusterInfo).Methods("GET")
	self.router.HandleFunc("/apiv1/engine/{id}/containers", getContainers).Methods("GET")
	self.router.HandleFunc("/apiv1/engines", engines).Methods("GET")
	self.router.HandleFunc("/apiv1/engine", createEngine).Methods("POST")
	self.router.HandleFunc("/apiv1/engine/{id}", getEngine).Methods("GET")
	// schedule a container out across the cluster
	self.router.HandleFunc("/apiv1/container", runContainer).Methods("POST")
	// stop and remove a container from a specific engine
	self.router.HandleFunc("/apiv1/engine/{engine}/container/{name}", destroyContainer).Methods("DELETE")

	return nil
}

func Serve() error {
	if err := http.ListenAndServe(self.conf.ListenAddr, self.router); err != nil {
		return err
	}
	return nil
}

func root(w http.ResponseWriter, r *http.Request) {
	jsonSuccess(w, "Hello world! GET /apiv1/")
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	jsonSuccess(w, c.config)
}
