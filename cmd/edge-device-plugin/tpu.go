package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type TPUDevicePlugin struct {
	name string
}

func (dp *TPUDevicePlugin) Start() error {
	return nil
}

func FindTPUs() []string {
	devices := []string{}
	// M.2/mPCIe TPUs
	pciedevices, _ := filepath.Glob("/sys/class/apex/apex*")
	for _, path := range pciedevices {
		devices = append(devices, "/dev/"+filepath.Base(path))
	}

	// USB TPUs
	usbdevices, _ := filepath.Glob("/sys/bus/usb/devices/*/idVendor")
	for _, path := range usbdevices {
		vendorid, _ := os.ReadFile(path)
		if strings.Contains(string(vendorid), "18d1") || strings.Contains(string(vendorid), "1a6e") {
			productid, _ := os.ReadFile(filepath.Dir(path) + "/idProduct")
			if strings.Contains(string(productid), "9302") || strings.Contains(string(productid), "089a") {
				// Only one per host supported
				devices = append(devices, "/dev/bus/usb")
				break
			}
		}
	}

	return devices
}

func CheckTPUHealth(device string) string {
	// Check M.2/mPCIe TPU temperature
	if strings.Contains(device, "usb") == false {
		path := "/sys/class/apex/" + filepath.Base(device)
		temp_b, _ := os.ReadFile(path + "/temp")
		trip_point0_temp_b, _ := os.ReadFile(path + "/trip_point0_temp")

		var temp int
		_, err := fmt.Sscan(string(temp_b), &temp)
		if err != nil {
			log.Println(err)
			return pluginapi.Unhealthy
		}

		var trip_point0_temp int
		_, err2 := fmt.Sscan(string(trip_point0_temp_b), &trip_point0_temp)
		if err2 != nil {
			log.Println(err2)
			return pluginapi.Unhealthy
		}

		if temp >= trip_point0_temp {
			log.Println("Device ", filepath.Base(device), " is overheating (", (temp / 1000), "C)")
			return pluginapi.Unhealthy
		}
	}
	return pluginapi.Healthy
}

func (dp *TPUDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	for {
		devs := []*pluginapi.Device{}
		for _, path := range FindTPUs() {
			dev := &pluginapi.Device{
				ID:     path,
				Health: CheckTPUHealth(path),
			}
			devs = append(devs, dev)
		}
		s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})

		time.Sleep(5 * time.Second)
	}
	return nil
}

func (dp *TPUDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
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

func (TPUDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{PreStartRequired: false, GetPreferredAllocationAvailable: false}, nil
}

func (TPUDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}

func (dp *TPUDevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return nil, nil
}

type TPULister struct {
}

func (l TPULister) GetResourceNamespace() string {
	return "coral.ai"
}

func (l TPULister) Discover(pluginListCh chan dpm.PluginNameList) {
	plugins := make(dpm.PluginNameList, 0)
	plugins = append(plugins, "tpu")
	pluginListCh <- plugins
}

func (l TPULister) NewPlugin(name string) dpm.PluginInterface {
	return &TPUDevicePlugin{
		name: name,
	}
}
