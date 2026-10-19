# Agent Instructions

## Project references

- Before product or implementation work, read `README.md`, `docs/command-spec.md`, `docs/acceptance-criteria.md`, `docs/markdown-terminal-reader-analysis.md`, and `docs/development-workflow.md`.
- Treat `docs/command-spec.md` as the source of truth for command and TUI behavior.
- Treat `docs/acceptance-criteria.md` as the source of truth for expected behavior and verification.
- Follow `docs/development-workflow.md` for implementation phases and verification gates.
- Use one ready ticket from `docs/tickets/README.md` as the scope for each implementation session; follow its blockers and acceptance checklist.
- If project references conflict, ask the user which decision is current; do not silently choose or change product behavior.

## Communication and decisions

- Communicate in Bahasa Indonesia and explain unfamiliar technical terms in simple language.
- Keep user-facing README and CLI help/error text in English as specified by project docs; keep development explanations and internal docs in Indonesian.
- For unresolved product or design decisions, run a grill-me question round and wait until shared understanding is confirmed before dependent implementation.
- Do not implement features that are explicitly outside the MVP without user approval.

## Project constraints

- The implementation language is Go.
- MVP is a read-only TUI that opens one file path; stdin input and editing are out of scope.
- Preserve the user's document; rendering must not execute raw HTML or make network requests.
- Open HTTP(S) links only after explicit user activation; never open links automatically.
- Sanitize terminal control characters from document content and restore terminal state on exit or error.
- Keep parser, layout, and TUI responsibilities separate.

## Development checks

- Format Go code: `gofmt -w cmd internal`.
- Verify the repository: `go test ./...`, `go vet ./...`, and `go build ./...`.
- CI runs build, test, and vet on Windows, macOS, and Linux; its formatting check runs on macOS/Linux.

<!-- antislop:start -->
## antislop

For TUI visual or interaction work, read `.agents/skills/antislop/SKILL.md` and follow its mode, direction, and delivery requirements. This pointer targets the core skill already present in the repository.
<!-- antislop:end -->
