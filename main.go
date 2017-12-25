package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/na-ga/k8s-election/election"
	"github.com/robfig/cron"

	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
)

var (
	namespace    = "default"
	electionName = "k8s-election"
	ttl          = time.Second * 5
	port         = 8989
)

func main() {

	// create elector
	identity, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to create identity: %s", err.Error())
	}
	client, err := newClient()
	if err != nil {
		log.Fatalf("Failed to connecting to the client: %s", err.Error())
	}
	callbacks := &leaderelection.LeaderCallbacks{
		OnStartedLeading: func(stop <-chan struct{}) { log.Println("Started leading") },
		OnStoppedLeading: func() { log.Println("Stopped leading") },
		OnNewLeader:      func(readerIdentity string) { log.Printf("Detected leader: readerId=%s", readerIdentity) },
	}
	elector, err := election.NewElectorWithCallbacks(namespace, electionName, identity, ttl, client, callbacks)
	if err != nil {
		log.Fatalf("Failed to create election: %s", err.Error())
	}

	// start elector
	go elector.Run()

	// start task worker
	worker := cron.New()
	defer worker.Stop()
	task := func() {
		if !elector.IsLeader() {
			return // do nothing if not leader
		}
		log.Println("I am a leader")
	}
	spec := "@every " + (time.Second * 5).String()
	worker.AddFunc(spec, task)
	worker.Start()

	// start http monitor
	errCh := make(chan error)
	go func() {
		webHandler := func(res http.ResponseWriter, req *http.Request) {
			ld := map[string]interface{}{
				"id":       identity,
				"leaderId": elector.GetLeader(),
				"isLeader": elector.IsLeader(),
			}
			data, err := json.Marshal(ld)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
				return
			}
			res.WriteHeader(http.StatusOK)
			res.Write(data)
		}
		http.HandleFunc("/", webHandler)
		log.Printf("Http server is running on :%d", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Printf("Http server is terminated: %s", err.Error())
			errCh <- err
		}
		log.Printf("Http server is terminated: %s", err.Error())
	}()

	// waiting signals or errors
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		close(signalCh)
		signal.Stop(signalCh)
	}()
	select {
	case err := <-errCh:
		log.Printf("Terminated with error: %s", err.Error())
	case <-signalCh:
		log.Println("Shutdown signal is received")
	}
}

func newClient() (*v1.CoreV1Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("error create cluster config: %s", err.Error())
	}
	return v1.NewForConfig(config)
}
