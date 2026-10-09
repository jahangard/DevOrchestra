# DevOrchestra

**A local-first, bilingual companion for OpenAI Codex CLI.**

A lightweight Go web application that runs locally on Windows, macOS and Linux. It opens an offline web UI in your default browser, manages multiple local workspaces, sends requests to an existing Codex CLI installation, and records structured execution history in SQLite.

> **Status:** early MVP on the `feat/codex-desktop-mvp` branch. Cross-platform builds and installers must be tested before a stable release. This is not yet a multi-agent orchestrator.

## User guides / راهنمای استفاده

- **فارسی:** [راهنمای کامل نصب و استفاده](docs/USER_GUIDE.fa.md)
- **English:** [Complete installation and usage guide](docs/USER_GUIDE.en.md)

Inside the application, open **Help / راهنمای استفاده** from the sidebar. Both the UI and embedded guide switch between Persian (RTL) and English (LTR) using the language button. Developer contact details are available under **About / درباره برنامه**.

## Current features

- Persian (RTL) and English (LTR) interface, with no frontend build dependencies.
- Project registration using local filesystem paths.
- Codex CLI detection, account sign-in status, and opt-in installation using npm.
- Non-interactive Codex execution with JSONL event history, stdout/stderr logging, job status and cancellation.
- Read-only sandbox by default, with explicit workspace-write selection.
- Prompt starter templates, save/reuse, and cost-free structured prompt enhancement.
- User-triggered import of Markdown/text templates from raw.githubusercontent.com, retaining source attribution.
- Local SQLite storage, including run history across app restarts.
- Installer recipes for Windows (NSIS), macOS (pkg), and Linux (deb); GitHub Actions builds each OS.

## Quick start

Prerequisites: Go 1.23+ and network access for Go modules. Node.js/npm is needed to install Codex through npm.

```bash
git clone https://github.com/jahangard/DevOrchestra.git
cd DevOrchestra
git switch feat/codex-desktop-mvp
go mod tidy
go run ./cmd/devorchestra
```

The app opens a random loopback-only address in the default browser. Do not expose the local server on a public network.

For Codex CLI installation, select **Settings → Install Codex CLI**, or use:

```bash
npm install -g @openai/codex
codex login
```

If you already have Codex CLI, DevOrchestra uses that installation and its existing authentication. It never asks for or stores your login credentials.

## Building

```bash
go test ./...
go build -trimpath -o devorchestra ./cmd/devorchestra
```

On Windows add `-ldflags "-H=windowsgui"` to suppress the console, if preferred. The program needs a web browser but does **not** require a running external web server or Node.js once Codex itself is installed.

### Packaging

- **Windows:** Build `dist/devorchestra-windows-amd64.exe` then run `makensis packaging/windows/setup.nsi` from repository root. NSIS offers an optional Codex installation component (requires preinstalled npm). The installer registers an uninstaller.
- **macOS:** Build `dist/devorchestra-darwin-arm64` on a compatible host, then run `bash packaging/macos/build-pkg.sh dist/devorchestra-darwin-arm64`.
- **Linux:** Build `dist/devorchestra-linux-amd64` then run `bash packaging/linux/build-deb.sh dist/devorchestra-linux-amd64`.

The GitHub Actions workflow produces downloadable build artifacts after successful CI runs. Packages are not signed or notarized yet.

## Security and privacy

Codex operates in the selected local project folder. Default sandbox is `read-only`; `workspace-write` requires explicit selection. There is no danger-full-access option in the GUI. The API binds to 127.0.0.1 only and requires a random per-launch request token and strict Host/Origin checks. Prompts and logs remain on the user's computer in their OS configuration directory. Those logs **may contain private code or secrets**: protect local backups and avoid sharing the SQLite file.

Imported prompts are untrusted documents; they are never executed automatically. Review external instructions and their licenses before using them. This MVP limits imports to user-specified GitHub raw Markdown/text URLs; automatic discovery and safe source refresh are on the roadmap.

## Documentation

- [Architecture & threat model](docs/ARCHITECTURE.md)
- [Roadmap](docs/ROADMAP.md)

## Developer / برنامه‌نویس

**برنامه نویسی : مهدی جهانگرد mehdi jahangard**

- Email: [mehdi.jahangard@gmail.com](mailto:mehdi.jahangard@gmail.com)
- Mobile: [+98 912 724 6727](tel:+989127246727)

## License

No license has been selected yet; public visibility does not grant blanket rights to redistribute this code.
