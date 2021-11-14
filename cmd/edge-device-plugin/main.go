package main

import (
	"log"
	"os"
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

func (dp *EdgeDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {

	devs := make([]*pluginapi.Device, 0)

	err := filepath.Walk("/dev/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}
		if !info.IsDir() && filepath.Base(path)[0:3] == "apex" {
			log.Println("Found device: ", path)
			dev := &pluginapi.Device{
				ID:     path,
				Health: pluginapi.Healthy,
			}
			devs = append(devs, dev)
		}
		return nil
	})

	if err != nil {
		log.Println(err)
	}

	s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})
	return nil
}

func (dp *EdgeDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {

	car := pluginapi.ContainerAllocateResponse{}

	err := filepath.Walk("/dev/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}
		if !info.IsDir() && filepath.Base(path)[0:3] == "apex" {
			log.Println("Allocating device: ", path)
			dev := &pluginapi.DeviceSpec{
				HostPath:      path,
				ContainerPath: path,
				Permissions:   "rw",
			}
			car.Devices = append(car.Devices, dev)
		}
		return nil
	})

	if err != nil {
		log.Println(err)
	}

	response := &pluginapi.AllocateResponse{
		ContainerResponses: []*pluginapi.ContainerAllocateResponse{
			&car,
		},
	}
	return response, nil
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
	l := Lister{}
	manager := dpm.NewManager(&l)
	manager.Run()
}
