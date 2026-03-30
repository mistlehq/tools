# General

## Scope

- This repository hosts open source provider-focused CLIs and shared packages for Mistle.
- Provider CLIs live directly in `packages/*`. Do not create package names like `foo-cli`; the CLI nature is already implied by the repository.
- `packages/core` contains shared CLI infrastructure.

## Tooling

- This repo uses `mise` to pin the development toolchain. Run `mise trust ./mise.toml` once, then `mise install`.
- Prefer running workspace commands through `mise exec -- ...` so the pinned Bun toolchain is used consistently.
- The canonical package manager is Bun.
- `oxfmt` is the formatter and `oxlint` is the linter. Do not add Biome, ESLint, or Prettier unless explicitly requested.
- `knip`, `jscpd`, and `typos` are part of the repository linting stack.
- Bun is used for standalone binary builds. TypeScript remains the source language.
- Prefer `tsgo` for typechecking and package builds. Use project references and build mode for workspace packages.
- Git hooks are managed with `lefthook`.
- Conventional Commits are enforced with `commitlint`.

## Dependencies & External APIs

- If you need to add a dependency, prefer the most maintained and widely adopted option with a clear API and active release cadence.
- Before adding a dependency for an external provider, check the provider's official SDKs and official API documentation first.
- Prefer small, focused dependencies over large framework-style additions unless the user explicitly wants a broader stack change.

## Fallback Behavior

- Do not write fallback behavior unless the user explicitly asks for fallback behavior in this task.
- Fail fast with explicit errors when required configuration, provider state, or runtime assumptions are missing.
- Do not silently switch auth modes, endpoints, output formats, or packaging paths.
- If a fallback is explicitly approved, make it obvious in code and cover it with explicit tests.

## Testing Philosophy

- Strict rule: do not use mocking, stubbing, faking, or simulated behavior in tests.
- Disallowed mocking APIs include `vi.fn`, `vi.spyOn`, `vi.mock`, `jest.*`, `sinon`, `nock`, `msw`, and equivalent libraries.
- Disallowed manual doubles include `Fake*`, `Stub*`, `Noop*`, and test-only implementations that do not match production behavior.
- Assert observable behavior instead of call counts or interaction patterns.
- Do not use fake timers or patched global time.
- Prefer pure unit tests for pure logic and real integration tests for dependency-bearing behavior.
- Test everything that materially changes behavior. For bug fixes, add a regression test.
- Unless the user asks otherwise, run only the tests relevant to the files you changed.

### Test Layout

- Unit tests: colocate as `src/**/*.test.ts`.
- Property tests: name them `*.property.test.ts` and colocate them with the unit-tested module.
- Integration tests: place them in `packages/*/integration/`.
- End-to-end or system-style CLI tests should exercise the real binary or spawned CLI process rather than internal helpers when practical.

### Property-Based Testing

- Use `fast-check` with Vitest via `@fast-check/vitest`.
- Keep property tests deterministic and replayable.
- Use explicit generator bounds and avoid heavy `.filter(...)` chains.
- Assert meaningful invariants such as round-trips, idempotence, canonicalization, stable ordering, or non-mutation.

## Workflows

- Run these commands in this order when validating a change:
  - `mise exec -- bun run format`
  - `mise exec -- bun run lint`
  - `mise exec -- bun run typecheck`
  - `mise exec -- bun run test`
- Also run `mise exec -- bun run build` when your change affects packaging, entrypoints, or release artifacts.
- Do not use `--no-verify` for commits or pushes.
- Always use Conventional Commits such as `feat:`, `fix:`, `chore:`, `refactor:`, `test:`, and `docs:`.
- Prefer small commits that keep related changes together.

## Releases

- This repo uses Changesets for versioning and changelog management.
- If a change affects a public package or binary in a user-visible way, add a changeset unless the user explicitly says not to.
- Do not bundle unrelated release notes into a single changeset.

## Pull Requests

- GH CLI is available if you need to open or update a PR.
- Before opening a PR, rebase onto the latest `main`.
- If you open a PR, stay with it until CI is green unless the remaining failure requires human intervention.

## Language Guidance

### TypeScript

- `any` and `as` are forbidden.
- Check `node_modules` for real external API types instead of guessing.
- Never use inline imports or dynamic imports for types. Use standard top-level imports.
- Do not remove intentional functionality just to satisfy types or tooling. Upgrade or adapt the dependency instead.
- Avoid immediately invoked function expressions. Prefer module scope or named functions.
- Avoid unnecessary inline closures when a named function is clearer.

### CLI Design

- Keep CLIs thin. Put reusable logic in shared packages rather than mixing common infrastructure into every provider package.
- Put universal CLI behavior such as `help`, `version`, shared output framing, and common process handling in `packages/core`.
- Parse arguments explicitly and fail with clear error messages on invalid input.
- Prefer machine-readable output modes only when they are intentional and specified by the command contract.
- Do not hide provider-specific constraints. Surface them directly in help text, errors, and docs.

### Provider Integrations

- Prefer official provider terminology for auth, scopes, resources, and identifiers.
- Do not assume one auth strategy fits all providers. Model provider-specific auth expectations explicitly.
- When supporting multiple auth modes for one provider, make the mode selection explicit rather than inferred from loosely shaped input.
- When a CLI is intended to run behind Mistle's credentialless proxy model, keep credential resolution and auth header injection outside the provider package.
- Keep config resolution opaque to command handlers. Use runtime adapters at the package boundary instead of loading env or files directly inside command modules.
