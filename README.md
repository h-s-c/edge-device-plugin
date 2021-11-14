# Kubernetes device plugin for low power edge computing accelerators

### Supported devices
- Coral Edge TPU (M.2, mPCIe)

### Upcoming devices
- Coral Edge TPU (USB)
- Intel Movidius VPU (USB)

### Send me some hardware
- Intel Movidius VPU (M.2, mPCIe)
- Kneron KL520 (M.2, mPCIe)

## Install DaemonSet:
```bash
kubectl create -f edge-device-plugin-daemonset.yaml
```

## Configure your pod:
```yaml
spec:
  containers:
    resources:
    limits:
      coral.ai/tpu: 2 # requesting 2 TPUs
```