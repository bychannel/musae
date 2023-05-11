#!/usr/bin/env bash
go run ./build.go -srcDir=./ -outDir=./cmd -protoc=protoc
mv ./cmd/proto_msg.pb.go ../proto_msg.pb.go