package controller

import (
	"context"
	"defaultscheduler/pkg/config"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//	"k8s.io/client-go/tools/clientcmd"
)


var ClientSet kubernetes.Interface
var podWatcher *PodWatcher
var nodeWatcher *NodeWatcher


func init(){
	flag.Parse()
	kubeconfig := config.Kubeconfig
	fmt.Println(kubeconfig)
	ClientSet = BuildNewClientSet(kubeconfig)
	podWatcher = NewPodWatcher(ClientSet)
	nodeWatcher = NewNodeWatcher(ClientSet)
	if PodWatcherInformers == nil {
		PodWatcherInformers = make(map[string]PodWatcherInformer)
	}
	if NodeWatcherInformers == nil {
		NodeWatcherInformers = make(map[string]NodeWatcherInformer)
	}
	stopch := make(chan struct{})
	go nodeWatcher.Controller.Run(stopch)
	go podWatcher.Controller.Run(stopch)
}


func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}


// Build a new kuberbetes clientset 
func BuildNewClientSet(kubeConfig string)kubernetes.Interface{
	config, err := GetClientConfig(kubeConfig)
	if err != nil {
		glog.Fatalf("Failed to load client config: %v", err)
	}
	ClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create connection: %v", err)
	}
	glog.Infof("Success build a kubernetes client using the kuebeconfig <%s>\n",kubeConfig)
	return ClientSet
}




func Test1() {
	
	
	fmt.Println("test in controller")
	return 
	glog.Fatal("errr")
	stopch := make(chan struct{})
	go podWatcher.Controller.Run(stopch)
	go nodeWatcher.Controller.Run(stopch)
	time.Sleep(3*time.Second)
	//stopch <- struct{}{}
	fmt.Println("exit the test")

	//node,_ := ClientSet.CoreV1().Nodes().Get(context.TODO(),"giga-node1",metav1.GetOptions{})
	//fmt.Println(node.GetLabels())
}


func GetPod(namespace,podname string ){
	pod,err := ClientSet.CoreV1().Pods(namespace).Get(context.TODO(),podname,metav1.GetOptions{})
	if err != nil {
		glog.Fatal(err)
	}
	j,_ := json.Marshal(pod)
	fmt.Println(string(j))
}