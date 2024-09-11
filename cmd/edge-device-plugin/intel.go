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

type IntelDevicePlugin struct {
	name string
}

func (dp *IntelDevicePlugin) Start() error {
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
				// Only one per host supported
				devices = append(devices, "/dev/bus/usb")
				break
			}
		}
	}

	return devices
}

func (dp *IntelDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	for {
		devs := []*pluginapi.Device{}
		for _, id := range FindVPUs() {
			dev := &pluginapi.Device{
				ID:     id,
				Health: pluginapi.Healthy,
			}
			devs = append(devs, dev)
		}
		s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})

		time.Sleep(5 * time.Second)
	}
	return nil
}

func (dp *IntelDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	responses := pluginapi.AllocateResponse{}
	for _, req := range r.ContainerRequests {
		response := pluginapi.ContainerAllocateResponse{}
		for _, id := range req.DevicesIDs {
			log.Println("Allocating device: ", id)
			dev := &pluginapi.DeviceSpec{
				HostPath:      id,
				ContainerPath: id,
				Permissions:   "rw",
			}
			response.Devices = append(response.Devices, dev)
		}
		responses.ContainerResponses = append(responses.ContainerResponses, &response)
	}

	return &responses, nil
}

func (IntelDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{PreStartRequired: false, GetPreferredAllocationAvailable: false}, nil
}

func (IntelDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}

func (dp *IntelDevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return nil, nil
}

type IntelLister struct {
}

func (l IntelLister) GetResourceNamespace() string {
	return "intel.com"
}

func (l IntelLister) Discover(pluginListCh chan dpm.PluginNameList) {
	plugins := make(dpm.PluginNameList, 0)
	plugins = append(plugins, "vpu")
	pluginListCh <- plugins
}

func (l IntelLister) NewPlugin(name string) dpm.PluginInterface {
	return &IntelDevicePlugin{
		name: name,
	}
}
