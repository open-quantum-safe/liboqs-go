FROM ubuntu:22.04 AS build

WORKDIR /app
RUN apt-get update && apt-get install --no-install-recommends -y \
  astyle=3.1-2build1 -V \
  ca-certificates=20211016 -V \
  cmake=3.22.1-1ubuntu1.22.04.1 -V \
  doxygen=1.9.1-2ubuntu2 -V \
  gcc=4:11.2.0-1ubuntu1 -V \
  gcc-multilib=4:11.2.0-1ubuntu1 -V \
  git=1:2.34.1-1ubuntu1.5 -V \
  golang=2:1.18~0ubuntu2 -V \
  graphviz=2.42.2-6 -V \
  libssl-dev=3.0.2-0ubuntu1.7 -V \
  ninja-build=1.10.1-1 -V \
  python3-pytest=6.2.5-1ubuntu2 -V \
  python3-pytest-xdist=2.5.0-1 -V \
  python3-yaml=5.4.1-1ubuntu1 -V \
  unzip=6.0-26ubuntu3.1 -V \
  valgrind=1:3.18.1-1ubuntu2 -V \
  xsltproc=1.1.34-4ubuntu0.22.04.1 -V \
&& git clone -b main https://github.com/open-quantum-safe/liboqs.git
WORKDIR /app/liboqs/build
RUN cmake -GNinja .. -DBUILD_SHARED_LIBS=ON -DOQS_BUILD_ONLY_LIB=ON \
&& ninja

FROM ubuntu:22.04
RUN apt-get update && apt-get install --no-install-recommends -y \
  ca-certificates=20211016 -V \
  gcc=4:11.2.0-1ubuntu1 -V \
  gcc-multilib=4:11.2.0-1ubuntu1 -V \
  git=1:2.34.1-1ubuntu1.5 -V \
  golang=2:1.18~0ubuntu2 -V \
  pkg-config=0.29.2-1ubuntu3 -V \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*
COPY --from=build /app/liboqs/build/lib /usr/local/lib
COPY --from=build /app/liboqs/build/include /usr/include
WORKDIR /home
RUN git clone https://github.com/open-quantum-safe/liboqs-go
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig:/home/liboqs-go/.config
WORKDIR /home/liboqs-go
RUN go test -c ./oqstests
CMD ["./oqstests.test"]
