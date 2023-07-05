FROM ubuntu:latest

# Install dependencies
RUN apt-get -y update && \
    apt-get install -y build-essential git cmake libssl-dev golang

# Get liboqs
RUN git clone --depth 1 --branch main https://github.com/open-quantum-safe/liboqs

# Install liboqs
RUN cmake -S liboqs -B liboqs/build -DBUILD_SHARED_LIBS=ON && \
    cmake --build liboqs/build --parallel 4 && \
    cmake --build liboqs/build --target install

# Enable a normal user
RUN useradd -m -c "Open Quantum Safe" oqs
USER oqs
WORKDIR /home/oqs

# Get liboqs-go
RUN git clone --depth 1 --branch main https://github.com/open-quantum-safe/liboqs-go.git

# Configure liboqs-go
ENV PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/home/oqs/liboqs-go/.config
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
