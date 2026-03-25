# Mailtrap CLI

Command-line interface for the [Mailtrap](https://mailtrap.io) email delivery platform. Send transactional and bulk emails, manage sandbox inboxes, contacts, templates, and more.

## Installation

### Homebrew (macOS / Linux)

```bash
brew install mailtrap/tap/mailtrap
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

## Configuration

Set your API token and account ID:

```bash
mailtrap configure --api-token YOUR_TOKEN
```

Or use environment variables:

```bash
export MAILTRAP_API_TOKEN=your-token
export MAILTRAP_ACCOUNT_ID=your-account-id
```

Or create `~/.mailtrap.yaml`:

```yaml
api-token: your-token
account-id: "your-account-id"
```

## Usage

### Send email

```bash
# Transactional
mailtrap send transactional \
  --from "you@yourdomain.com" \
  --to "recipient@example.com" \
  --subject "Hello" \
  --text "Hello from Mailtrap CLI"

# Bulk
mailtrap send bulk \
  --from "you@yourdomain.com" \
  --to "recipient@example.com" \
  --subject "Newsletter" \
  --html "<h1>Hello</h1>"
```

### Sandbox testing

```bash
# Send test email to sandbox
mailtrap sandbox-send single \
  --inbox-id 12345 \
  --from "test@example.com" \
  --to "recipient@example.com" \
  --subject "Test" \
  --text "Hello"

# List messages in sandbox inbox
mailtrap messages list --inbox-id 12345

# Get spam score
mailtrap messages spam-score --inbox-id 12345 --id 67890
```

### Manage resources

```bash
# Domains
mailtrap domains list
mailtrap domains get --id 123

# Templates
mailtrap templates list
mailtrap templates create --name "Welcome" --subject "Hello" --text "Welcome!"

# Contacts
mailtrap contacts create --email "user@example.com"
mailtrap contact-lists list

# Projects & Inboxes
mailtrap projects list
mailtrap inboxes list
```

### Output formats

```bash
# Table (default)
mailtrap domains list

# JSON
mailtrap domains list --output json

# Text
mailtrap domains list --output text
```

## Commands

| Group | Commands |
|-------|----------|
| **Sending** | `send transactional`, `send bulk`, `send batch-transactional`, `send batch-bulk` |
| **Sandbox** | `sandbox-send single`, `sandbox-send batch` |
| **Domains** | `domains list`, `domains get`, `domains create`, `domains delete` |
| **Templates** | `templates list`, `templates get`, `templates create`, `templates update`, `templates delete` |
| **Suppressions** | `suppressions list`, `suppressions delete` |
| **Stats** | `stats get`, `stats by-domain`, `stats by-category`, `stats by-esp`, `stats by-date` |
| **Email Logs** | `email-logs list`, `email-logs get` |
| **Projects** | `projects list`, `projects get`, `projects create`, `projects update`, `projects delete` |
| **Inboxes** | `inboxes list`, `inboxes get`, `inboxes create`, `inboxes update`, `inboxes delete`, `inboxes clean`, `inboxes mark-read`, `inboxes reset-credentials`, `inboxes toggle-email-username`, `inboxes reset-email-username` |
| **Messages** | `messages list`, `messages get`, `messages update`, `messages delete`, `messages forward`, `messages spam-score`, `messages html-analysis`, `messages headers`, `messages html`, `messages text`, `messages source`, `messages raw`, `messages eml` |
| **Attachments** | `attachments list`, `attachments get` |
| **Contacts** | `contacts get`, `contacts create`, `contacts update`, `contacts delete`, `contacts import`, `contacts import-status`, `contacts export`, `contacts export-status`, `contacts create-event` |
| **Contact Lists** | `contact-lists list`, `contact-lists get`, `contact-lists create`, `contact-lists update`, `contact-lists delete` |
| **Contact Fields** | `contact-fields list`, `contact-fields get`, `contact-fields create`, `contact-fields update`, `contact-fields delete` |
| **Accounts** | `accounts list` |
| **Account Access** | `account-access list`, `account-access remove` |
| **Permissions** | `permissions resources`, `permissions bulk-update` |
| **Tokens** | `tokens list`, `tokens get`, `tokens create`, `tokens delete`, `tokens reset` |
| **Billing** | `billing usage` |
| **Organizations** | `organizations list-sub-accounts`, `organizations create-sub-account` |
| **Config** | `configure`, `completion [bash\|zsh\|fish\|powershell]` |

## Releasing

Releases are automated via [GoReleaser](https://goreleaser.com). To create a release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

This triggers the GitHub Actions workflow which builds binaries for Linux, macOS, and Windows (amd64/arm64), creates a GitHub Release, and updates the Homebrew tap.

### Prerequisites for Homebrew tap

1. Create a repo `mailtrap/homebrew-tap`
2. Add a `HOMEBREW_TAP_GITHUB_TOKEN` secret to this repo with write access to the tap repo
