# `Dockerfile for liboqs-go`

This is a multistage Dockerfile for building liboqs and setting up liboqs-go.

You can pull the container from Dockerhub with:
```bash
docker pull openquantumsafe/go
```
You can build the container with:
```bash
docker build -t openquantumsafe/go .
```
You can start the container and run the tests with:
```bash
docker run openquantumsafe/go
```
If you want to mount your current directory into the container and use it as an development environment:
```bash
 docker run --rm -it --workdir=/app -v ${PWD}:/app openquantumsafe/go /bin/bash 
```
# `Contribution`

If you want to contribute, please make sure that the dockerfile passes all hadolint tests.
After intalling hadolint, you can check it with:
```
docker run --rm -i hadolint/hadolint < Dockerfile
```
