liboqs-go: Go bindings for liboqs
=================================

[![GitHub actions](https://github.com/open-quantum-safe/liboqs-go/actions/workflows/go.yml/badge.svg)](https://github.com/open-quantum-safe/liboqs-go/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-quantum-safe/liboqs-go)](https://goreportcard.com/report/github.com/open-quantum-safe/liboqs-go)
[![Documentation](https://godoc.org/github.com/open-quantum-safe/liboqs-go/oqs?status.svg)](https://pkg.go.dev/github.com/open-quantum-safe/liboqs-go/oqs)

---

**liboqs-go** offers a Go wrapper for
the [Open Quantum Safe](https://openquantumsafe.org/) [liboqs](https://github.com/open-quantum-safe/liboqs/)
C library, which is a C library for quantum-resistant cryptographic algorithms.

The wrapper is written in Go, hence in the following it is assumed that you have
access to a Go compliant environment. liboqs-go has been extensively tested on
Linux, macOS and Windows platforms. Continuous integration is provided via
CircleCI and AppVeyor.

## <a name="pre-requisites"></a>Pre-requisites

liboqs-go depends on the [liboqs](https://github.com/open-quantum-safe/liboqs) C
library; liboqs must first be compiled as a Linux/macOS/Windows library (i.e.
using `ninja install` with `-DBUILD_SHARED_LIBS=ON` during configuration), see
the specific platform building instructions below.

In addition, we assume you have access to:

- POSIX compliant system (UNIX/Linux/macOS) or Windows
- Go version 1.18 or later
- a standard C compliant compiler (`gcc`/`clang` etc.)
- `pkg-config` (use `sudo apt-get install pkg-config` to install on
  Ubuntu/Debian-based Linux platforms or install it via a third-party compiler
  such as  [MSYS2](https://www.msys2.org/) on Windows)
- if using Windows, you need a C compiler supported by `cgo` added to
  your `PATH` environment variable; currently, the best supported ones are
  provided by [MSYS2](https://www.msys2.org/)
  and [`tdm-gcc`](https://jmeubank.github.io/tdm-gcc/)
  ; [Cygwin](https://www.cygwin.com/) is **not yet supported** by `cgo`; we
  recommend using MSYS2 since it also contains `pkg-config` as a package; to
  install `gcc` and `pkg-config` under MSYS2, please type in a MSYS2 terminal
  window `pacman -Ss mingw64/mingw-w64-x86_64-gcc mingw64/mingw-w64-x86_64-pkg-config`
  , then add the corresponding installation location (in our
  case, `C:\msys64\mingw64\bin`) to your `PATH` environment variable.

### Functional restrictions

Please note that on some platforms not all algorithms are supported:

- Darwin: The Rainbow and Classic-McEliece algorithm families as well as HQC-256 do not work.
- Windows: The Rainbow and Classic-McEliece algorithm families do not work.

<a name="contents"></a>Contents
----

liboqs-go is a Go package. The project contains the following files and
directories:

- **`oqs/oqs.go`: main package file for the wrapper**
- `.config/liboqs.pc`: `pkg-config` configuration file needed by `cgo`
- `examples`: usage examples, including a client/server KEM over TCP/IP
- `oqstests`: unit tests

<a name="usage"></a>Usage
----

The examples in
the [`examples`](https://github.com/open-quantum-safe/liboqs-go/tree/main/examples)
directory are self-explanatory and provide more details about the wrapper's API.

<a name="posix"></a>Running/building on POSIX (Linux/UNIX-like) platforms
----

### <a name="install-oqs"></a>Installing liboqs

First, you must build liboqs according to
the [liboqs building instructions](https://github.com/open-quantum-safe/liboqs#linuxmacos)
with shared library support enabled (add `-DBUILD_SHARED_LIBS=ON` to the `cmake`
command), followed (optionally) by a `sudo ninja install` to ensure that the
shared library is visible system-wide (by default it installs
under `/usr/local/include` and `/usr/local/lib` on Linux/macOS).

You may need to set the `LD_LIBRARY_PATH` (`DYLD_LIBRARY_PATH` on macOS)
environment variable to point to the path to liboqs library directory, e.g.

    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib

assuming `liboqs.so.*` were installed in `/usr/local/lib` (true if you
ran `sudo ninja install` after building liboqs).

### <a name="install-wrapper"></a>Installing the wrapper

Download/clone the `liboqs-go` wrapper repository in the directory of your
choice, e.g. `$HOME`, by typing in a terminal/console

    cd $HOME && git clone https://github.com/open-quantum-safe/liboqs-go

Next, you must modify the following lines
in [`$HOME/liboqs-go/.config/liboqs.pc`](https://github.com/open-quantum-safe/liboqs-go/tree/main/.config/liboqs.pc)

    LIBOQS_INCLUDE_DIR=/usr/local/include
    LIBOQS_LIB_DIR=/usr/local/lib

so they correspond to your C liboqs include/lib installation directories.

Finally, you must add/append the `$HOME/liboqs-go/.config` directory to
the `PKG_CONFIG_PATH` environment variable, i.e.

    export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$HOME/liboqs-go/.config

Once you have configured your system as directed above,
simply `import "github.com/open-quantum-safe/liboqs-go/oqs"` in the Go
application of your choice, initialize the application module
with `go mod init <module_name>`, and finally run it with `go run <module_name>`
or build it with `go build <module_name>`.

To run the examples from the terminal/console, first change directory
to `liboqs-go` by typing in a terminal/console

    cd $HOME/liboqs-go

then run the example(s) by typing e.g.

    go run examples/kem/kem.go 

Replace `go run` with `go build` if you intend to build the corresponding
executable `$HOME/liboqs-go/kem`.

To run the unit tests from the terminal/console, type (still from
inside `$HOME/liboqs-go`) `go test -v ./oqstests` and to build the unit test
executable from the terminal/console, type

    go test -c ./oqstests

which will build the `$HOME/liboqs-go/oqstests.test` executable.

<a name="windows"></a>Running/building on Windows
----

We assume that `liboqs` is installed under `C:\some\dir\liboqs` and
was [successfully built](https://github.com/open-quantum-safe/liboqs#windows)
in`C:\some\dir\liboqs\build`
, and that `liboqs-go` is installed under `C:\some\dir\liboqs-go` (
replace `\some\dir` with your corresponding paths). Ensure that the liboqs
shared library `oqs.dll` is visible system-wide. Use the "Edit the system
environment variables" Control Panel tool or type in a Command Prompt

	set PATH="%PATH%;C:\some\dir\liboqs\build\bin"

of course replacing the paths with the ones corresponding to your system.

As mentioned in the [Pre-requisites](#a-namepre-requisitesapre-requisites) section, we assume you
have installed the  [MSYS2](https://www.msys2.org/) C compiler (which contains
also `pkg-config`). Make sure that the executable `gcc` and `pkg-config` from
MSYS2 are system-wide visible, using e.g. the "Edit the system environment
variables" Control Panel tool to add their corresponding directory to the `PATH`
environment variable. In our case, MSYS2's `gcc` and `pkg-config` are located
under `C:\msys64\mingw64\bin`
,so we need to add that directory to the `PATH`. **Very important:** make sure
that the `PATH` entry to the `gcc` and `pkg-config` provided by `MSYS2`comes **
before** any other (if any) `gcc` and `pkg-config` executables you may have
installed (e.g. such as the ones provided by [Cygwin](https://www.cygwin.com)).
To verify, type into a Command Prompt `gcc --version`, and you should get an
output like

> gcc (Rev3, Built by MSYS2 project) 9.1.0

Next, similarly to the [POSIX instructions](#a-nameposixarunningbuilding-on-posix-linuxunix-like-platforms), modify the
corresponding
lines
in [`liboqs-go\.config\liboqs.pc`](https://github.com/open-quantum-safe/liboqs-go/tree/main/.config/liboqs.pc)
to point to the correct locations, **using forward slashes `/` and not
back-slashes**, e.g.

    LIBOQS_INCLUDE_DIR=C:/some/dir/liboqs/build/include
    LIBOQS_LIB_DIR=C:/some/dir/liboqs/build/lib

Finally, add the `liboqs-go\.config`
directory to the `PKG_CONFIG_PATH` environment variable, using the "Edit the
system environment variables" Control Panel tool or by typing in the Command
Prompt

```bash
set PKG_CONFIG_PATH="C:\some\dir\liboqs-go\.config"
```

Once those steps are completed, you can now test whether everything went OK by
running one of the examples and/or unit tests. First change directory
to `liboqs-go` by typing in a Command Prompt

```bash
cd C:\some\dir\liboqs-go
```

followed by e.g.

```bash
go run examples\kem\kem.go
```

and/or

```bash
go test -v .\oqstests
```

If you do not get any errors and the example (unit tests) is (are) successfully
run, then your installation was successful. For more details about command-line
configuration under Windows see the AppVeyor CI configuration
file [`appveyor.yml`](https://github.com/open-quantum-safe/liboqs-go/tree/main/appveyor.yml).


<a name="Docker containers"></a>Docker container
----

There are two ways to run liboqs-go in a Docker container. You can pull the
official liboqs-go Docker container from Dockerhub:

```bash
docker pull openquantumsafe/go
```

Or you can build the container yourself:

``` bash
cd docker
docker build -t openquantumsafe/go .
```

Run the tests with:

```bash
docker run openquantumsafe/go
```

Maybe you want to use the Docker container as development environment.
Mount your current project in the Docker container with:

```bash
 docker run --rm -it --workdir=/app -v ${PWD}:/app openquantumsafe/go /bin/bash 
```

---

<a name="documentation"></a>Documentation
----

The `liboqs-go` wrapper is fully documented using the Go standard documentation
conventions. For example, to read the full documentation about
the `oqs.Signature.Verify` method, type in a terminal/console

	go doc $HOME/liboqs-go/oqs.Signature.Verify

For the RNG-related function, type e.g.

    go doc $HOME/liboqs-go/oqs/rand.RandomBytes

For automatically-generated documentation in HTML format,
click [here](https://pkg.go.dev/github.com/open-quantum-safe/liboqs-go/oqs).

For the RNG-related documentation,
click [here](https://pkg.go.dev/github.com/open-quantum-safe/liboqs-go/oqs/rand)
.

<a name="limitations"></a>Limitations and security
----

liboqs is designed for prototyping and evaluating quantum-resistant
cryptography. Security of proposed quantum-resistant algorithms may rapidly
change as research advances, and may ultimately be completely insecure against
either classical or quantum computers.

We believe that the NIST Post-Quantum Cryptography standardization project is
currently the best avenue to identifying potentially quantum-resistant
algorithms. liboqs does not intend to "pick winners", and we strongly recommend
that applications and protocols rely on the outcomes of the NIST standardization
project when deploying post-quantum cryptography.

We acknowledge that some parties may want to begin deploying post-quantum
cryptography prior to the conclusion of the NIST standardization project. We
strongly recommend that any attempts to do make use of so-called **hybrid
cryptography**, in which post-quantum public-key algorithms are used alongside
traditional public key algorithms (like RSA or elliptic curves) so that the
solution is at least no less secure than existing traditional cryptography.

Just like liboqs, liboqs-go is provided "as is", without warranty of any kind.
See [LICENSE](https://github.com/open-quantum-safe/liboqs-go/blob/main/LICENSE)
for the full disclaimer.

License
-------

liboqs-go is licensed under the MIT License;
see [LICENSE](https://github.com/open-quantum-safe/liboqs-go/blob/main/LICENSE)
for details.

Team
----

The Open Quantum Safe project is led
by [Douglas Stebila](https://www.douglas.stebila.ca/research/)
and [Michele Mosca](http://faculty.iqc.uwaterloo.ca/mmosca/) at the University
of Waterloo.

liboqs-go was developed by [Vlad Gheorghiu](http://vsoftco.github.io) at
softwareQ Inc. and University of Waterloo.

### Support

Financial support for the development of Open Quantum Safe has been provided by
Amazon Web Services and the Canadian Centre for Cyber Security.

We'd like to make a special acknowledgement to the companies who have dedicated
programmer time to contribute source code to OQS, including Amazon Web Services,
evolutionQ, softwareQ, and Microsoft Research.

Research projects which developed specific components of OQS have been supported
by various research grants, including funding from the Natural Sciences and
Engineering Research Council of Canada (NSERC); see the source papers for
funding acknowledgments.
