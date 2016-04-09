/*
Copyright 2014 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sensoraccess

import (
	"net/http"
	"time"

	"github.com/golang/glog"

	//"fmt"
	//"net"
	//"os"
	//"strconv"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/cache"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/quota"
	"k8s.io/kubernetes/pkg/runtime"
	utilruntime "k8s.io/kubernetes/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/util/wait"
	"k8s.io/kubernetes/pkg/util/workqueue"
	"k8s.io/kubernetes/pkg/watch"
)

// SensorAccessControllerOptions holds options for creating a quota controller
type SensorAccessControllerOptions struct {
	// Must have authority to list all quotas, and update quota status
	KubeClient clientset.Interface
	// Controls full recalculation of quota usage
	ResyncPeriod controller.ResyncPeriodFunc
	// Knows how to calculate usage
	Registry quota.Registry
	// List of GroupKind objects that should be monitored for replenishment at
	// a faster frequency than the quota controller recalculation interval
	GroupKindsToReplenish []unversioned.GroupKind
}

// SensorAccessController is responsible for tracking quota usage status in the system
type SensorAccessController struct {
	// Must have authority to list all resources in the system, and update quota status
	kubeClient clientset.Interface
	// An index of resource quota objects by namespace
	rqIndexer cache.Indexer
	// Watches changes to all resource quota
	rqController *framework.Controller
	// SensorAccess objects that need to be synchronized
	queue *workqueue.Type
	// To allow injection of syncUsage for testing.
	syncHandler func(key string) error
	// function that controls full recalculation of quota usage
	resyncPeriod controller.ResyncPeriodFunc
	// knows how to calculate usage
	registry quota.Registry

	//conn *net.TCPConn
}

func NewSensorAccessController(options *SensorAccessControllerOptions) *SensorAccessController {
	// build the resource quota controller
	rq := &SensorAccessController{
		kubeClient:   options.KubeClient,
		queue:        workqueue.New(),
		resyncPeriod: options.ResyncPeriod,
		registry:     options.Registry,
	}

	//tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:12892")

	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	//}

	//qnon, err := net.DialTCP("tcp", nil, tcpAddr)

	//rq.conn = qnon

	// set the synchronization handler
	rq.syncHandler = rq.syncSensorAccessFromKey

	// build the controller that observes quota
	rq.rqIndexer, rq.rqController = framework.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (runtime.Object, error) {
				return rq.kubeClient.Core().SensorAccesses(api.NamespaceAll).List(options)
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return rq.kubeClient.Core().SensorAccesses(api.NamespaceAll).Watch(options)
			},
		},
		&api.SensorAccess{},
		rq.resyncPeriod(),
		framework.ResourceEventHandlerFuncs{
			AddFunc: rq.enqueueSensorAccess,
			UpdateFunc: func(old, cur interface{}) {
				// We are only interested in observing updates to quota.spec to drive updates to quota.status.
				// We ignore all updates to quota.Status because they are all driven by this controller.
				// IMPORTANT:
				// We do not use this function to queue up a full quota recalculation.  To do so, would require
				// us to enqueue all quota.Status updates, and since quota.Status updates involve additional queries
				// that cannot be backed by a cache and result in a full query of a namespace's content, we do not
				// want to pay the price on spurious status updates.  As a result, we have a separate routine that is
				// responsible for enqueue of all resource quotas when doing a full resync (enqueueAll)
				oldSensorAccess := old.(*api.SensorAccess)
				curSensorAccess := cur.(*api.SensorAccess)
				if curSensorAccess.Access == oldSensorAccess.Access {
					return
				}
				rq.enqueueSensorAccess(curSensorAccess)
			},
			// This will enter the sync loop and no-op, because the controller has been deleted from the store.
			// Note that deleting a controller immediately after scaling it to 0 will not work. The recommended
			// way of achieving this is by performing a `stop` operation on the controller.
			DeleteFunc: rq.enqueueSensorAccess,
		},
		cache.Indexers{"namespace": cache.MetaNamespaceIndexFunc},
	)

	return rq
}

// enqueueAll is called at the fullResyncPeriod interval to force a full recalculation of quota usage statistics
func (rq *SensorAccessController) enqueueAll() {
	defer glog.V(4).Infof("Resource quota controller queued all resource quota for full calculation of usage")
	for _, k := range rq.rqIndexer.ListKeys() {
		rq.queue.Add(k)
	}
}

// obj could be an *api.SensorAccess, or a DeletionFinalStateUnknown marker item.
func (rq *SensorAccessController) enqueueSensorAccess(obj interface{}) {
	key, err := controller.KeyFunc(obj)
	if err != nil {
		glog.Errorf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	rq.queue.Add(key)
}

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (rq *SensorAccessController) worker() {
	for {
		func() {
			key, quit := rq.queue.Get()
			if quit {
				return
			}
			defer rq.queue.Done(key)
			err := rq.syncHandler(key.(string))
			if err != nil {
				utilruntime.HandleError(err)
				rq.queue.Add(key)
			}
		}()
	}
}

// Run begins quota controller using the specified number of workers
func (rq *SensorAccessController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	go rq.rqController.Run(stopCh)
	// the workers that chug through the quota calculation backlog
	for i := 0; i < workers; i++ {
		go wait.Until(rq.worker, time.Second, stopCh)
	}
	// the timer for how often we do a full recalculation across all quotas
	go wait.Until(func() { rq.enqueueAll() }, rq.resyncPeriod(), stopCh)
	<-stopCh
	glog.Infof("Shutting down SensorAccessController")
	rq.queue.ShutDown()
}

// syncSensorAccessFromKey syncs a quota key
func (rq *SensorAccessController) syncSensorAccessFromKey(key string) (err error) {
	startTime := time.Now()
	defer func() {
		glog.V(4).Infof("Finished syncing resource quota %q (%v)", key, time.Now().Sub(startTime))
	}()

	obj, exists, err := rq.rqIndexer.GetByKey(key)
	if !exists {
		glog.Infof("Resource quota has been deleted %v", key)
		return nil
	}
	if err != nil {
		glog.Infof("Unable to retrieve resource quota %v from store: %v", key, err)
		rq.queue.Add(key)
		return err
	}
	quota := *obj.(*api.SensorAccess)
	return rq.syncSensorAccess(quota)
}

// syncSensorAccess runs a complete sync of resource quota status across all known kinds
func (rq *SensorAccessController) syncSensorAccess(resourceQuota api.SensorAccess) (err error) {

	//_, err = rq.conn.Write([]byte(strconv.Itoa(int(resourceQuota.Access))))
	_, err = http.Get("http://localhost:12892")
	_, err = rq.kubeClient.Core().SensorAccesses(resourceQuota.Namespace).Update(&resourceQuota)

	return err
}
