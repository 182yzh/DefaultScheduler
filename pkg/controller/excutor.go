package controller

// excuter is a mod for do some operations to k8s, such as lebels pod, delete jobs and others.

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)




func CS_PatchLabelToPod(ClientSet kubernetes.Interface,namespace,podname string, patchData map[string]interface{})bool{
    playLoadBytes, _ := json.Marshal(patchData)
    _, err := ClientSet.CoreV1().Pods(namespace).Patch(context.TODO(),podname, types.StrategicMergePatchType, playLoadBytes,metav1.PatchOptions{})
    if err != nil {
        glog.Errorf("[UpdatePodByPodSn] %v pod Patch fail %v", podname, err)
        return false
	}
	glog.V(1).Infof("success patch label %s to pod %s",string(playLoadBytes),podname)
    return true
}


func CS_BindPodToNode(ClientSet kubernetes.Interface, Namespace,Podname,Nodename string) (bool,error){
	//clusterinfo.WaitForResReadyLocal(Namespace,Podname,Nodename)
	err := ClientSet.CoreV1().Pods(Namespace).Bind(context.TODO(),&v1.Binding{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: Podname,
		},
		Target: v1.ObjectReference{
			Namespace: Namespace,
			Name:      Nodename,
		}},
		metav1.CreateOptions{})
	if err != nil {
		glog.Errorf("Could not bind pod:%s to nodeName:%s, error: %v",Podname, Nodename, err)
		return false,err	
	} else {
		glog.Infof("Bind pod to node : %s to %s",Podname, Nodename)
		return true,nil
	}
}


func CS_RemovePod(ClientSet kubernetes.Interface,Podname,Namespace,Jobname string) (bool,error){
	var graceTime int64 = 0
	if Jobname != Podname[:len(Podname)-6]{
		glog.Fatal("Fatal Error, this task and job is not mapping")
	}
	propagationPolicy := metav1.DeletePropagationBackground
	err := ClientSet.BatchV1().Jobs(Namespace).Delete(context.TODO(), Jobname, metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,GracePeriodSeconds : &graceTime})
	if err != nil {
		glog.Errorf("delete jobs error: %v")
		return false,err
	}
	err = ClientSet.CoreV1().Pods(Namespace).Delete(context.TODO(),Podname, metav1.DeleteOptions{PropagationPolicy: &propagationPolicy,GracePeriodSeconds : &graceTime})
	if err != nil {
		glog.Errorf("Could not delete pod:%s in namespace:%s, error: %v", Podname, "default", err)
		return false,err
	}
	return true,nil
}

func CS_GetAllPods(Clientset kubernetes.Interface)(*v1.PodList,error){
	pods,err := ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	return pods,err
}

func GetAllPods()(*v1.PodList,error){
	return CS_GetAllPods(ClientSet)
}

func CS_GetAllNodes(Clientset kubernetes.Interface)(*v1.NodeList,error){
	nodes,err := ClientSet.CoreV1().Nodes().List(context.TODO(),metav1.ListOptions{})
	return nodes,err
}

func PatchLabelToPod(Namespace,Podname string, patchData map[string]interface{})bool{
    return CS_PatchLabelToPod(ClientSet, Namespace, Podname, patchData)
}

func BindPodToNode(Namespace,Podname,Nodename string) (bool,error){
	return CS_BindPodToNode(ClientSet,Namespace,Podname,Nodename)
}

func RemovePod(Podname,Namespace,Jobname string) (bool,error){
	pod,err  := ClientSet.CoreV1().Pods(Namespace).Get(context.TODO(),Podname,metav1.GetOptions{})
	if err != nil {
		glog.Error(err)
		glog.Error("Can not find the pod for " + Podname)
		return false, err
	}
	strtasknum,ok := pod.Labels["gangschedulenumber"]
	if ok == false {
		glog.Fatal("Node do not have sgangschedulenumber label, "+Podname)
	}
	tasknum,err := strconv.Atoi(strtasknum)
	if err !=nil {
		glog.Fatal("Node do not have sgangschedulenumber label, "+Podname)
	}
	if tasknum > 1 {
		glog.Error("should not preempt gang schedule task")
		return false,nil
	}
	return CS_RemovePod(ClientSet,Podname,Namespace,Jobname)
}

func GetAllNodes()(*v1.NodeList,error){
	return CS_GetAllNodes(ClientSet)
}


