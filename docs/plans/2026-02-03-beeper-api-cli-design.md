# Beeper API CLI Enhancement Plan — Design

Date: 2026-02-03

## Summary
Implement the full enhancement plan across messaging, desktop control, downloads, chat management, search/filters, streaming modes, reminders, safety/agent features, contacts, uploads, and optional advanced backends. Preserve backwards compatibility and LLM-friendly defaults.

## Goals
- Deliver all features in `ENHANCEMENT_PLAN.md` and gaps in `docs/LIBRARY_COMPARISON.md`.
- Keep existing UX stable; add new behavior via flags/subcommands.
- Maintain robust error handling and LLM-friendly output.
- Ensure live API integration works with Beeper Desktop.

## Non-Goals
- Breaking CLI changes to existing command signatures.
- Full SQLite offline implementation unless the Beeper Desktop DB path and schema are stable.

## Architecture
- `cmd/`: CLI surface with Cobra commands and flag validation.
- `internal/api`: HTTP client and models for new endpoints.
- `internal/output`: output formats (json/text/markdown/jsonl/template) and richer message rendering.
- `internal/config`: config updates (aliases, safety flags, auth extensions).
- New helpers in `internal/` for date parsing, IO input, downloads/uploads, safety checks.

## Features
### Messaging Enhancements
- `send --reply-to` adds `m.relates_to` with reply metadata.
- `send --stdin` and `send --file` are mutually exclusive with `--message`.
- Messages output includes reply context where available.

### Desktop Focus/Control
- New `focus` command using Desktop API focus endpoint.
- Supports `--chat-id`, `--message-id`, `--draft`, `--draft-stdin`, `--draft-file`, and attachments.

### Attachments
- New `download` command for MXC URLs.
- `messages list`/`search` add `--download-attachments` and `--download-dir`.
- Streaming downloads with retries, resuming, collision handling.

### Chat Management
- `chats archive` with `--unarchive`.
- `chats archive-read` with `--dry-run`/`--confirm`.
- `chats create` for DM/group with `--message` optional.
- `alias` command group with local storage and global resolution.

### Search & Filtering
- `--after`/`--before` using ISO or natural language.
- `--sender`, `--media-type`, `--chat-type` filters.
- `--context` and `--context-window` for surrounding messages.
- `--highlight` for text output.

### Streaming & Watch
- `messages tail` polling-based follow with filters, `--interval`.
- `messages wait` with exit code 0 on match, 1 on timeout.
- `inbox watch` for unread messages across chats.

### Reminders & Scheduling
- `reminders` group for set/list/clear.
- Local persistence and triggering behaviors (notify/focus/draft).

### Safety & Output Enhancements
- `jsonl` output for streaming.
- Template formatting via `--format` (Go templates).
- `--readonly` and command allowlist for agent safety.
- `doctor` diagnostics and `status` summary.

### Contacts & Uploads
- `contacts search` if API endpoint is available.
- `upload` or `send --attach` if Desktop API supports uploads.

### Optional Advanced Features
- SQLite read-only backend under `--source db` if stable.
- OAuth/token convenience behind a feature flag.

## Data Flow
Commands validate flags, resolve aliases, apply safety restrictions, then call API client. Output is formatted by `internal/output` and can include filters and annotations. If server-side filtering is unavailable, perform client-side filtering with explicit output notes.

## Error Handling
- Structured error JSON preserved.
- Safety blocks return clear error messages and nonzero exit codes.
- Parsing errors provide examples and expected formats.

## Testing Strategy
- Unit tests for date parsing, alias resolution, IO input, JSONL/template formatting, download collision logic, safety allowlist.
- API client tests for new endpoints against live API.
- Command tests for flag validation and output formats.
- Integration tests via `BEEPER_API_URL`, `BEEPER_TOKEN`, and `BEEPER_TEST_CHAT_ID`.

## Rollout
Implement in phases with incremental documentation updates to `README.md`, `API.md`, and `QUICKSTART.md`. Each phase lands with tests and compatibility checks.
