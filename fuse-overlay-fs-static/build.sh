#!/bin/bash

set -ex
set -o pipefail

# rm -rf fuse-overlayfs
# git clone https://github.com/containers/fuse-overlayfs.git
# cd fuse-overlayfs 
# 
# 
# ./autogen.sh
# ./configure LDFLAGS="-static"
cd fuse-overlayfs && make clean && make build

# mv fuse-overlayfs ../fuse-overlayfs-bin