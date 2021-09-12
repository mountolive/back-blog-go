#! /bin/bash

BASE_DIR="grpc"

PKG_DIR="$(pwd)/$BASE_DIR"

mkdir -p $PKG_DIR

cd ../proto/user

protoc -I=. --go_out="$PKG_DIR" --go-grpc_out="$PKG_DIR" *.proto
