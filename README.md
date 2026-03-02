# BrowserAuditor

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-blue)](#supported-platforms)
[![Release](https://img.shields.io/github/v/release/r4j3sh-com/BrowserAuditor?display_name=tag)](https://github.com/r4j3sh-com/BrowserAuditor/releases)
[![Repo](https://img.shields.io/badge/GitHub-BrowserAuditor-181717?logo=github)](https://github.com/r4j3sh-com/BrowserAuditor)

BrowserAuditor is a local-only browser profile inventory and system audit tool written in Go. It helps security teams, incident responders, and recovery workflows identify browser profile locations, artifact files, and host metadata without extracting secrets or sending data over the network.

Repository: `https://github.com/r4j3sh-com/BrowserAuditor.git`

## Description

BrowserAuditor scans supported browser profile directories and generates a structured JSON report with browser profile paths, artifact file locations, and host metadata. It is designed for browser artifact discovery, recovery preparation, forensic triage, and local system auditing while keeping collection limited to metadata and file presence.

## Why BrowserAuditor

- Local-only JSON output
- No network transmission or cloud upload
- No cookie decryption
- No password extraction
- No browsing history parsing
- Cross-platform browser path discovery for macOS, Linux, and Windows
- Useful for browser profile inventory, recovery planning, and forensic triage

## Key Features

- Enumerates browser profile directories for Chrome, Brave, Chromium, and Firefox
- Detects known browser artifact files such as `Bookmarks`, `History`, `Cookies`, `Login Data`, `places.sqlite`, and `logins.json`
- Captures host metadata including OS, architecture, hostname, username, home directory, and local IP addresses
- Writes timestamped JSON output by default so reports are not overwritten
- Supports custom output paths with `-out`
- Uses only the Go standard library

## Supported Browsers

- Google Chrome
- Brave Browser
- Chromium
- Mozilla Firefox

## Supported Platforms

- macOS
- Linux
- Windows

## Security Model

BrowserAuditor is intentionally limited to inventory and audit use cases.

- It records file locations and metadata presence only
- It does not open browser databases to extract credentials or content
- It does not exfiltrate data through HTTP, email, Telegram, or any other channel
- It writes reports to the local filesystem only

## Example Use Cases

- Browser profile discovery during endpoint triage
- Recovery preparation before system migration or cleanup
- Security inventory of browser artifacts on developer workstations
- Local audit reporting for incident response workflows
- Identifying which browser profiles and storage files exist on a machine

## Installation

### Download Prebuilt Releases

Prebuilt binaries are published on the GitHub Releases page:

`https://github.com/r4j3sh-com/BrowserAuditor/releases`

Release artifacts are generated for:

- macOS `amd64`
- macOS `arm64`
- Linux `amd64`
- Linux `arm64`
- Windows `amd64`
- Windows `arm64`

### macOS

Download the matching `darwin` archive from Releases, then extract and run:

```bash
tar -xzf BrowserAuditor_VERSION_darwin_arm64.tar.gz
chmod +x BrowserAuditor
./BrowserAuditor
```

For Intel Macs, use the `darwin_amd64` archive instead of `darwin_arm64`.

### Linux

Download the matching `linux` archive from Releases, then extract and run:

```bash
tar -xzf BrowserAuditor_VERSION_linux_amd64.tar.gz
chmod +x BrowserAuditor
./BrowserAuditor
```

For ARM Linux systems, use the `linux_arm64` archive.

### Windows

Download the matching `windows` archive from Releases, extract it, and run:

```powershell
Expand-Archive .\BrowserAuditor_VERSION_windows_amd64.zip
.\BrowserAuditor.exe
```

For Windows on ARM, use the `windows_arm64` archive.

### Build From Source

Clone the repository:

```bash
git clone https://github.com/r4j3sh-com/BrowserAuditor.git
cd BrowserAuditor
```

Build the binary:

```bash
go build -o BrowserAuditor .
```

On Windows, build with:

```powershell
go build -o BrowserAuditor.exe .
```

## Usage

Run with an auto-generated timestamped output filename:

```bash
./BrowserAuditor
```

Example default filename:

```text
browser-audit-report-20260302-032507Z.json
```

Write to a custom JSON file:

```bash
./BrowserAuditor -out snapshot.json
```

Run directly with Go:

```bash
go run . -out browser-auditor.json
```

## Multi-OS Release Process

BrowserAuditor includes a GitHub Actions release workflow that builds cross-platform binaries and attaches them to a GitHub Release.

### Supported Release Targets

- `darwin/amd64`
- `darwin/arm64`
- `linux/amd64`
- `linux/arm64`
- `windows/amd64`
- `windows/arm64`

### How To Publish A Release

Create and push a semantic version tag:

```bash
git tag v1.0.0
git push origin v1.0.0
```

That workflow will:

- build binaries for every supported OS and architecture
- package Unix builds as `.tar.gz`
- package Windows builds as `.zip`
- generate `SHA256SUMS`
- create a GitHub Release and upload all archives

You can also publish a GitHub Release from the GitHub UI. The workflow now supports both:

- tag push events
- GitHub Release `published` events

You can also trigger the same workflow manually from GitHub Actions using `workflow_dispatch`, and you must provide a release tag such as `v1.0.0`.

## Output

BrowserAuditor writes a JSON report with:

- Tool name and version
- Generation timestamp
- Machine identifier
- System metadata
- Browser inventory
- Profile directories
- Existing artifact file paths
- Notes describing collection limits

Example output:

```json
{
  "tool_name": "BrowserAuditor",
  "tool_version": "1.0",
  "generated_at": "2026-03-02T02:57:46Z",
  "machine_id": "c0176db1e92a9adaabc9db02c9bbb217",
  "system_info": {
    "os": "darwin",
    "arch": "arm64",
    "hostname": "host.local",
    "username": "user",
    "home_dir": "/Users/user",
    "local_ips": ["192.168.1.10"]
  },
  "browser_audit": [],
  "notes": [
    "Local-only recovery snapshot."
  ]
}
```

Actual results depend on the browser profiles and artifact files present on the system being audited.

## Project Structure

- `browserauditor.go` - CLI entrypoint, output naming, snapshot generation, JSON writing
- `getdata.go` - system metadata collection and browser profile discovery
- `browser-auditor.json` - example generated output
- `.github/workflows/release.yml` - multi-OS GitHub Release pipeline

## Privacy

BrowserAuditor is built for privacy-conscious local analysis.

- No external API calls
- No remote reporting
- No credential extraction
- No secret collection

## Roadmap Ideas

- File metadata and hashing for detected artifacts
- Snapshot diffing between two reports
- Extension inventory for supported browsers
- CSV export and concise terminal summaries
- Automated tests for path detection and profile discovery

## License

This project is licensed under the MIT License. See [`LICENSE`](LICENSE).

## SEO Keywords

BrowserAuditor, browser audit tool, browser profile inventory, browser artifact scanner, forensic browser triage, local browser audit, browser recovery tool, Chrome profile audit, Firefox profile audit, incident response browser inventory, browser forensics Go tool
