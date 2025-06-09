package main

import (
	"flag"
	"log"
	"sync"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

func main() {
	flag.Parse()

	log.Println("Edge device plugin for Kubernetes")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		coralmanager := dpm.NewManager(CoralLister{})
		coralmanager.Run()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		intelmanager := dpm.NewManager(IntelLister{})
		intelmanager.Run()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		raspberrypimanager := dpm.NewManager(RasberrypiLister{})
		raspberrypimanager.Run()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		sonoffmanager := dpm.NewManager(SonoffLister{})
		sonoffmanager.Run()
	}()
	wg.Wait()
}
