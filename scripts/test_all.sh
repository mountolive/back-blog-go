#! /bin/bash

echo "TESTING USER PROJECT"
cd user || return
make test
cd ..
echo "TESTING POST PROJECT"
cd post || return
make test
cd ..
