package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type VPUDevicePlugin struct {
	name string
	stop chan bool
}

func (dp *VPUDevicePlugin) Start() error {
	dp.stop = make(chan bool)
	dp.stop <- false
	return nil
}

func (dp *VPUDevicePlugin) Stop() error {
	dp.stop <- true
	return nil
}

func FindVPUs() []string {
	devices := []string{}
	// USB VPUs
	usbdevices, _ := filepath.Glob("/sys/bus/usb/devices/*/idVendor")
	for _, path := range usbdevices {
		vendorid, _ := os.ReadFile(path)
		if strings.Contains(string(vendorid), "03e7") {
			productid, _ := os.ReadFile(filepath.Dir(path) + "/idProduct")
			if strings.Contains(string(productid), "2485") || strings.Contains(string(productid), "f63b") {
				// Only one NCS2 per host supported
				devices = append(devices, "ncs2")
			}
		}
	}

	return devices
}

func (dp *VPUDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	for {
		if <-dp.stop {
			break
		}

		devs := []*pluginapi.Device{}
		for _, id := range FindTPUs() {
			dev := &pluginapi.Device{
				ID:     id,
				Health: pluginapi.Healthy,
			}
			devs = append(devs, dev)
		}
		s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})

		time.Sleep(5 * time.Second)
		log.Println("Loop 2")
	}
	return nil
}

func (dp *VPUDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	responses := pluginapi.AllocateResponse{}
	for _, req := range r.ContainerRequests {
		response := pluginapi.ContainerAllocateResponse{}
		for _, id := range req.DevicesIDs {
			if id == "ncs2" {
				log.Println("Allocating device: ", id)
				dev := &pluginapi.DeviceSpec{
					HostPath:      "/dev/bus/usb",
					ContainerPath: "/dev/bus/usb",
					Permissions:   "rw",
				}
				response.Devices = append(response.Devices, dev)
			}
		}
		responses.ContainerResponses = append(responses.ContainerResponses, &response)
	}

	return &responses, nil
}

func (VPUDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{PreStartRequired: false, GetPreferredAllocationAvailable: false}, nil
}

func (VPUDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}

func (dp *VPUDevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return nil, nil
}

type VPULister struct {
}

func (l VPULister) GetResourceNamespace() string {
	return "intel.com"
}

func (l VPULister) Discover(pluginListCh chan dpm.PluginNameList) {
	plugins := make(dpm.PluginNameList, 0)
	plugins = append(plugins, "vpu")
	pluginListCh <- plugins
}

func (l VPULister) NewPlugin(name string) dpm.PluginInterface {
	return &VPUDevicePlugin{
		name: name,
	}
}
