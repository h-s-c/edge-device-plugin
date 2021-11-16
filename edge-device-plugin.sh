#!/bin/bash
/bin/edge-device-plugin --accelerator=tpu &
/bin/edge-device-plugin --accelerator=vpu &
wait -n
exit $?