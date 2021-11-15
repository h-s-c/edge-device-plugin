package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type EdgeDevicePlugin struct {
	name string
}

func (dp *EdgeDevicePlugin) Start() error {
	return nil
}

func FindDevices() []string {
	matches, err := filepath.Glob("/sys/class/apex/apex*")
	if err != nil {
		log.Println(err)
	}

	devices := []string{}
	for _, path := range matches {
		devices = append(devices, filepath.Base(path))
	}
	return devices
}

func CheckDeviceHealth(device string) string {
	// Check TPU temperature
	temp_b, _ := os.ReadFile("/sys/class/apex/" + device + "/temp")
	trip_point0_temp_b, _ := os.ReadFile("/sys/class/apex/" + device + "/trip_point0_temp")

	temp, _ := strconv.ParseInt(string(temp_b), 10, 64)
	trip_point0_temp, _ := strconv.ParseInt(string(trip_point0_temp_b), 10, 64)

	if temp >= trip_point0_temp {
		log.Println("Device ", device, " is overheating (", (temp / 1000), "C)")
		return pluginapi.Unhealthy
	} else {
		return pluginapi.Healthy
	}

}

func (dp *EdgeDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	for {
		devs := []*pluginapi.Device{}
		for _, path := range FindDevices() {
			dev := &pluginapi.Device{
				ID:     path,
				Health: CheckDeviceHealth(path),
			}
			devs = append(devs, dev)
		}
		s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})

		time.Sleep(5 * time.Second)
	}
	return nil
}

func (dp *EdgeDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	responses := pluginapi.AllocateResponse{}
	for _, req := range r.ContainerRequests {
		response := pluginapi.ContainerAllocateResponse{}
		for _, id := range req.DevicesIDs {
			log.Println("Allocating device: ", id)
			dev := &pluginapi.DeviceSpec{
				HostPath:      "/dev/" + id,
				ContainerPath: "/dev/" + id,
				Permissions:   "rw",
			}
			response.Devices = append(response.Devices, dev)
		}
		responses.ContainerResponses = append(responses.ContainerResponses, &response)
	}

	return &responses, nil
}

func (EdgeDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{PreStartRequired: false, GetPreferredAllocationAvailable: false}, nil
}

func (EdgeDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}

func (dp *EdgeDevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return nil, nil
}

type Lister struct {
}

func (l Lister) GetResourceNamespace() string {
	return "coral.ai"
}

func (l Lister) Discover(pluginListCh chan dpm.PluginNameList) {
	plugins := make(dpm.PluginNameList, 0)
	plugins = append(plugins, "tpu")
	pluginListCh <- plugins
}

func (l Lister) NewPlugin(name string) dpm.PluginInterface {
	return &EdgeDevicePlugin{
		name: name,
	}
}

func main() {
	flag.Parse()
	log.Println("Edge device plugin for Kubernetes")
	l := Lister{}
	manager := dpm.NewManager(&l)
	manager.Run()
}
