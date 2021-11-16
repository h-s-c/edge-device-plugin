# Kubernetes device plugin for low power edge computing accelerators

### Supported devices
- Coral Edge TPU (USB\*, M.2, mPCIe)
- Intel Movidius VPU (Neural Compute Stick 2\*)

\*only one USB accelerator per host supported

### Pull request or hardware welcome
- Intel Movidius VPU (M.2, mPCIe)
- Kneron KL520 (M.2, mPCIe)

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