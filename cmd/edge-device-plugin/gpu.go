package main

import (
	"log"
	"os"
	"time"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type GPUDevicePlugin struct {
	name string
}

func (dp *GPUDevicePlugin) Start() error {
	return nil
}

func (dp *GPUDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	devs := []*pluginapi.Device{}
	if _, err := os.Stat("/sys/class/vchiq"); err == nil || os.IsExist(err) {
		path := "/dev/vchiq"
		dev := &pluginapi.Device{
			ID:     path,
			Health: pluginapi.Healthy,
		}
		devs = append(devs, dev)
	}

	for {
		s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})

		time.Sleep(5 * time.Second)
	}
	return nil
}

func (dp *GPUDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
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

func (GPUDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{PreStartRequired: false, GetPreferredAllocationAvailable: false}, nil
}

func (GPUDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}

func (dp *GPUDevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return nil, nil
}

type GPULister struct {
}

func (l GPULister) GetResourceNamespace() string {
	return "broadcom.com"
}

func (l GPULister) Discover(pluginListCh chan dpm.PluginNameList) {
	plugins := make(dpm.PluginNameList, 0)
	plugins = append(plugins, "gpu")
	pluginListCh <- plugins
}

func (l GPULister) NewPlugin(name string) dpm.PluginInterface {
	return &GPUDevicePlugin{
		name: name,
	}
}
