# Kubernetes device plugin for low power edge computing accelerators

### Supported devices
- Coral Edge TPU (USB, M.2, mPCIe)

### Upcoming devices
- Intel Movidius VPU (USB)

### Send me some hardware
- Intel Movidius VPU (M.2, mPCIe)
- Kneron KL520 (M.2, mPCIe)

## Install DaemonSet:
```bash
kubectl create -f edge-device-plugin-daemonset.yaml
```

## Configure your Pod:
```yaml
resources: 
  requests:
    coral.ai/tpu: 1 # requesting 1 TPUs
  limits:
    coral.ai/tpu: 1 # requesting 1 TPUs
```