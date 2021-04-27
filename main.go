package main

import (
	"defaultscheduler/pkg/config"
	"defaultscheduler/pkg/controller"
	"defaultscheduler/pkg/simulator"
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
)

func main(){
	
	fmt.Printf("info")
	config.Test()
	flag.Parse()
	glog.Infof("test")
	controller.Test1()
	simulator.SimInit()
	//time.Sleep(100*time.Second)
}