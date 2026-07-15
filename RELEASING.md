# Releasing liboqs-go

Releases track the liboqs version and publish automatically on green tests.

## Versioning

- Tags: `vMAJOR.MINOR.PATCH`, e.g. `v0.16.0` ([Go convention](https://go.dev/doc/modules/version-numbers); the `v` is required for `go get`).
- Version matches the liboqs release built/tested against.
- liboqs tags have no `v` (`0.16.0`); the workflow adds it.

## Triggers — `.github/workflows/release.yml`

- `repository_dispatch` type `liboqs-release`, `client_payload.version` = liboqs tag.
- `workflow_dispatch` with a `version` input.
- Push of a `v*` tag.

Flow:

- Reuses `go.yml` (Linux/macOS/Windows) against the matching liboqs version.
- On all-green: tags `vX.Y.Z`, publishes a release with notes from the `## Version X.Y.Z` section of `CHANGES.md` (auto-generated if absent).
- On failure: nothing releases — fix and re-trigger.

## Upstream liboqs trigger

liboqs dispatches a `liboqs-release` `repository_dispatch` (payload `version` = the release tag) to this repo on `release: published`, from its `downstream-release.yml`.

- Upstream liboqs PR: <!-- TODO: link liboqs PR -->

## Before a release

- Add `## Version X.Y.Z - <date>` to `CHANGES.md`.
- Bump `Version:` in `.config/*.pc` and `.config-static/*.pc`.
- Update `RELEASE.md` (version, date, release URL).
