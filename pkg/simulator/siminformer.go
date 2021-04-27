package simulator

import (
	"defaultscheduler/pkg/config"
	"defaultscheduler/pkg/controller"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
)

type SimInformer struct{
	m *Monitor
	//map jobname to jobnumber
	jobs map[string]int
	mutex *sync.Mutex
}

//PodWatcherRegister
func (si *SimInformer)AddFunction(obj interface{})  {
	
	pod,ok := obj.(*v1.Pod)
	if ok ==false {
		glog.Fatal("Error can not invert to v1.pod")
	}
	if pod.Namespace != config.GetNameSpace() {
		return
	}
	gpu := GetGpuNeed(pod)
	tmp,ok := pod.Labels["gangschedulenumber"]
	if ok == false {
		glog.Fatal("pod have no label gangschedulenumber")
	}
	task,err := strconv.Atoi(tmp)
	if err  != nil {
		glog.Fatalf("can not invert %s to int",tmp)
	}
	timenow := time.Now()
	formattime := timenow.Format("2006-01-02 13:04:05")
	glog.Errorf("SCHEDULE_EVENT:SUBMIT %s %d %d %s",pod.Name,task,gpu,formattime)
	return
}

func (si *SimInformer)UpdateFunction(newobj,oldobj interface{})  {
	newpod,ok := newobj.(*v1.Pod)
	if ok ==false {
		glog.Fatal("Error can not invert to v1.pod")
	}
	oldpod,ok := oldobj.(*v1.Pod)
	if ok ==false {
		glog.Fatal("Error can not invert to v1.pod")
	}
	if newpod.Status.Phase == v1.PodRunning && oldpod.Status.Phase != v1.PodRunning{
		gpu := GetGpuNeed(newpod)
		tmp,ok := newpod.Labels["gangschedulenumber"]
		if ok == false {
			glog.Fatal("pod have no label gangschedulenumber")
		}
		task,err := strconv.Atoi(tmp)
		if err  != nil {
			glog.Fatalf("can not invert %s to int",tmp)
		}
		timenow := time.Now()
		formattime := timenow.Format("2006-01-02 13:04:05")
		glog.Errorf("SCHEDULE_EVENT:PLACE %s %d %d %s",newpod.Name,task,gpu,formattime)		
	}

	return
}


func (si *SimInformer)DeleteFunction(obj interface{})  {
	print(obj)
	pod,ok := obj.(*v1.Pod)
	if ok == false{
		glog.Fatal("can not invert to pod")
	}
	jobname,ok := pod.Labels["job-name"]
	if ok == false {
		glog.Fatal("job not hava job name")
	}
	tasknumber,err := strconv.Atoi(pod.Labels["gangschedulenumber"])
	si.mutex.Lock()
	if err !=nil{
		glog.Fatal(err.Error())
	}
	if _,ok := si.jobs[jobname];ok{
		si.jobs[jobname] += 1	
	} else {
		si.jobs[jobname] = 1 
	}
	temok := (si.jobs[jobname] == tasknumber)
	si.mutex.Unlock()
	if temok == false {
		return 
	}
	gpu := GetGpuNeed(pod)
	if temok && pod.Status.Phase == v1.PodSucceeded {
		si.mutex.Lock()
		delete(si.jobs,jobname)
		timenow := time.Now()
		formattime := timenow.Format("2006-01-02 13:04:05")
		glog.Error("SCHEDULE_EVENT:COMPLETED %s %d %d %s",pod.Name,tasknumber,gpu,formattime)
		si.mutex.Unlock()
		si.m.AddGPU( int64(tasknumber) * gpu*-1)
	} 
	if temok &&  pod.Status.Phase != v1.PodSucceeded {
		if _,ok := pod.Labels["job-name"];ok==false{
			glog.Fatal("Pod do not have job name")
		}
		go ReSubmitJob(pod.Labels["job-name"])
	}
	
	return 
}

func GetGpuNeed(pod *v1.Pod)int64{
	if pod == nil {
		glog.Error("pod is nil")
		return 0
	}
	ans := int64(0)
	//bytepod,_ := json.Marshal(pod)
	//fmt.Println(string(bytepod))
	for _,con := range pod.Spec.Containers {
		gpu,ok := con.Resources.Requests["nvidia.com/gpu"]
		if ok ==false {
			glog.Fatalf("pod %s have no res is gpu",pod.Name)
		}
		num,ok := gpu.AsInt64()
		if ok == false  {
			glog.Fatal(gpu ,"cant invet to int64")
		}
		ans += num
	}
	return ans
}

func SimInit()*SimInformer{
	time.Sleep(5*time.Second)
	si := new(SimInformer)
	si.mutex = new(sync.Mutex)
	si.jobs = make(map[string]int)
	m := new(Monitor)
	m.Init("Test/cluster_job_log")
	si.m = m
	glog.Infof("start init")
	controller.PodWatcherRegister("simpodwatcher",si)
	if config.GetSimulateMode() == "gpunumber"{
		si.m.BeginGpuMonitor()
	} else if config.GetSimulateMode() == "timestep" {
		si.m.BeginTimeMonitor()
	} else {
		glog.Fatal("Error cannot see this simulator mod")
	}
	return si
}