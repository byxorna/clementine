package main

import (
	"flag"
	"log"

	"github.com/byxorna/clementine/config"
	"github.com/byxorna/clementine/controller"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "conf", "", "config file")
	flag.Parse()
}

func main() {
	log.Printf("Loading config\n")
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Setting up clementine\n")
	if err := controller.Setup(*config); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("HTTP server listening on %s\n", config.ListenAddr)
	//todo: change controller.Serve(). This is nonsense
	if err := controller.Serve(); err != nil {
		log.Printf("There was an error running clementine\n")
		log.Fatal(err)
	}
}
