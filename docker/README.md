# `Dockerfile for liboqs-go`

This is a multistage Dockerfile for building liboqs and setting up liboqs-go.

You can build the container with:
```cmd
docker build -t liboqs-go .

```
You can start the container and run the kem example with:
```cmd
docker run liboqs-go

```
If you want to use the container in an interactive shell:
```
docker run -it liboqs-go /bin/bash
```
# `Contribution`

If you want to contribute, please make sure that the dockerfile passes all hadolint tests.
After intalling hadolint, you can check it with:
```
docker run --rm -i hadolint/hadolint < Dockerfile
```

