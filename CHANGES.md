# Changelog for liboqs-go

## Version 0.16.0 - July 15, 2026

- Updated compatibility for liboqs 0.16.0
- **Breaking change:** SPHINCS+ has been removed from liboqs 0.16.0. Users must migrate to SLH-DSA ([FIPS 205](https://csrc.nist.gov/pubs/fips/205/final))
- Added support for the ML-DSA external-mu variants `ML-DSA-44-extmu`, `ML-DSA-65-extmu`, and `ML-DSA-87-extmu`, which take a fixed 64-byte mu digest in place of the message
- Fixed a panic in `Sign`/`Verify`/`SignWithCtxStr`/`VerifyWithCtxStr` when passed an empty message, signature, or context string (zero-length slice dereference)
- liboqs 0.16.0 renames the existing FrodoKEM parameter sets to ephemeral FrodoKEM (`efrodokem_*`) and adds the salted variant under the original `frodokem_*` names
- Added an automated release workflow (`release.yml`) that is triggered by upstream liboqs and publishes a version-matched release once tests pass; see [RELEASING.md](https://github.com/open-quantum-safe/liboqs-go/blob/main/RELEASING.md)
- Made the CI workflow reusable, pinned all GitHub Actions to commit hashes, and added an actionlint + zizmor lint gate for the workflow files
- Switched release tags to the Go module `vX.Y.Z` convention (e.g. `v0.16.0`), so the module can be fetched with `go get github.com/open-quantum-safe/liboqs-go@v0.16.0`

## Version 0.15.0 - January 20, 2026

- Updated compatibility for liboqs 0.15.0
- **Breaking change:** Dilithium has been removed from liboqs 0.15.0.  Users must migrate to ML-DSA (FIPS 204)
- SPHINCS+ is still supported in liboqs 0.15.0 but will be removed in liboqs 0.16.0 and replaced by SLH-DSA
- liboqs 0.15.0 adds support for NTRU (re-added), updated CROSS to version 2.2, and includes SLH-DSA implementation
- Updated GitHub Actions workflow to increase the test timeout from the default 10 minutes to 30 minutes.

## Version 0.12.0 - January 15, 2025

- Fixes [Issue #44](https://github.com/open-quantum-safe/liboqs-go/issues/44). The API that NIST has introduced in [FIPS 204](https://csrc.nist.gov/pubs/fips/204/final) for ML-DSA includes a context string of length >= 0. Added new API for signing with a context string:

> `func (sig *Signature) SignWithCtxStr(message []byte, context []byte) ([]byte, error)`

> `func (sig *Signature) VerifyWithCtxStr(message []byte, signature []byte, context []byte, publicKey []byte) (bool, error)`

- Updated examples to use `ML-KEM` and `ML-DSA` as the defaults
- Removed the `oqs.rand` package and moved the `RandomBytes` family of functions from `oqs.rand` to the main `oqs` package to avoid warnings about linking liboqs twice

## Version 0.10.0 - March 27, 2024

- Bumped Go version to 1.21
- Replaced CHANGES by [CHANGES.md](https://github.com/open-quantum-safe/liboqs-go/blob/main/CHANGES.md), as we now use Markdown format to keep track of changes in new releases
- Removed the NIST PRNG as the latter is no longer exposed by liboqs' public API
- Added the [.config-static](https://github.com/open-quantum-safe/liboqs-go/tree/main/.config-static) pkg-config configuration directory for linking statically against liboqs, see [README.md](https://github.com/open-quantum-safe/liboqs-go/blob/main/README.md) for more details

## Version 0.9.0 - October 30, 2023

- No modifications, release bumped to match the latest release of liboqs

## Version 0.8.0 - July 5, 2023

- This is a maintenance release, minor fixes
- Minimalistic Docker support
- Go minimum required version bumped to 1.15
- Removed AppVeyor and CircleCI, all continuous integration is now done via GitHub actions

## Version 0.7.2 - August 26, 2022

- Added liboqs library version retrieval function `LiboqsVersion() string`

## Version 0.7.1 - January 5, 2022

- Release numbering updated to match liboqs
- Switched continuous integration from Travis CI to CircleCI, we now support macOS & Linux (CircleCI) and Windows (AppVeyor)

## Version 0.4.0 - November 28, 2020

- Bugfixes
- Renamed 'master' branch to 'main'

## Version 0.3.0 - June 10, 2020

- Full Windows support and AppVeyor continuous integration
- Minor fixes

## Version 0.2.2 - December 10, 2019

- Changed panics to errors in the API

## Version 0.2.1 - November 7, 2019

- Added a client/server KEM over TCP/IP example

## Version 0.2.0 - November 2, 2019

- Minor API change to account for Go naming conventions
- Concurrent unit testing

## Version 0.1.2 - October 31, 2019

- Added support for RNGs from `<oqs/rand.h>`

## Version 0.1.1 - October 24, 2019

- Added support for Go modules

## Version 0.1.0 - October 22, 2019

- Initial release
