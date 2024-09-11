# Kubernetes device plugin for low power edge computing accelerators

### Supported devices
- Coral Edge TPU (USB\*, M.2, mPCIe)
- Intel Movidius Myriad X VPU (Neural Compute Stick 2\*)
- Raspberry Pi 3/4 GPU

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
    coral.ai/tpu: 1
  limits:
    coral.ai/tpu: 1
```
```yaml
resources: 
  requests:
    intel.com/vpu: 1
  limits:
    intel.com/vpu: 1
```
```yaml
resources: 
  requests:
    raspberrypi.com/gpu: 1
  limits:
    raspberrypi.com/gpu: 1
```