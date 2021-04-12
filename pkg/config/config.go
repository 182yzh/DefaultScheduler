package config

import (
	"flag"
	"fmt"
)

var Kubeconfig string
var AutoMoveFinishJob bool
var NameSpace string
var TotalGPU int64
var SimulateMode string 
var GPULimit int64
var TimeStep int64
var DebugMode bool


func init(){
	flag.StringVar(&Kubeconfig, "kubeconfig", "", "the path of the kubernetes config")
	flag.BoolVar(&AutoMoveFinishJob,"automovefinishjob",true,"if the option is true,  compelete jobs are removed autoly by k8s")
	flag.StringVar(&NameSpace,"namespace","default","the scheduler managers namespace")
	flag.Int64Var(&TotalGPU,"totalgpu",1280,"the gpus that the cluster have")
	flag.StringVar(&SimulateMode,"simulatemode","gpunumber","the simulator mod,if you choose gpunumber you should also use gpulimit ans totalgpu")
	flag.Int64Var(&GPULimit,"gpulimit",64,"the gpu limit that excess")
	flag.Int64Var(&TimeStep,"timestep",50,"the time step that timestep mode need")
	flag.BoolVar(&DebugMode,"debugmode",false,"debug log more info")
}

func IsDebugMode()bool{
	return DebugMode
}

func GetTimeStep()int64{
	return TimeStep
}
func GetGPULimit()int64{
	return GPULimit
}

func GetSimulateMode()string{
	return SimulateMode
}


func Test(){
	fmt.Println("test in config")
}

func GetTotalGPU()int64{
	return TotalGPU
}

func GetNameSpace()string {
	return NameSpace
}