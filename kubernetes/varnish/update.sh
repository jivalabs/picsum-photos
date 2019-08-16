#!/usr/bin/env bash

set -ex

docker build . -t DMarby/picsum-photos-varnish:latest
docker push DMarby/picsum-photos-varnish:latest
