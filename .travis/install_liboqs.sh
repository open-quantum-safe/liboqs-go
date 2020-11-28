#!/bin/sh

git clone --branch main --single-branch --depth 1 https://github.com/open-quantum-safe/liboqs
cd liboqs
mkdir build && cd build
cmake -GNinja -DBUILD_SHARED_LIBS=ON -DOQS_BUILD_ONLY_LIB=ON ..
ninja
sudo ninja install
