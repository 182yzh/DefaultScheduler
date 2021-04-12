package controller

import (
	//"fmt"
	"github.com/golang/glog"
)

type PodWatcherInformer interface {
	AddFunction(obj interface{})
	UpdateFunction(newobj,oldobj interface{})
	DeleteFunction(obj interface{})
}


type NodeWatcherInformer interface {
	AddFunction(obj interface{})
	UpdateFunction(newobj,oldobj interface{})
	DeleteFunction(obj interface{})
}

var PodWatcherInformers map[string]PodWatcherInformer
var NodeWatcherInformers map[string]NodeWatcherInformer





// Register add a new informer to the process
func NodeWatcherRegister(informerID string, informer PodWatcherInformer) (bool, string) {
	if _, haskey := NodeWatcherInformers[informerID]; haskey {
		return false, "the informerID <" + informerID + "> has been registered\n"
	}
	NodeWatcherInformers[informerID] = informer
	glog.Infof("Register informers for informerID <%s>", informerID)
	nodes,err := GetAllNodes()
	if err != nil{
		glog.Fatal("nodes is err",err.Error())
	}
	for  i,_ := range nodes.Items {
		informer.AddFunction(&(nodes.Items[i]))
	}
	
	return true, ""
}

// Revoke the informer for informerID
func NodeWatcherRevoke(informerID string) (bool, string) {
	if _, haskey := NodeWatcherInformers[informerID]; !haskey {
		return false, "The informerID <" + informerID + "> was not registered\n"
	}
	delete(NodeWatcherInformers, informerID)
	glog.Infof("Revoke informers for informerID <%s>", informerID)
	return true, ""
}


// Register add a new informer to the process
func PodWatcherRegister(informerID string, informer PodWatcherInformer) (bool, string) {
	if _, haskey := PodWatcherInformers[informerID]; haskey {
		return false, "the informerID <" + informerID + "> has been registered\n"
	}
	PodWatcherInformers[informerID] = informer
	glog.Infof("Register informers for informerID <%s>", informerID)
	pods,err:= GetAllPods()
	if err != nil{
		glog.Fatal(err)
	}
	for i,_ := range pods.Items{
		informer.AddFunction(&(pods.Items[i]))
	}
	return true, ""
}

// Revoke the informer for informerID
func PodWatcherRevoke(informerID string) (bool, string) {
	if _, haskey := PodWatcherInformers[informerID]; !haskey {
		return false, "The informerID <" + informerID + "> was not registered\n"
	}
	delete(PodWatcherInformers, informerID)
	glog.Infof("Revoke informers for informerID <%s>", informerID)
	return true, ""
}
