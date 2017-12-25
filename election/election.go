package election

import (
	"log"
	"os"
	"time"

	api "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	cli "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

//
const componentName = "k8s-leader-elector"

// NewElector returns leader elector that using config maps resource lock.
func NewElector(namespace, configMapName, identity string, ttl time.Duration, client cli.CoreV1Interface) (*leaderelection.LeaderElector, error) {
	return NewElectorWithCallbacks(namespace, configMapName, identity, ttl, client, nil)
}

// NewElectorWithCallbacks returns leader elector that using config maps resource lock.
func NewElectorWithCallbacks(namespace, configMapName, identity string, ttl time.Duration, client cli.CoreV1Interface, callbacks *leaderelection.LeaderCallbacks) (*leaderelection.LeaderElector, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	broadcaster := record.NewBroadcaster()
	broadcaster.StartLogging(log.Printf)
	broadcaster.StartRecordingToSink(&cli.EventSinkImpl{Interface: client.Events(namespace)})
	recorder := broadcaster.NewRecorder(scheme.Scheme, api.EventSource{Component: componentName, Host: hostname})
	cmLock := &resourcelock.ConfigMapLock{
		Client: client,
		ConfigMapMeta: meta.ObjectMeta{
			Namespace: namespace,
			Name:      configMapName,
		},
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      identity,
			EventRecorder: recorder,
		},
	}
	if callbacks == nil {
		callbacks = NewDefaultCallbacks()
	}
	config := leaderelection.LeaderElectionConfig{
		Lock:          cmLock,
		LeaseDuration: ttl,
		RenewDeadline: ttl / 2,
		RetryPeriod:   ttl / 4,
		Callbacks:     *callbacks,
	}
	return leaderelection.NewLeaderElector(config)
}

// NewDefaultCallbacks returns default leader callbacks.
func NewDefaultCallbacks() *leaderelection.LeaderCallbacks {
	return &leaderelection.LeaderCallbacks{
		OnStartedLeading: func(stop <-chan struct{}) {},
		OnStoppedLeading: func() {},
		OnNewLeader:      func(identity string) {},
	}
}
