package simulator

import (
	"defaultscheduler/pkg/config"
	"flag"
	"sort"
	"strconv"
	"sync"

	"bytes"
	//"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	//"sort"
	"github.com/golang/glog"
	//"k8s.io/apimachinery/pkg/util/wait"
	//"sync"

	"time"
)

type Monitor struct{
	queue OrderedJobInfos
	total int64
	inputfile *os.File
	TimeChangeRate int64
	CurGPU int64
	mutex *sync.Mutex
}
func (m *Monitor)AddGPU(x int64){
	m.mutex.Lock()
	glog.Infof("gpu is %d,x is %d",m.CurGPU,x)
	m.CurGPU += x
	if m.CurGPU < 0 {
		glog.Fatal("Error,gpu number is less than 0")
	}
	m.mutex.Unlock()
}


func (m *Monitor)ParseFile(filepath string){
	m.queue = OrderedJobInfos{}
	times := make([]int64,0,0)
	ifile,err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	m.inputfile = ifile
	dec := json.NewDecoder(m.inputfile)
	_,err = dec.Token()
	if err != nil {
		fmt.Println(err)
	}
	
	count := 0;
	for ;dec.More();{
		count +=1;
		jinfo := new(JobInfo)
		if err = dec.Decode(jinfo); err != nil {
			fmt.Println(count)
			panic(err)
		}
		if len(jinfo.Attempts) == 0{
			continue
		}
		att1 := jinfo.Attempts[0]
		att2 := jinfo.Attempts[len(jinfo.Attempts)-1]
		//att1 := att2
		startTime := att1.Start_time
		endTime  := att2.End_time
		submitTime := jinfo.Submitted_time
		if startTime == "None" || endTime == "None" || submitTime == "None" {
			continue
		}
		gpunum := jinfo.GetGpuNeedNum()
		if gpunum == 0 {
			continue
		}
		times = append(times,int64(jinfo.ExcuteTime().Seconds()))
		m.queue = append(m.queue,jinfo)
	}
	fmt.Println(m.queue.Len(),count);
	
	var mx,sum,mn int64 = 0,0,111111111
	for _,v := range times{
		sum += v
		if mx < v {mx=v}
		if mn > v {mn=v}
	}
	fmt.Println("Info: average Time,max Time,min Time",sum/int64(len(times)),mx,mn)
}

func (m *Monitor)Init(filepath string){
	m.total = config.GetTotalGPU()
	m.CurGPU = 0
	m.TimeChangeRate = 1000000000
	m.ParseFile(filepath)
	m.mutex = new(sync.Mutex)
	sort.Sort(m.queue)
}



// GPU numbers
func (m *Monitor)BeginGpuMonitor(){
	order := 1
	cnt := 0
	len := m.queue.Len()
	//m.AddGPU(10)
	glog.Error(m.queue.Len())
	for {
		//time.Sleep(time.Second*1)
		if cnt >= len {
			return 
		}

		m.mutex.Lock()
		gpu := m.CurGPU
		m.mutex.Unlock()
		lim := config.TotalGPU
		time.Sleep(20*time.Millisecond)
		if gpu >= lim {
			time.Sleep(1*time.Second)
		}
		glog.Infof("gpu is %d",gpu)
		if gpu < m.total + lim && cnt < len{
			jinfo := m.queue[cnt]
			task := jinfo.GetTaskNum()
			gpunum := jinfo.GetGpuNeedNum()/task
			name := fmt.Sprintf("default-app%d-%d-%d",order,task,gpunum)
			etime := int64(jinfo.ExcuteTime().Seconds())
			etime /= 40000
			//3600s 之上 
			if etime > 3600*10 {
				etime = etime/20
				if etime > 3600{
					etime = 3600;
				}
			} else if etime > 3600 { 
				etime = 1500 + int64((etime-3600)/18)
			} else if etime > 1800 {
				etime = 900 + int64((etime-1800)/3)
			} else if etime > 600{
				etime = 500 +int64((etime-600)/3)
			} else if etime > 300 {
				etime = 300 + int64((etime-300))
			} else {
				etime += 120
				if etime > 300{
					etime  = 300
				}
			}

			order ++
			glog.Infof("%s %d %d %d\n",name,gpunum,task,etime)
			cmd := exec.Command("./Test/create_file.sh",name,fmt.Sprintf("%d",gpunum),fmt.Sprintf("%d",task),fmt.Sprintf("%d",etime))
			err := cmd.Run()
			if err != nil{
				fmt.Printf("jobname %s gpunum %d task number %d excute time %d\n",name,gpunum,task,etime);
				glog.Fatalf(err.Error())
				panic(err)
			}
			go  cmd.Wait()
			cmd = exec.Command("kubectl","create","-f",fmt.Sprintf("yaml/"+name+".yaml"),"--kubeconfig="+config.Kubeconfig)
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil{
				//fmt.Println(stdout.String(),stderr.String())
				glog.Errorf(stdout.String()+"/"+stderr.String())
				fmt.Println(err)
				glog.Fatalf(err.Error())
			}
			go cmd.Wait()
			//time.Sleep(20 * time.Millisecond)
			m.AddGPU(task*gpunum)
			//	fmt.Println(gpu)
			glog.Infof("create job %s",name)
			cnt ++
		}
	}
}



/*
func Initmain(){
	m := new(Monitor)
	m.Init("/home/yzh/SchedulerFrame/SimulatorData/cluster_job_log")
	glog.Infof("start init")
	//controller.PodWatcherRegister("simpodwatcher",si)
	if config.GetSimulateMode() == "gpunumber"{
		m.BeginGpuMonitor()
	} else if config.GetSchedulerName() == "timestep" {
		m.BeginTimeMonitor()
	} else {
		glog.Fatal("Error cannot see this simulator mod")
	}
	return 
}
*/

func Oldmain(){

	fmt.Println("this is  a test ")
	flag.Parse()
	m := new(Monitor)
	m.Init("../tmp-test/joblog6.log")
	//m.Init("../tmp-test/JobLog100.log")
	//go m.CheckCompleted()
	//m.CheckCompleted()
	//comexit := make(chan int)
	stopchan := make(chan struct{})
	//go Execute_main(stopchan)
	tmp := make(chan int)
	go m.BeginGpuMonitor()
	time.Sleep(5 * time.Second)
	//go m.GetAllPodInfo(tmp)
	
	
	//go m.CheckCompleted(comexit)
	//m.BeginMonitor()
	<-tmp
	stopchan<-struct{}{}
	//<-comexit
//	m.GetNodeInfo("yzh-ubuntu-node-2","yzh-ubuntu-node-80")
}


func ReSubmitJob(jobname string){
	time.Sleep(3*time.Second)
	cmd := exec.Command("kubectl","create","-f",fmt.Sprintf("../yaml/"+jobname+".yaml"),"--kubeconfig="+config.Kubeconfig)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil{
		fmt.Println(stdout.String(),stderr.String())
		glog.Errorf(stdout.String()+"/"+stderr.String())
		fmt.Println(err)
		glog.Fatalf(err.Error())
	}
}


func (m *Monitor) BeginTimeMonitor(){
	order := 1
	cnt := 0
	curtime := StringTime2Time("2017-10-10 00:00:00")
	//fmt.Println(timerate,"++++",string(timerate))
	secondtime,err := time.ParseDuration(strconv.FormatInt(config.GetTimeStep(),10)+"s")
	if (err != nil) {
		glog.Infof(err.Error())
		glog.Fatalf(err.Error())
	}
	t := time.NewTicker(1*time.Second)
	defer t.Stop()
	for {
		<- t.C
		cnt += 1;
		if cnt >= m.queue.Len() {
			break
		}
		curtime = curtime.Add(secondtime)
		
		for {
			if cnt >= m.queue.Len() {
				break
			}
			jinfo := m.queue[cnt]
			submittime := StringTime2Time( jinfo.Submitted_time)
			
			if submittime.Sub(curtime) > 0 {
				break
			} 

			task := jinfo.GetTaskNum()
			gpunum := jinfo.GetGpuNeedNum()/task
			name := fmt.Sprintf("default-app%d-%d-%d",order,task,gpunum)
			etime := int64(jinfo.ExcuteTime().Seconds())
			if etime < 25{
				etime = 20
			}
			order ++
			glog.Infof("%s %d %d %d\n",name,gpunum,task,etime)
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd := exec.Command("./Test/create_file.sh",name,fmt.Sprintf("%d",gpunum),fmt.Sprintf("%d",task),fmt.Sprintf("%d",etime))
			//cmd := exec.Command("../tmp-test/create_file.sh",name,fmt.Sprintf("%d",gpunum),fmt.Sprintf("%d",task),fmt.Sprintf("%d",etime))
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil{
				fmt.Printf("jobname %s gpunum %d task number %d excute time %d\n",name,gpunum,task,etime);	
				glog.Fatalf(stdout.String()+"/"+stderr.String()+"/"+err.Error())
				panic(err)
			}
			cmd.Wait()
			cmd = exec.Command("kubectl","create","-f",fmt.Sprintf("../yaml/"+name+".yaml"),"--kubeconfig="+config.Kubeconfig)
			//cmd = exec.Command("kubectl","create","-f",fmt.Sprintf("../tmp-test/"+name+".yaml"))
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			_ = stdout.String()
			if err != nil{
				glog.Errorf("name is %s",name)
				glog.Errorf(stdout.String()+"/"+stderr.String())
				fmt.Println(stdout.String()+"/"+stderr.String())
				glog.Fatalf(stdout.String()+"/"+stderr.String())
			}
			cmd.Wait()
			glog.Infof("create job %s",name)
			cnt ++
		}
	}
}
