package main

import (
	"flag"
	"log"
	"path/filepath"

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
	if len(matches) > 0 {
		for _, path := range matches {
			log.Println("Found device: ", filepath.Base(path))
			devices = append(devices, filepath.Base(path))
		}
	}
	return devices
}

func (dp *EdgeDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	devs := []*pluginapi.Device{}
	for _, path := range FindDevices() {
		dev := &pluginapi.Device{
			ID:     path,
			Health: pluginapi.Healthy,
		}
		devs = append(devs, dev)
	}

	s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})
	return nil
}

func (dp *EdgeDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	response := pluginapi.AllocateResponse{}

	for _, req := range r.ContainerRequests {
		car := pluginapi.ContainerAllocateResponse{}

		for _, path := range req.DevicesIDs {
			log.Println("Allocating device: ", path)
			dev := &pluginapi.DeviceSpec{
				HostPath:      "/dev/" + path,
				ContainerPath: "/dev/" + path,
				Permissions:   "rw",
			}
			car.Devices = append(car.Devices, dev)
		}

		response.ContainerResponses = append(response.ContainerResponses, &car)
	}

	return &response, nil
}

func (EdgeDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

func (EdgeDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

func (dp *EdgeDevicePlugin) GetPreferredAllocation(ctx context.Context, request *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return &pluginapi.PreferredAllocationResponse{}, nil
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
