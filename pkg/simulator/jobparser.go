package simulator

import (

	//"io/ioutil"

	//"sort"

	"github.com/golang/glog"

	//"k8s.io/apimachinery/pkg/util/wait"
	//"sync"

	"time"
)

type DetailInfo struct {
    Ip string `json:"ip"`
    Gpus []string  `json:"gpus"`
}

type AttemptInfo struct{
    Start_time string `json:"start_time"`
    End_time string `json:"end_time"`
    Detail []DetailInfo `json:"detail"`
}

type JobInfo struct{
    Status string `json:"status"`
    Vc string `json:"vc"`
    Jobid string `json:"jobid"`
    Attempts []AttemptInfo   `json:"attempts"`
    Submitted_time string `json:"submitted_time"`
	User string `json:"user"`
	Test string `json:"server"`
}

func (jif *JobInfo)GetGpuNeedNum()int64{
	var ans int64 = 0
	for _,det := range jif.Attempts[len(jif.Attempts)-1].Detail{
		ans+=int64(len(det.Gpus))
	}
	return ans
}

func (jif *JobInfo)GetTaskNum()int64{
	return int64(len(jif.Attempts[len(jif.Attempts)-1].Detail))
}

func (jif *JobInfo)ExcuteTime()time.Duration{
	att1 := jif.Attempts[0]
	att2 := jif.Attempts[len(jif.Attempts)-1]
	startTime:=StringTime2Time(att1.Start_time)
	endTime := StringTime2Time(att2.End_time)
	ans := endTime.Sub(startTime)
	return ans
}



func StringTime2Time(s string)time.Time{
	t,err := time.ParseInLocation("2006-01-02 15:04:05",s,time.Local)
	if err != nil {
		glog.Fatal(err)
	}
	return t
}


type OrderedJobInfos []*JobInfo

func (a OrderedJobInfos) Len() int           { return len(a) }
func (a OrderedJobInfos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OrderedJobInfos) Less(i, j int) bool { 
	aist := StringTime2Time(a[i].Submitted_time)
	ajst := StringTime2Time(a[j].Submitted_time)
	return aist.UnixNano() < ajst.UnixNano() 
}
