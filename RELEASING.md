# Releasing liboqs-go

Releases track the liboqs version and publish automatically on green tests.

## Versioning

- Tags: `vMAJOR.MINOR.PATCH`, e.g. `v0.16.0` ([Go convention](https://go.dev/doc/modules/version-numbers); the `v` is required for `go get`).
- Version matches the liboqs release built/tested against.
- liboqs tags have no `v` (`0.16.0`); the workflow adds it.

## Triggers — `.github/workflows/release.yml`

- `repository_dispatch` type `liboqs-go-release`, `client_payload.version` = liboqs tag.
- `workflow_dispatch` with a `version` input.
- Push of a `v*` tag.

Flow:

- Reuses `go.yml` (Linux/macOS/Windows) against the matching liboqs version.
- Release candidates (`X.Y.Z-rcN`, e.g. from a liboqs RC) only run the tests — no tag, no release.
- On all-green: verifies main is prepped for the version (`.pc` files, `RELEASE.md`, a `## Version X.Y.Z` section in `CHANGES.md` — see below), then tags `vX.Y.Z` and publishes a release with the `CHANGES.md` notes.
- On failure (tests or prep check): nothing releases — fix and re-trigger.

## Upstream liboqs trigger

`release.yml` listens for a `liboqs-go-release` `repository_dispatch` (`client_payload.version` = liboqs release tag), sent from liboqs on `release: published`.

- Upstream liboqs PR: [open-quantum-safe/liboqs#2494](https://github.com/open-quantum-safe/liboqs/pull/2494)

## Before a release

- Add `## Version X.Y.Z - <date>` to `CHANGES.md`.
- Bump `Version:` in `.config/*.pc` and `.config-static/*.pc`.
- Update `RELEASE.md` (version, date, release URL).

`release.yml` enforces this checklist: it refuses to publish until all three match
the version being released, so an upstream dispatch against an unprepped main
fails loudly instead of tagging.
