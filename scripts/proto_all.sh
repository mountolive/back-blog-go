#! /bin/bash

echo "GENERATING PROTOS USER PROJECT"
cd user || return
make proto-gen
cd ..
echo "GENERATING PROTOS POST PROJECT"
cd post || return
make proto-gen
cd ..
echo "GENERATING PROTOS GATEWAY PROJECT"
cd gateway || return
make build
