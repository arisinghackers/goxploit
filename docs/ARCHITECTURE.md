# Architecture

## Package Boundaries

- `pkg/msfrpc`
  - Low-level RPC transport and authentication.
  - Owns request encoding/decoding and HTTP behavior.
- `pkg/msfrpc/generated`
  - Generator-owned wrappers.
  - Map-based responses (`map[string]interface{}`).
  - Must not be edited manually.
- `pkg/metasploit`
  - Typed, stable SDK for application code.
  - Adapts low-level RPC responses into typed models.
- `internal/generator`
  - Scraper and code generation implementation.
- `cmd/generator`
  - CLI command used by contributors and CI for regeneration.

## Design Rules

- Runtime code must not depend on scraping docs.
- Typed SDK changes go in `pkg/metasploit`, not in generated files.
- Generated files can be overwritten at any time.
- Add typed slices endpoint-by-endpoint with tests.

## Workflow

1. Update generator internals in `internal/generator` if generation behavior changes.
2. Run `make generate` to refresh `pkg/msfrpc/generated`.
3. Add typed adapters in `pkg/metasploit`.
4. Run `make check`.
