liboqs-go version 0.7.2
=======================

About
-----

The **Open Quantum Safe (OQS) project** has the goal of developing and prototyping quantum-resistant cryptography. More information on OQS can be found on our website: https://openquantumsafe.org/ and on Github at https://github.com/open-quantum-safe/.

**liboqs** is an open source C library for quantum-resistant cryptographic algorithms. See more about liboqs at [https://github.com/open-quantum-safe/liboqs/](https://github.com/open-quantum-safe/liboqs/), including a list of supported algorithms.

**liboqs-go** is an open source Go wrapper for the liboqs C library for quantum-resistant cryptographic algorithms. Details about liboqs-go can be found in [README.md](https://github.com/open-quantum-safe/liboqs-go/blob/main/README.md). See in particular limitations on intended use.

Release notes
=============

This release of liboqs-go was released on August 26, 2022. Its release page on 
GitHub is https://github.com/open-quantum-safe/liboqs-go/releases/tag/0.7.2.

What's New
----------

This is the 10th release of liboqs-go. 

This is an incremental minor change release.

For a list of changes see [CHANGES.txt](https://github.com/open-quantum-safe/liboqs-go/blob/main/CHANGES.txt).

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
