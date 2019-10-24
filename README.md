liboqs-go: Go bindings for liboqs
=================================

[![Build status - Linux/macOS](https://api.travis-ci.com/open-quantum-safe/liboqs-go.svg?branch=master)](https://travis-ci.com/open-quantum-safe/liboqs-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-quantum-safe/liboqs-go)](https://goreportcard.com/report/github.com/open-quantum-safe/liboqs-go)
[![Documentation](https://godoc.org/github.com/open-quantum-safe/liboqs-go/oqs?status.svg)](https://godoc.org/github.com/open-quantum-safe/liboqs-go/oqs)

---

**liboqs-go** offers a Go wrapper for the master branch of [Open Quantum Safe](https://openquantumsafe.org/) [liboqs](https://github.com/open-quantum-safe/liboqs/) C library, which is a C library for quantum-resistant cryptographic algorithms.

The wrapper is written in Go, hence in the following it is assumed that you have access to a Go compliant environment. liboqs-go has been extensively tested on Linux and macOS systems. Continuous integration is provided via Travis CI.

## Pre-requisites

liboqs-go depends on the [liboqs](https://github.com/open-quantum-safe/liboqs) C library; liboqs master branch must first be compiled as a UNIX/Linux/macOS library, see the liboqs [specific platform building instructions](https://github.com/open-quantum-safe/liboqs#quickstart).

In addition, we assume you have access to:

- a POSIX compliant system (UNIX/Linux/macOS). For now, `cgo` is not fully supported under Windows due to various ABI issues; we will add Windows support when it becomes available.
- Go version 1.7 or later
- a standard C compliant compiler (`gcc`/`clang` etc.)
- `pkg-config` (use `sudo apt-get install pkg-config` to install on Ubuntu/Debian-based Linux platforms)

Contents
--------

liboqs-go is a Go package. The project contains the following files and directories:

 - **`oqs/oqs.go`: main package file for the wrapper**
 - `.config/liboqs.pc`: `pkg-config` configuration file needed by `cgo`
 - `examples/kem.go`: key encapsulation example
 - `examples/sig.go`: signature example
 - `oqstests`: unit tests

Usage
-----

The examples in the [`examples`](https://github.com/open-quantum-safe/liboqs-go/tree/master/examples) directory are self-explanatory and provide more details about the wrapper's API.

Running/building
----------------
First, you must build the master branch of liboqs according to the [liboqs building instructions](https://github.com/open-quantum-safe/liboqs#building), followed (optionally) by a `sudo make install` to ensure that the compiled library is system-wide visible (by default it installs under `/usr/local/include` and `/usr/local/lib` under Linux/macOS).

### Using Go modules (requires Go 1.11 or later)
Download/clone the `liboqs-go` wrapper repository in the directory of your choice, e.g. `$HOME`, by typing in a terminal/console

    cd $HOME && git clone https://github.com/open-quantum-safe/liboqs-go
    
Next, you must `cd $HOME/liboqs-go` and modify the following lines in [`.config/liboqs.pc`](https://github.com/open-quantum-safe/liboqs-go/tree/master/.config/liboqs.pc)

    LIBOQS_INCLUDE_DIR=/usr/local/include
    LIBOQS_LIB_DIR=/usr/local/lib

so they correspond to your C liboqs include/lib installation directories.

Finally, you must add/append the `.config` directory to the `PKG_CONFIG_PATH` environment variable, i.e.

    export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$HOME/liboqs-go/.config
    
Once you have configured your system as directed above, simply `import "github.com/open-quantum-safe/liboqs-go/oqs"` in the Go application of your choice, initialize the application module with `go mod init <module_name>`, and finally run it with `go run <module_name>` or build it with `go build <module_name>`.

To run the examples from the terminal/console, type (from `$HOME/liboqs-go`)

    go run examples/kem.go 
    
or 
    
    go run examples/sig.go

Replace `go run` with `go build` if you intend to build the corresponding executables; in this case they will be built in the directory you ran the `go build` command from. 

To run the unit tests from the terminal/console, type (from `$HOME/liboqs-go`)
	
	go test -v ./oqstests
	
To build the unit test executable from the terminal/console, type (from the directory in which you want to build the executable)

    go test -c /path/to/liboqs-go/oqstests
    
This will build the `oqstests.test` executable in the directory of your choice above.


### Using the traditional Go get
Install the latest version of the `liboqs-go` wrapper by typing 

    go get github.com/open-quantum-safe/liboqs-go/oqs

in a terminal/console. This will install the wrapper in one of your `$GOPATH` directories. In my case the Go package manager installs the wrapper in `$HOME/go/src/github.com/open-quantum-safe/liboqs-go`. To update a previously installed Go wrapper, type

    go get -u github.com/open-quantum-safe/liboqs-go/oqs
    
in a terminal/console.

To simplify the instructions to follow, export the path to the wrapper in the `LIBOQSGO_INSTALL_PATH` environment variable by typing in a terminal/console

    export LIBOQSGO_INSTALL_PATH=/your/path/to/liboqs-go
    
In my case `LIBOQSGO_INSTALL_PATH` is set to `$HOME/go/src/github.com/open-quantum-safe/liboqs-go`.

Next, you must `cd $LIBOQSGO_INSTALL_PATH` and modify the following lines in [`.config/liboqs.pc`](https://github.com/open-quantum-safe/liboqs-go/tree/master/.config/liboqs.pc)

    LIBOQS_INCLUDE_DIR=/usr/local/include
    LIBOQS_LIB_DIR=/usr/local/lib

so they correspond to your C liboqs include/lib installation directories.

Finally, you must add/append the `$LIBOQSGO_INSTALL_PATH/.config` directory to the `PKG_CONFIG_PATH` environment variable, i.e.

    export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$LIBOQSGO_INSTALL_PATH/.config

If running/building on Linux, you may need to set the `LD_LIBRARY_PATH` environment variable to point to the path
to liboqs' library directory, i.e.

    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
            
assuming `liboqs.so.*` were installed in `/usr/local/lib` (true assuming you ran `make install` during your liboqs setup).
 
Once you have configured your system as directed above, simply `import "github.com/open-quantum-safe/liboqs-go/oqs"` in the Go application of your choice and run it with `go run <application_name.go>` or build it with `go build <applicaiton_name.go>`.

To run the examples from the terminal/console, type 

    go run $LIBOQSGO_INSTALL_PATH/examples/kem.go 
    
or 
    
    go run $LIBOQSGO_INSTALL_PATH/examples/sig.go

Replace `go run` with `go build` if you intend to build the corresponding executables; in this case they will be built in the directory you ran the `go build` command from. 

To run the unit tests from the terminal/console, type
	
	go test -v github.com/open-quantum-safe/liboqs-go/oqstests
	
To build the unit test executable from the terminal/console, type (from the directory in which you want to build the executable)

    go test -c $LIBOQSGO_INSTALL_PATH/oqstests/*
    
This will build the `oqstests.test` executable in the directory of your choice above.

Documentation
-------------

The `oqs` Go package is fully documented using the Go standard documentation conventions. For example, to read the full documentation about the `oqs.Signature.Verify` method, type in a terminal/console

    go doc github.com/open-quantum-safe/liboqs-go/oqs.Signature.Verify
 
For [GoDoc](https://godoc.org) automatically-generated documentation in HTML format, click [here](https://godoc.org/github.com/open-quantum-safe/liboqs-go/oqs).

Limitations and security
------------------------

liboqs is designed for prototyping and evaluating quantum-resistant cryptography. Security of proposed quantum-resistant algorithms may rapidly change as research advances, and may ultimately be completely insecure against either classical or quantum computers.

We believe that the NIST Post-Quantum Cryptography standardization project is currently the best avenue to identifying potentially quantum-resistant algorithms. liboqs does not intend to "pick winners", and we strongly recommend that applications and protocols rely on the outcomes of the NIST standardization project when deploying post-quantum cryptography.

We acknowledge that some parties may want to begin deploying post-quantum cryptography prior to the conclusion of the NIST standardization project. We strongly recommend that any attempts to do make use of so-called **hybrid cryptography**, in which post-quantum public-key algorithms are used alongside traditional public key algorithms (like RSA or elliptic curves) so that the solution is at least no less secure than existing traditional cryptography.

Just like liboqs, liboqs-go is provided "as is", without warranty of any kind. See [LICENSE](https://github.com/open-quantum-safe/liboqs-go/blob/master/LICENSE) for the full disclaimer.

License
-------

liboqs-go is licensed under the MIT License; see [LICENSE](https://github.com/open-quantum-safe/liboqs-go/blob/master/LICENSE) for details.

Team
----

The Open Quantum Safe project is led by [Douglas Stebila](https://www.douglas.stebila.ca/research/) and [Michele Mosca](http://faculty.iqc.uwaterloo.ca/mmosca/) at the University of Waterloo.

liboqs-go was developed by [Vlad Gheorghiu](http://vsoftco.github.io) at evolutionQ and University of Waterloo.

### Support

Financial support for the development of Open Quantum Safe has been provided by Amazon Web Services and the Tutte Institute for Mathematics and Computing.  

We'd like to make a special acknowledgement to the companies who have dedicated programmer time to contribute source code to OQS, including Amazon Web Services, evolutionQ, and Microsoft Research.  

Research projects which developed specific components of OQS have been supported by various research grants, including funding from the Natural Sciences and Engineering Research Council of Canada (NSERC); see the source papers for funding acknowledgments.
