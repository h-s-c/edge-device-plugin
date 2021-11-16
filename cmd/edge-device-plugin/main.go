package main

import (
	"flag"
	"log"
	"sync"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

var finished chan bool

func main() {
	flag.Parse()

	log.Println("Starting Edge device plugin for Kubernetes")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		tpumanager := dpm.NewManager(TPULister{})
		tpumanager.Run()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		vpumanager := dpm.NewManager(VPULister{})
		vpumanager.Run()
	}()
	wg.Wait()

	log.Println("Exiting Edge device plugin for Kubernetes")
}
