package main

import (
	"flag"
	"log"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

func main() {
	accelerator := flag.String("accelerator", "tpu", "tpu/vpu")
	flag.Parse()

	log.Println("Edge device plugin for Kubernetes")

	if *accelerator == "tpu" {
		tpumanager := dpm.NewManager(TPULister{})
		tpumanager.Run()
	} else if *accelerator == "vpu" {
		vpumanager := dpm.NewManager(VPULister{})
		vpumanager.Run()
	}
}
