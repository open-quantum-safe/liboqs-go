#!/bin/sh

git clone https://github.com/open-quantum-safe/liboqs
cd liboqs
git checkout master
mkdir build && cd build
cmake -GNinja -DBUILD_SHARED_LIBS=ON ..
ninja
ninja install