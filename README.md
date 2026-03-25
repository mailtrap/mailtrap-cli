# Mailtrap CLI

[![CI](https://github.com/mailtrap/mailtrap-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/mailtrap/mailtrap-cli/actions/workflows/ci.yml)
[![Release](https://github.com/mailtrap/mailtrap-cli/actions/workflows/release.yml/badge.svg)](https://github.com/mailtrap/mailtrap-cli/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Command-line interface for the [Mailtrap](https://mailtrap.io) email delivery platform. Send transactional and bulk emails, manage sandboxes, contacts, templates, and more.

## Installation

### Homebrew (macOS / Linux)

```bash
brew install mailtrap/cli/mailtrap
```

### Download binary

Download the latest release from [GitHub Releases](https://github.com/mailtrap/mailtrap-cli/releases) and add it to your `PATH`.

### Go install

```bash
go install github.com/mailtrap/mailtrap-cli@latest
```

### Build from source

```bash
git clone https://github.com/mailtrap/mailtrap-cli.git
cd mailtrap-cli
make build
```

## Quick Start

### 1. Configure

```bash
mailtrap configure --api-token YOUR_TOKEN
```

This validates your token, detects your account ID, and saves both to `~/.mailtrap.yaml`.

You can also use environment variables:

```bash
export MAILTRAP_API_TOKEN=your-token
```

### 2. Send an email

```bash
mailtrap send transactional \
  --from "you@yourdomain.com" \
  --to "recipient@example.com" \
  --subject "Hello" \
  --text "Hello from Mailtrap CLI"
```

## Usage

### Transactional & bulk sending

```bash
# Transactional (single)
mailtrap send transactional \
  --from "App <noreply@yourdomain.com>" \
  --to user@example.com \
  --subject "Welcome!" \
  --html "<h1>Welcome</h1>"

# Bulk
mailtrap send bulk \
  --from "newsletter@yourdomain.com" \
  --to subscriber@example.com \
  --subject "Weekly Update" \
  --text "Here's what's new..."

# Batch (from JSON file)
mailtrap send batch-transactional --file emails.json
```

### Sandbox testing

```bash
# Send test email to a sandbox
mailtrap sandbox-send single \
  --sandbox-id 12345 \
  --from "test@example.com" \
  --to "recipient@example.com" \
  --subject "Test" \
  --text "Hello"

# Inspect messages
mailtrap messages list --sandbox-id 12345
mailtrap messages spam-score --sandbox-id 12345 --id 67890
```

### Manage resources

```bash
# Domains
mailtrap domains list
mailtrap domains create --name "yourdomain.com"

# Templates
mailtrap templates list
mailtrap templates create --name "Welcome" --subject "Hello {{name}}" --body-html "<h1>Hi!</h1>"

# Contacts
mailtrap contacts create --email "user@example.com" --first-name "John"
mailtrap contact-lists list

# Sandboxes & projects
mailtrap projects list
mailtrap sandboxes list
```

### Output formats

```bash
# Table (default)
mailtrap domains list

# JSON (for scripting)
mailtrap domains list --output json

# Text
mailtrap domains list --output text
```

## Commands

| Group | Commands |
|-------|----------|
| **Sending** | `send transactional`, `send bulk`, `send batch-transactional`, `send batch-bulk` |
| **Sandbox Send** | `sandbox-send single`, `sandbox-send batch` |
| **Domains** | `domains list`, `domains get`, `domains create`, `domains delete` |
| **Templates** | `templates list`, `templates get`, `templates create`, `templates update`, `templates delete` |
| **Suppressions** | `suppressions list`, `suppressions delete` |
| **Stats** | `stats get`, `stats by-domain`, `stats by-category`, `stats by-esp`, `stats by-date` |
| **Email Logs** | `email-logs list`, `email-logs get` |
| **Contacts** | `contacts get`, `contacts create`, `contacts update`, `contacts delete`, `contacts import`, `contacts export`, `contacts import-status`, `contacts export-status`, `contacts create-event` |
| **Contact Lists** | `contact-lists list`, `contact-lists get`, `contact-lists create`, `contact-lists update`, `contact-lists delete` |
| **Contact Fields** | `contact-fields list`, `contact-fields get`, `contact-fields create`, `contact-fields update`, `contact-fields delete` |
| **Projects** | `projects list`, `projects get`, `projects create`, `projects update`, `projects delete` |
| **Sandboxes** | `sandboxes list`, `sandboxes get`, `sandboxes create`, `sandboxes update`, `sandboxes delete`, `sandboxes clean`, `sandboxes mark-read`, `sandboxes reset-credentials`, `sandboxes toggle-email`, `sandboxes reset-email` |
| **Messages** | `messages list`, `messages get`, `messages update`, `messages delete`, `messages forward`, `messages spam-score`, `messages html-analysis`, `messages headers`, `messages html`, `messages text`, `messages source`, `messages raw`, `messages eml` |
| **Attachments** | `attachments list`, `attachments get` |
| **Accounts** | `accounts list` |
| **Account Access** | `account-access list`, `account-access remove` |
| **Permissions** | `permissions resources`, `permissions bulk-update` |
| **Tokens** | `tokens list`, `tokens get`, `tokens create`, `tokens delete`, `tokens reset` |
| **Billing** | `billing usage` |
| **Organizations** | `organizations list-sub-accounts`, `organizations create-sub-account` |
| **Config** | `configure`, `completion [bash\|zsh\|fish\|powershell]` |

## Shell Completion

```bash
# Bash
mailtrap completion bash > /etc/bash_completion.d/mailtrap

# Zsh
mailtrap completion zsh > "${fpath[1]}/_mailtrap"

# Fish
mailtrap completion fish > ~/.config/fish/completions/mailtrap.fish
```

## Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes
4. Push and open a Pull Request

## License

[MIT](LICENSE)
