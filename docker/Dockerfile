FROM ubuntu:22.04 AS build

WORKDIR /app
RUN apt-get update && apt-get install -y astyle cmake gcc ninja-build libssl-dev python3-pytest python3-pytest-xdist unzip xsltproc doxygen graphviz python3-yaml valgrind git \
&& git clone -b main https://github.com/open-quantum-safe/liboqs.git
RUN mkdir -p /app/liboqs/build
WORKDIR /app/liboqs/build
RUN cmake -GNinja .. -DBUILD_SHARED_LIBS=ON -DOQS_BUILD_ONLY_LIB=ON -DOQS_DIST_BUILD=ON \
&& ninja

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y astyle cmake gcc ninja-build libssl-dev python3-pytest python3-pytest-xdist unzip xsltproc doxygen graphviz python3-yaml valgrind git golang \
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
