# DevOrchestra architecture (MVP)

## Design constraints

- Local-first Go application; opens the user's default browser at a loopback-only address.
- The HTTP backend serves embedded offline HTML/CSS/JavaScript with no frontend toolchain.
- Persistent state uses SQLite (modernc.org/sqlite, pure Go, no CGO).
- Codex is always a separate CLI child process; the application never handles a ChatGPT password or token.
- Requests use `codex exec --json` and stdin prompts, with an explicit read-only or workspace-write sandbox.
- Each job persists JSONL events and stderr to SQLite. Only one job at a time in the first iteration.
- API requires a random per-launch token in a custom HTTP header. No remote listener.
- External prompt imports are manually initiated, HTTPS-only and restricted to raw.githubusercontent.com Markdown/text files. Imported text is untrusted until reviewed.
- No arbitrary shell execution endpoint. Never interpolate user text into a shell command.
- A native Wails shell is a planned optional future enhancement; the MVP uses the system browser to keep cross-platform builds lightweight.

## Storage

Application data lives in the OS user's configuration directory, under DevOrchestra:
`devorchestra.db`. This location is not the checked-out project tree.

Tables:
- `projects`: local working directories.
- `templates`: prompt templates, including source provenance.
- `runs`: prompt, project, sandbox, status, timestamps and Codex thread ID.
- `run_events`: chronological JSONL/stderr events by run.

On restart, jobs left in the running state are marked interrupted.

## Security notes

- localhost does NOT by itself prevent local browser attacks. Token-check every API, pin HTTP Host, and avoid permissive CORS.
- Treat stdout/stderr and fetched prompt text as potentially sensitive or malicious.
- The UI uses textContent, never innerHTML, for untrusted output.
- Imported prompts are never executed automatically.
- Do not silently set danger-full-access, disable sandboxing, or auto-approve changes.
- No credentials are stored by DevOrchestra; users authenticate via `codex login`.
- Logs may contain proprietary code or secrets: user must control retention and backup.
- Installer's optional npm install requires Node/npm already present and the user's explicit selection.
- Signing/notarization is required before broad production distribution.
