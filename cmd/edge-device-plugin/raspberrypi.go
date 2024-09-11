package main

import (
	"log"
	"os"
	"time"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type RasberrypiDevicePlugin struct {
	name string
}

func (dp *RasberrypiDevicePlugin) Start() error {
	return nil
}

func (dp *RasberrypiDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	devs := []*pluginapi.Device{}
	// Detect Raspberry Pi
	if _, err := os.Stat("/sys/class/vc-mem"); err == nil || os.IsExist(err) {
		dev := &pluginapi.Device{
			ID:     "/dev/vchiq /dev/vcsm-cma /dev/video10 /dev/video11 /dev/video12 /dev/dri",
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

func (dp *RasberrypiDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	responses := pluginapi.AllocateResponse{}
	for _, req := range r.ContainerRequests {
		response := pluginapi.ContainerAllocateResponse{}
		for _, id := range req.DevicesIDs {
			log.Println("Allocating devices: ", id)

			// OpenMAX Video
			dev1 := &pluginapi.DeviceSpec{
				HostPath:      "/dev/vchiq",
				ContainerPath: "/dev/vchiq",
				Permissions:   "rw",
			}
			response.Devices = append(response.Devices, dev1)
			dev2 := &pluginapi.DeviceSpec{
				HostPath:      "/dev/vcsm-cma",
				ContainerPath: "/dev/vcsm-cma",
				Permissions:   "rw",
			}
			response.Devices = append(response.Devices, dev2)
			// V4L2 Video
			dev3 := &pluginapi.DeviceSpec{
				HostPath:      "/dev/video10",
				ContainerPath: "/dev/video10",
				Permissions:   "rw",
			}
			response.Devices = append(response.Devices, dev3)
			dev4 := &pluginapi.DeviceSpec{
				HostPath:      "/dev/video11",
				ContainerPath: "/dev/video11",
				Permissions:   "rw",
			}
			response.Devices = append(response.Devices, dev4)
			dev5 := &pluginapi.DeviceSpec{
				HostPath:      "/dev/video12",
				ContainerPath: "/dev/video12",
				Permissions:   "rw",
			}
			response.Devices = append(response.Devices, dev5)
			// GPU
			dev6 := &pluginapi.DeviceSpec{
				HostPath:      "/dev/dri",
				ContainerPath: "/dev/dri",
				Permissions:   "rw",
			}
			response.Devices = append(response.Devices, dev6)

		}
		responses.ContainerResponses = append(responses.ContainerResponses, &response)
	}

	return &responses, nil
}

func (RasberrypiDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{PreStartRequired: false, GetPreferredAllocationAvailable: false}, nil
}

func (RasberrypiDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}

func (dp *RasberrypiDevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return nil, nil
}

type RasberrypiLister struct {
}

func (l RasberrypiLister) GetResourceNamespace() string {
	return "raspberrypi.com"
}

func (l RasberrypiLister) Discover(pluginListCh chan dpm.PluginNameList) {
	plugins := make(dpm.PluginNameList, 0)
	plugins = append(plugins, "gpu")
	pluginListCh <- plugins
}

func (l RasberrypiLister) NewPlugin(name string) dpm.PluginInterface {
	return &RasberrypiDevicePlugin{
		name: name,
	}
}
