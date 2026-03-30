# mistle-tools

Open source provider CLIs and future skills for Mistle.

## Stack

- TypeScript source
- Bun workspaces
- `turbo` task orchestration
- `changesets` for versioning and changelogs
- `oxfmt` for formatting
- `oxlint` for linting
- `knip`, `jscpd`, and `typos` for repository hygiene
- `tsgo` for typechecking
- Bun for package management and standalone binary builds

## Layout

- `packages/core`: shared CLI shell and common helpers
- `packages/*`: provider CLIs and shared libraries

## Commands

```sh
mise trust ./mise.toml
mise install
bun install
bun run lint
bun run typecheck
bun run test
bun run build
```

## Current Packages

- `@mistle-tools/core`: shared CLI shell and metadata helpers
- `@mistle-tools/jira`: placeholder Jira CLI scaffold

## Binary Builds

Standalone binaries are produced with Bun directly from the package entrypoints.
The source code remains plain TypeScript and the workspace remains Bun-managed.

Typechecking uses `tsgo`. The `build` task emits the final standalone binaries.

## Toolchain

`mise` pins the external toolchain for local development:

- Bun

Trust the repo-local `mise.toml` once, then `mise install` will provision the
toolchain versions expected by the workspace.

## Git Hooks

`bun install` installs local Git hooks via `lefthook`.

- `pre-commit`: runs `lint-staged`
- `commit-msg`: validates Conventional Commits with `commitlint`
- `pre-push`: runs `bun run lint`, `bun run typecheck`, and `bun run test`
