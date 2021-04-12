package controller

import (
	"context"

	"github.com/golang/glog"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type NodeWatcher struct {
	ClientSet    kubernetes.Interface
	Controller   cache.Controller
}




// NewNodeWatcher initialize a NodeWatcher.
func NewNodeWatcher(client kubernetes.Interface )*NodeWatcher {
	glog.Infof("Starting NodeWatcher...")
	nodeWatcher := &NodeWatcher{
		ClientSet: client,
	}
	_, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(alo metav1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().Nodes().List(context.TODO(),alo)
			},
			WatchFunc: func(alo metav1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Nodes().Watch(context.TODO(),alo)
			},
		},
		&v1.Node{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err != nil {
					glog.Errorf("AddFunc: error getting key %v", err)
				}
				nodeWatcher.enqueueNodeAddition(key, obj)
			},
			UpdateFunc: func(old, new interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(new)
				//glog.Infof("SCHESULE UpdateFunc,key:%v",key)
				if err != nil {
					glog.Errorf("UpdateFunc: error getting key %v", err)
				}
				nodeWatcher.enqueueNodeUpdate(key, old, new)
			},
			DeleteFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err != nil {
					glog.Errorf("DeleteFunc: error getting key %v", err)
				}
				nodeWatcher.enqueueNodeDeletion(key, obj)
			},
		},
	)
	nodeWatcher.Controller = controller
	return nodeWatcher
}




func (pw *NodeWatcher) enqueueNodeAddition(key interface{}, obj interface{}) {
	pod := obj.(*v1.Node)
	_ = pod.Name
	for _ , informer := range NodeWatcherInformers {
		go informer.AddFunction(obj)
	}
	glog.V(6).Infof("------V6 enqueueNodeAddition for pod %s",pod.Name)
}

func (pw *NodeWatcher) enqueueNodeDeletion(key interface{}, obj interface{}) {
	pod := obj.(*v1.Node)
	_ = pod.Name
	glog.Infof(pod.Name)
	for _ , informer := range NodeWatcherInformers {
		go informer.DeleteFunction(obj)
	}
	glog.V(6).Infof("V6 enqueueNodeDeletion for pod %s",pod.Name)
}

func (pw *NodeWatcher) enqueueNodeUpdate(key, oldObj, newObj interface{}) {
	oldNode := oldObj.(*v1.Node)
	newNode := newObj.(*v1.Node)
	for _ , informer := range NodeWatcherInformers {
		go informer.UpdateFunction(newObj,oldObj)
	}
	glog.V(102).Infof("V6 enqueueNodeUpdate, newPod Name :%s, Status: %s; oldPod Name :%s, Status :%s ",
				newNode.Name, newNode.Status.Phase,oldNode.Name,oldNode.Status.Phase)
}


