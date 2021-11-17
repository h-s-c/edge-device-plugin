# Kubernetes device plugin for low power edge computing accelerators

### Supported devices
- Coral Edge TPU (USB\*, M.2, mPCIe)
- Intel Movidius Myriad X VPU (Neural Compute Stick 2\*)
- Broadcom VideoCore (Raspberry Pi 3/4\*)

\*only one accelerator per host supported\

### Send me some hardware
- Intel Movidius Myriad X VPU (M.2, mPCIe)
- Kneron KL520/KL720 (USB, M.2, mPCIe)

## Install DaemonSet
### Helm chart
```bash
helm repo add edge-device-plugin https://h-s-c.github.io/edge-device-plugin
helm install edge-device-plugin edge-device-plugin/edge-device-plugin
```

### Manually
```bash
kubectl create -f edge-device-plugin-daemonset.yaml
```

## Configure your Pod:
```yaml
resources: 
  requests:
    coral.ai/tpu: 1 # requesting 1 TPU
  limits:
    coral.ai/tpu: 1 # requesting 1 TPU
```
```yaml
resources: 
  requests:
    intel.com/vpu: 1 # requesting 1 VPU
  limits:
    intel.com/vpu: 1 # requesting 1 VPU
```
```yaml
# for hardware accelerated video encoding/decoding via /dev/vchiq
resources: 
  requests:
    broadcom.com/gpu: 1 # requesting 1 GPU
  limits:
    broadcom.com/gpu: 1 # requesting 1 GPU
```