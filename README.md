liboqs-go: Go bindings for liboqs
===================================

---

**liboqs-go** offers a Go wrapper for the master branch of [Open Quantum Safe](https://openquantumsafe.org/) [liboqs](https://github.com/open-quantum-safe/liboqs/) C library, which is a C library for quantum-resistant cryptographic algorithms.

The wrapper is written in Go, hence in the following it is assumed that you have access to a Go compliant environment. liboqs-go has been extensively tested on Linux, macOS and Windows systems. Continuous integration is provided via Travis CI and AppVeyor.

## Pre-requisites

liboqs-go depends on the [liboqs](https://github.com/open-quantum-safe/liboqs) C library; liboqs master branch must first be compiled as a Linux/macOS/Windows library, see the liboqs [specific platform building instructions](https://github.com/open-quantum-safe/liboqs#quickstart).

Contents
--------

liboqs-go is a Go package. The project contains the following files
and folders:

 - **`src/oqs/oqs.go`: main package file for the wrapper**
 - `examples/kem.go`: key encapsulation example
 - `examples/sig.go`: signature example
 - `unit_tests`: unit tests 

Usage
-----
The examples in the [`examples`](https://github.com/open-quantum-safe/liboqs-go/tree/master/examples) folder are self-explanatory and provide more details about the wrapper's API.

Building
--------------------------------------------

First, you must build the master branch of liboqs according to the [liboqs building instructions](https://github.com/open-quantum-safe/liboqs#building), followed (optionally) by a `sudo make install` to ensure that the compiled library is system-wide visible (by default it installs under `/usr/local/include` and `/usr/local/lib` under Linux/macOS).

Next, to use the wrapper, you must first add/append `liboqs-go` to the `GOPATH` environment variable, and add/append `liboqs-go/config` to the `PKG_CONFIG_PATH` environment variable.

Then simple `import "oqs"` in your Go program and run with `go run <program.go>` or build an executable with `go build <program.go>`. To run the examples from the terminal/command prompt, type (from the project root directory)

    go run examples/sig.go 
    
or 
    
    go run examples/kem.go

from the root of the project folder. Replace `go run` with `go build` if you intend to build the corresponding executables.

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
