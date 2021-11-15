package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type EdgeDevicePlugin struct {
	name    string
	devices []string
}

func (dp *EdgeDevicePlugin) Start() error {
	return nil
}

func FindDevices(devices []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}
		if !info.IsDir() && len(filepath.Base(path)) >= 4 && filepath.Base(path)[0:3] == "apex" {
			log.Println("Found device: ", path)
			devices = append(devices, filepath.Base(path))
		}
		return nil
	}
}

func (dp *EdgeDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	err := filepath.Walk("/dev/", FindDevices(dp.devices))
	if err != nil {
		log.Println(err)
	}

	devs := []*pluginapi.Device{}
	for _, path := range dp.devices {
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
			dev := &pluginapi.DeviceSpec{
				HostPath:      path,
				ContainerPath: path,
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
