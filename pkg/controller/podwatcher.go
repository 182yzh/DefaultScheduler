package controller

import (
	"context"

	"github.com/golang/glog"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type PodWatcher struct {
	ClientSet    kubernetes.Interface
	Controller   cache.Controller
}

// NewPodWatcher initialize a PodWatcher.
func NewPodWatcher(client kubernetes.Interface )*PodWatcher {
	glog.Infof("Starting PodWatcher...")
	podWatcher := &PodWatcher{
		ClientSet: client,
	}
	schedulerSelector := fields.Everything()
	podSelector := labels.Everything()
	_, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(alo metav1.ListOptions) (runtime.Object, error) {
				alo.FieldSelector = schedulerSelector.String()
				alo.LabelSelector = podSelector.String()
				return client.CoreV1().Pods("").List(context.TODO(),alo)
			},
			WatchFunc: func(alo metav1.ListOptions) (watch.Interface, error) {
				alo.FieldSelector = schedulerSelector.String()
				alo.LabelSelector = podSelector.String()
				return client.CoreV1().Pods("").Watch(context.TODO(),alo)
			},
		},
		&v1.Pod{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err != nil {
					glog.Errorf("AddFunc: error getting key %v", err) 
				}
				podWatcher.enqueuePodAddition(key, obj)
			},
			UpdateFunc: func(old, new interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(new)
				if err != nil {
					glog.Errorf("UpdateFunc: error getting key %v", err)
				}
				podWatcher.enqueuePodUpdate(key, old, new)
			},
			DeleteFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err != nil {
					glog.Errorf("DeleteFunc: error getting key %v", err)
				}
				podWatcher.enqueuePodDeletion(key, obj)
			},
		},
	)
	podWatcher.Controller = controller
	return podWatcher
}




func (pw *PodWatcher) enqueuePodAddition(key interface{}, obj interface{}) {
	pod := obj.(*v1.Pod)
	_ = pod.Name
	for _ , informer := range PodWatcherInformers {
		go informer.AddFunction(obj)
	}
	glog.V(106).Infof("V6 enqueuePodAddition for pod %s",pod.Name)
}

func (pw *PodWatcher) enqueuePodDeletion(key interface{}, obj interface{}) {
	pod := obj.(*v1.Pod)
	_ = pod.Name
	glog.Infof(pod.Name)
	for _ , informer := range PodWatcherInformers {
		go informer.DeleteFunction(obj)
	}
	glog.V(6).Infof("V6 enqueuePodDeletion for pod %s",pod.Name)
}

func (pw *PodWatcher) enqueuePodUpdate(key, oldObj, newObj interface{}) {
	oldPod := oldObj.(*v1.Pod)
	newPod := newObj.(*v1.Pod)
	for _ , informer := range PodWatcherInformers {
		go informer.UpdateFunction(newObj,oldObj)
	}
	glog.V(6).Infof("V6 enqueuePodUpdate, newPod Name :%s, Status: %s; oldPod Name :%s, Status :%s ",
				newPod.Name, newPod.Status.Phase,oldPod.Name,oldPod.Status.Phase)
}


