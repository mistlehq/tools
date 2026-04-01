## Scope

- This repository hosts open source provider-focused CLIs which are thin wrappers around provider APIs.
- Provider binaries live under `cmd/*`.

## Tooling

- This repo uses `mise` to pin the development toolchain. Run `mise trust ./mise.toml` once, then `mise install`.
- Prefer running commands through `mise exec -- ...` so the pinned Go toolchain is used consistently.
- Keep the toolchain minimal.
- Always use Conventional Commits such as `feat:`, `fix:`, `chore:`, `refactor:`, `test:`, and `docs:`.

## Dependencies & External APIs

- Avoid introducing dependencies as much as possible.
- If you need to add a dependency, prefer the most maintained and widely adopted option with a clear API and active release cadence.
- Before adding a dependency for an external provider, check the provider's official SDKs and official API documentation first.
- Prefer small, focused dependencies over large framework-style additions.

## Abstractions

- Avoid excessive abstractions that provide no immediate benefit
- If abstractions are needed, consider surfacing them first as suggestions instead

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

- Unit tests: colocate them with the Go package they cover as `*_test.go`.
- Integration tests: use Go's normal `_test.go` structure and keep them near the package they exercise unless a broader layout is clearly needed.

## Workflows

- Run these commands in this order when validating a change:
  - `mise exec -- go test ./...`
- Also run `mkdir -p dist && mise exec -- go build -o dist/<name> ./cmd/<name>` when your change affects entrypoints or binary packaging.
- Prefer small commits that keep related changes together.

## Pull Requests

- GH CLI is available if you need to open or update a PR.
- Before opening a PR, rebase onto the latest `main`.
- If you open a PR, stay with it until CI is green unless the remaining failure requires human intervention.

## Go Guidance

- Prefer the standard library first.
- Keep packages small and package-centric rather than recreating TypeScript-style folder layering.
- Start flat. Extract shared packages only when a second binary actually needs the same behavior.
- Parse arguments explicitly and fail with clear error messages on invalid input.
