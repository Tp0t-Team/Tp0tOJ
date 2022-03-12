#!/bin/bash

ARCHIVE=$(awk '/^__CONFIG_BELOW__/ {print NR + 1; exit 0; }' "$0")
tail -n+$ARCHIVE "$0" > registries-config.yaml
if [ $? -ne 0 ];then
        echo "get config fail"
        exit 1
fi
