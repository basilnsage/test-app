#!/bin/bash

services=(postService commentService queryService moderationService eventBus)
cwd="$(pwd)"
for service in "${services[@]}"
do
    cd $service
    sed -i'.bak' -e '/basilnsage\/test-app\/shared/d' go.mod
    # go get -u github.com/basilnsage/test-app/shared@master
    go mod tidy
    cd $cwd
done
