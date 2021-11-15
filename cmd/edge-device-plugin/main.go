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
	for _, path := range matches {
		log.Println("Found device: ", filepath.Base(path))
		devices = append(devices, filepath.Base(path))
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
	log.Println("Allocating devices")

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
