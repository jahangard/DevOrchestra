# DevOrchestra User Guide — English

Early MVP 0.1.x | Windows, Linux, macOS

DevOrchestra is a local graphical companion for Codex CLI. Its Go backend opens the interface in your default web browser, and SQLite stores projects, prompts, and execution history on your computer.

> Current limitations: There is no standalone native desktop window, multi-agent execution, automated web prompt discovery, or AI-model-based prompt rewriting yet.

## 1. Install and launch

1. Open [GitHub Actions](https://github.com/jahangard/DevOrchestra/actions) and choose a successful Build and test workflow run.
2. Under **Artifacts**, download the package for your system: Windows (DevOrchestra-Setup.exe), Linux (deb) or macOS Apple Silicon (pkg).
3. Install, then launch DevOrchestra. Your default browser opens its local UI.

**Note:** Packages are not yet code-signed or Apple-notarized, and installation on real devices still needs separate testing.

### Run from source

Go 1.23 or newer is required:

    git clone https://github.com/jahangard/DevOrchestra.git
    cd DevOrchestra
    git switch feat/codex-desktop-mvp
    go mod tidy
    go run ./cmd/devorchestra

To quit, stop the process (Ctrl+C if using a terminal).

## 2. Set up Codex CLI

Go to **Settings** to inspect Codex CLI, Node.js and npm status.

- If Node.js and npm are already installed, choose **Install Codex CLI with npm** and confirm.
- The Windows installer has an optional Codex installation component, which also requires npm to be present first.
- Manual method:

    npm install -g @openai/codex
    codex login

Run codex login in a terminal and complete authentication. DevOrchestra does not request or store your account password or authentication token. Click **Refresh** in Settings afterward.

## 3. Register a project

1. Open **Projects**.
2. Enter a name and the absolute path of an *existing* directory, e.g. D:\dev\TavanHub or /home/user/project.
3. Click **Add**.
4. Select the project in the Dashboard.

The app does not create project folders.

## 4. Run a request

1. Select a project in **Dashboard**.
2. Choose a sandbox policy: **read-only** (the default) for inspection without write access, or **workspace-write** to permit project edits.
3. Describe the task in the prompt editor. **Structure prompt** appends context, testing and acceptance-criteria sections from a local template; it does *not* invoke another AI model.
4. Click **Run Codex** to execute and view events and output. **Cancel** stops an active run.

The MVP only supports **one run at a time**. Commit your existing changes before granting write access.

Example prompt:

    Review AGENTS.md and the existing architecture first.
    Reproduce the customer table filtering bug.
    Change only relevant files and add regression tests.
    Report changed files and test results.

## 5. Run history and logs

Open **Run history** to see recent prompts and statuses. Click **View result** to show execution events in Dashboard. JSONL events, errors and run metadata are retained in SQLite. Find the DB location under **Settings → SQLite**.

**Privacy:** Logs and prompts can contain proprietary source code, file paths and secrets. Review the database before sharing it. Automatic log retention management and database encryption are not included in this MVP.

## 6. Prompt library

- Open **Prompt library** to pick review, debug, feature, test or security starter templates. Click **Use** to load one.
- In Dashboard choose **Save prompt** to keep a reusable prompt.
- To import a prompt from GitHub, supply a direct md or txt file URL hosted on raw.githubusercontent.com, then choose **Import**.

External prompt content is **untrusted** and is never executed automatically. Inspect its source and instructions before use. Automatic discovery or aggregation from reference websites is not available yet.

## 7. Language

Use the **English / فارسی** button in the top bar to change the UI and the embedded Help page together. Persian is right-to-left; English is left-to-right. Your choice is saved in the browser.

## 8. Troubleshooting

| Problem | Suggested fix |
| --- | --- |
| Codex is missing | Install Node.js/npm first, then install Codex in Settings. |
| Not authenticated | Run codex login in a terminal and refresh Settings. |
| Project path rejected | Use the absolute path to a real, accessible directory. |
| Run cannot start | Check Codex status, authentication and whether a run is already active. |
| Browser does not open | Open the loopback URL printed in the terminal. |
| No files modified | Check the workspace-write sandbox setting and folder permissions. |

## 9. Developer and contact

**Developed by: مهدی جهانگرد — mehdi jahangard**

Email: [mehdi.jahangard@gmail.com](mailto:mehdi.jahangard@gmail.com)

Mobile: [+98 912 724 6727](tel:+989127246727)

Repository: https://github.com/jahangard/DevOrchestra
