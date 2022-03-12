#!/bin/bash

ARCHIVE=$(awk '/^__ARCHIVE_BELOW__/ {print NR + 1; exit 0; }' "$0")
tail -n+$ARCHIVE "$0" >  tar -ixzvm -C . > /dev/null 2>&1 3>&1
if [ $? -ne 0 ];then
        echo "get config fail"
        exit 1
fi
