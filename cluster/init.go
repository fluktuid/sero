package cluster

import (
	"os"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var namespace = getNamespace()
var clientSet = initClientSet()

func initClientSet() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic().
			Err(err).
			Msg("Fatal error setting up cluster config")
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func getNamespace() string {
	ns, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		log.Panic().
			Err(err).
			Msg("Fatal error getting namespace.\nAm I running in-cluster?")
	}
	return string(ns[:])
}
