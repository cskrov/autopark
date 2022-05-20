#!/bin/sh

GOOS=linux go build -o autopark-x86-64-linux .
GOARCH=arm64 GOOS=linux go build -o autopark-aarch64-linux .
