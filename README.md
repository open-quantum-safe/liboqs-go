# liboqs-go: Go bindings for liboqs

[![GitHub actions](https://github.com/open-quantum-safe/liboqs-go/actions/workflows/go.yml/badge.svg)](https://github.com/open-quantum-safe/liboqs-go/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-quantum-safe/liboqs-go)](https://goreportcard.com/report/github.com/open-quantum-safe/liboqs-go)
[![Documentation](https://godoc.org/github.com/open-quantum-safe/liboqs-go/oqs?status.svg)](https://pkg.go.dev/github.com/open-quantum-safe/liboqs-go/oqs)

---

## About

The **Open Quantum Safe (OQS) project** has the goal of developing and
prototyping quantum-resistant cryptography.

**liboqs-go** offers a Go wrapper for the
[Open Quantum Safe](https://openquantumsafe.org/)
[liboqs](https://github.com/open-quantum-safe/liboqs/) C library, which is a C
library for quantum-resistant cryptographic algorithms.

liboqs-go is a Go package, hence in the following it is assumed that you have
access to a Go compliant environment. liboqs-go has been extensively tested on
Linux, macOS and Windows platforms. Continuous integration is provided via
GitHub actions.

The project contains the following files and directories:

- **`oqs/oqs.go`: main package file for the wrapper**
- `.config/liboqs-go.pc`: `pkg-config` configuration file needed by `cgo`
- `.config-static/liboqs-go.pc`: `pkg-config` configuration file needed by
  `cgo` when linking statically against liboqs
- `examples`: usage examples, including a client/server KEM over TCP/IP
- `oqstests`: unit tests

---

## Pre-requisites

- [liboqs](https://github.com/open-quantum-safe/liboqs)
- [git](https://git-scm.com/)
- [CMake](https://cmake.org/)
- C compiler, e.g., [gcc](https://gcc.gnu.org/)
  , [clang](https://clang.llvm.org)
  , [MSYS2](https://www.msys2.org/) etc.
- [Go 1.21 or later](https://go.dev/)
- `pkg-config` (use `sudo apt-get install pkg-config` to install on
  Ubuntu/Debian-based Linux platforms or install it
  via a third-party compiler such as [MSYS2](https://www.msys2.org/) on
  Windows)
- If using Windows, you need a C compiler supported by `cgo` added to your
  `PATH` environment variable; currently, the best supported ones are provided
  by [MSYS2](https://www.msys2.org/)
  and [`tdm-gcc`](https://jmeubank.github.io/tdm-gcc/);
  [Cygwin](https://www.cygwin.com/) is **not yet supported**
  by `cgo`; we recommend using MSYS2 since it also contains `pkg-config` as a
  package; to install `gcc` and `pkg-config` under MSYS2, please execute in a
  MSYS2 terminal window
  `pacman -S mingw64/mingw-w64-x86_64-gcc mingw64/mingw-w64-x86_64-pkg-config`,
  then add the corresponding installation location (e.g,
  `C:\msys64\mingw64\bin`) to your `PATH` environment variable by executing
  `set PATH=%PATH%;C:\msys64\mingw64\bin`.

  **Very important:** make sure that
  the `PATH` entry to the `gcc` and `pkg-config` provided by `MSYS2`comes
  **before** any other (if any) `gcc` and `pkg-config` executables you may have
  installed (e.g. such as the ones provided
  by [Cygwin](https://www.cygwin.com)). To verify, type into a Command Prompt
  `gcc --version`, and you should get an output such as

  > gcc (Rev3, Built by MSYS2 project) 9.1.0

---

## Functional restrictions

Please note that on some platforms not all algorithms are supported:

- macOS/Darwin: The Rainbow and Classic-McEliece algorithm families as well as
  HQC-256 do not work.
- Windows: The Rainbow and Classic-McEliece algorithm families do not work.

---

## Installation

In the rest of this document, we assume you execute commands from inside the
`$HOME` directory on UNIX-like systems, or from inside the `%USERPROFILE%` on
Windows.

### Configure, build and install liboqs

Execute in a Terminal/Console/Administrator Command Prompt

```shell
git clone --depth=1 https://github.com/open-quantum-safe/liboqs
cmake -S liboqs -B liboqs/build -DBUILD_SHARED_LIBS=ON
cmake --build liboqs/build --parallel 8
cmake --build liboqs/build --target install
```

The last line may require prefixing it by `sudo` on UNIX-like systems.
Change `--parallel 8` to match the number of available cores on your system.

On UNIX-like platforms, you may need to set
the `LD_LIBRARY_PATH` (`DYLD_LIBRARY_PATH` on macOS) environment variable to
point to the path to liboqs' library directory, e.g.,

```shell
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
```

On Windows platforms, **you must ensure** that the liboqs shared
library `oqs.dll` is visible system-wide, and that the following environment
variable are being set. Use the "Edit the system environment variables" Control
Panel tool or execute in a Command Prompt, e.g.,

```shell
set PATH=%PATH%;C:\Program Files (x86)\liboqs\bin
```

You can change liboqs' installation directory by configuring the build to use
an alternative path, e.g., `C:\liboqs`, by replacing the first CMake line above
by

```shell
cmake -S liboqs -B liboqs/build -DCMAKE_INSTALL_PREFIX="C:\liboqs" -DBUILD_SHARED_LIBS=ON
```

### Configure and install the wrapper

Execute in a Terminal/Console/Administrator Command Prompt

```shell
git clone --depth=1 https://github.com/open-quantum-safe/liboqs-go
```

Next, you must modify the following lines in
[`$HOME/liboqs-go/.config/liboqs-go.pc`](https://github.com/open-quantum-safe/liboqs-go/tree/main/.config/liboqs-go.pc)

    LIBOQS_INCLUDE_DIR=/usr/local/include
    LIBOQS_LIB_DIR=/usr/local/lib

so they correspond to your liboqs include/lib installation directories. On
Windows, **using forward slashes `/` and not
back-slashes**, e.g.,

    LIBOQS_INCLUDE_DIR=C:/Program Files (x86)/liboqs/bin
    LIBOQS_LIB_DIR=C:/Program Files (x86)/liboqs/lib

Finally, you must add/append the `$HOME/liboqs-go/.config` directory to the
`PKG_CONFIG_PATH` environment variable, i.e., on UNIX-like systems execute in a
terminal

```shell
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$HOME/liboqs-go/.config
```

or, on Windows platforms, use the "Edit the system environment variables"
Control Panel tool or execute in a Command Prompt

```shell
set PKG_CONFIG_PATH=%PKG_CONFIG_PATH%;$HOME/liboqs-go/.config
```

### Linking statically against liboqs - excluding macOS/OS X platforms

Replace `.config` with `.config-static` when setting the `PKG_CONFIG_PATH`
environment variable above. This assumes that you previously compiled and
installed the static version of liboqs, i.e., you did not pass
`-DBUILD_SHARED_LIBS=ON` to CMake when configuring liboqs above.

Note that `.config-static/liboqs-go.pc` links statically against OpenSSL as
well. In case you don't have OpenSSL installed, remove the `-lcrypto` from the
last line of `.config-static/liboqs-go.pc`, and make sure you compiled liboqs
without OpenSSL, i.e., pass the `-DOQS_USE_OPENSSL=OFF` CMake flag when
configuring liboqs, otherwise you will get linker errors.

**Important:** Ensure that you run `go clean -cache` before building or
running, so `pkg-config` refreshes its cache.

### Linking statically against liboqs - macOS/OS X platforms

The macOS/OS X linker does not allow choosing static vs dynamic linking when
both static and dynamic versions of a library are installed. In this case, the
dynamic version will always be chosen by the linker. Hence, to link statically
agains liboqs on macOS/OS X, make sure you have not installed the dynamic
version of liboqs anywhere on your system, and use the `.config` (not
`.config-static`) when setting the `PKG_CONFIG_PATH` environment variable.

**Important:** Ensure that you run `go clean -cache` before building or
running.

### Run the examples

From inside the `liboqs-go` directory, execute

```shell
go run examples/kem/kem.go
go run examples/sig/sig.go
go run examples/rand/rand.go
```

### Build executables

Replace `go run` by `go build`, e.g., `go build examples/kem/kem.go`.

**Note** go binaries produced on macOS arm64 are not code-signed properly. See
[https://github.com/golang/go/issues/63997](https://github.com/golang/go/issues/63997).

To fix, run

```shell
codesign -f -s - path/to/executable
```

### Run the unit tests

From inside the `liboqs-go` directory, execute

```shell
cd liboqs-go
go test -v ./oqstests
```

On Windows, you may need to replace forward-slashes `/` by back-slashes `\`.

---

## Usage in standalone applications

liboqs-go can be imported into Go programs with

```go
import (
    "github.com/open-quantum-safe/liboqs-go/oqs"
)
```

The examples in the
[`examples`](https://github.com/open-quantum-safe/liboqs-go/tree/main/examples)
directory are self-explanatory and provide more details about the wrapper's
API.

---

## Documentation

The `liboqs-go` wrapper is fully documented using the Go standard documentation
conventions. For example, to read the full documentation about the
`oqs.Signature.Verify` method, execute from inside the `liboqs-go` directory

```shell
go doc liboqs-go/oqs.Signature.Verify
```

For the RNG-related function, execute e.g.

```shell
go doc liboqs-go/oqs/rand.RandomBytes
```

For automatically-generated documentation in HTML format,
click [here](https://pkg.go.dev/github.com/open-quantum-safe/liboqs-go/oqs).

For the RNG-related documentation, click
[here](https://pkg.go.dev/github.com/open-quantum-safe/liboqs-go/oqs/rand).

---

## Docker

A self-explanatory minimalistic Docker file is provided in
[`Dockerfile`](https://github.com/open-quantum-safe/liboqs-go/tree/main/Dockerfile).

Build the image by executing

```shell
docker build -t oqs-go .
```

Run, e.g., the key encapsulation example by executing

```shell
docker run -it oqs-go sh -c "cd liboqs-go && go run examples/kem/kem.go"
```

Or, run the unit tests with

```shell
docker run -it oqs-go sh -c "cd liboqs-go && go test -v ./oqstests"
```

In case you want to use the Docker container as a development environment,
mount your current project in the Docker container with

```shell
docker run --rm -it --workdir=/app -v ${PWD}:/app oqs-go /bin/bash
```

---

## Limitations and security

liboqs is designed for prototyping and evaluating quantum-resistant
cryptography. Security of proposed quantum-resistant algorithms may rapidly
change as research advances, and may ultimately be completely insecure against
either classical or quantum computers.

We believe that the NIST Post-Quantum Cryptography standardization project is
currently the best avenue to identifying potentially quantum-resistant
algorithms. liboqs does not intend to "pick winners", and we strongly recommend
that applications and protocols rely on the outcomes of the NIST
standardization project when deploying post-quantum cryptography.

We acknowledge that some parties may want to begin deploying post-quantum
cryptography prior to the conclusion of the NIST standardization project. We
strongly recommend that any attempts to do make use of so-called
**hybrid cryptography**, in which post-quantum public-key algorithms are used
alongside traditional public key algorithms (like RSA or elliptic curves) so
that the solution is at least no less secure than existing traditional
cryptography.

Just like liboqs, liboqs-go is provided "as is", without warranty of any kind.
See [LICENSE](https://github.com/open-quantum-safe/liboqs-go/blob/main/LICENSE)
for the full disclaimer.

---

## License

liboqs-go is licensed under the MIT License;
see [LICENSE](https://github.com/open-quantum-safe/liboqs-go/blob/main/LICENSE)
for details.

---

## Team

The Open Quantum Safe project is led by
[Douglas Stebila](https://www.douglas.stebila.ca/research/) and
[Michele Mosca](http://faculty.iqc.uwaterloo.ca/mmosca/) at the University of
Waterloo.

liboqs-go was developed by [Vlad Gheorghiu](https://vsoftco.github.io) at
[softwareQ Inc.](https://www.softwareq.ca) and at the University of Waterloo.

---

## Support

Financial support for the development of Open Quantum Safe has been provided by
Amazon Web Services and the Canadian Centre for Cyber Security.

We'd like to make a special acknowledgement to the companies who have dedicated
programmer time to contribute source code to OQS, including Amazon Web
Services, evolutionQ, softwareQ, and Microsoft Research.

Research projects which developed specific components of OQS have been
supported by various research grants, including funding from the Natural
Sciences and Engineering Research Council of Canada (NSERC); see the source
papers for funding acknowledgments.
