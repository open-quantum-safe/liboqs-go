FROM ubuntu:22.04 AS build

WORKDIR /app
RUN apt-get update && apt-get install -y astyle cmake gcc ninja-build libssl-dev python3-pytest python3-pytest-xdist unzip xsltproc doxygen graphviz python3-yaml valgrind git \
&& git clone -b main https://github.com/open-quantum-safe/liboqs.git
RUN mkdir -p /app/liboqs/build
WORKDIR /app/liboqs/build
RUN cmake -GNinja .. -DBUILD_SHARED_LIBS=ON -DOQS_BUILD_ONLY_LIB=ON -DOQS_DIST_BUILD=ON \
&& ninja install

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y golang \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*
COPY --from=build /usr/local/lib/liboqs* /usr/local/lib
COPY --from=build /usr/local/include/oqs /usr/local/include/oqs
WORKDIR /home
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig:/home/liboqs-go/.config
WORKDIR /home/liboqs-go
COPY . .
RUN ls -al
#RUN go test -v ./oqstests
CMD ["go", "test", "-v", "./oqstests"]
