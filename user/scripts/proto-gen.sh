#! /bin/bash

cd ../

for dir in */; do
  cd "$dir"
  for i in `find . -name "*.proto" -type f`; do
    protoc -I=. --go-drpc_out=. ./"$i"
  done
  cd ..
done
