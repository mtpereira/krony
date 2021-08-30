package main

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func connectToK8s() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panicf("InClusterConfig: %v", err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicf("NewForConfig: %v", err.Error())
	}
	return clientset
}

func listJobs(clientset *kubernetes.Clientset) []batchv1.Job {
	jobs, err := clientset.BatchV1().Jobs("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Panicf("main: %v", err.Error())
	}
	return jobs.Items
}

type server struct {
	Router *http.ServeMux
}

func newServer() *server {
	s := server{Router: http.NewServeMux()}
	s.Router.HandleFunc("/jobs/list", handleListJobs())
	return &s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func handleListJobs() http.HandlerFunc {
	clientset := connectToK8s()
	return func(w http.ResponseWriter, r *http.Request) {
		jobs := listJobs(clientset)
		respondWithJSON(w, http.StatusOK, jobs)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func run() error {
	s := newServer()
	log.Println("Starting API server...")
	err := http.ListenAndServe(":9090", s)
	if err != nil {
		return errors.Wrap(err, "setup http server")
	}
	log.Println("Stopping...")
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
