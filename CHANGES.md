# Version 0.10.0 - March 27, 2024

- Bumped Go version to 1.21
- Replaced CHANGES by
  [CHANGES.md](https://github.com/open-quantum-safe/liboqs-go/blob/main/CHANGES.md),
  as we now use Markdown format to keep track of changes in new releases
- Removed the NIST PRNG as the latter is no longer exposed by liboqs' public
  API
- Added the
  [.config-static](https://github.com/open-quantum-safe/liboqs-go/tree/main/.config-static)
  pkg-config configuration directory for linking statically against liboqs, see
  [README.md](https://github.com/open-quantum-safe/liboqs-go/blob/main/README.md)
  for more details

# Version 0.9.0 - October 30, 2023

- No modifications, release bumped to match the latest release of liboqs

# Version 0.8.0 - July 5, 2023

- This is a maintenance release, minor fixes
- Minimalistic Docker support
- Go minimum required version bumped to 1.15
- Removed AppVeyor and CircleCI, all continuous integration is now done via
  GitHub actions

# Version 0.7.2 - August 26, 2022

- Added liboqs library version retrieval function `LiboqsVersion() string`

# Version 0.7.1 - January 5, 2022

- Release numbering updated to match liboqs
- Switched continuous integration from Travis CI to CircleCI, we now support
  macOS & Linux (CircleCI) and Windows (AppVeyor)

# Version 0.4.0 - November 28, 2020

- Bugfixes
- Renamed 'master' branch to 'main'

# Version 0.3.0 - June 10, 2020

- Full Windows support and AppVeyor continuous integration
- Minor fixes

# Version 0.2.2 - December 10, 2019

- Changed panics to errors in the API

# Version 0.2.1 - November 7, 2019

- Added a client/server KEM over TCP/IP example

# Version 0.2.0 - November 2, 2019

- Minor API change to account for Go naming conventions
- Concurrent unit testing

# Version 0.1.2 - October 31, 2019

- Added support for RNGs from `<oqs/rand.h>`

# Version 0.1.1 - October 24, 2019

- Added support for Go modules

# Version 0.1.0 - October 22, 2019

- Initial release
