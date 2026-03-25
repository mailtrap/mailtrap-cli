# Mailtrap CLI Skill

The `mailtrap` CLI is the official command-line interface for the Mailtrap email platform.
It covers transactional/bulk email sending, sandbox testing, domain management, templates,
contacts, statistics, and account administration.

## Authentication

The CLI requires an API token. Set it via any of these (highest priority first):

1. `--api-token` flag on any command
2. `MAILTRAP_API_TOKEN` environment variable
3. `~/.mailtrap.yaml` config file (`api-token` key)

Most commands also require an account ID, set the same way (`--account-id`, `MAILTRAP_ACCOUNT_ID`, or `account-id` in config).

Run `mailtrap configure --api-token <token>` to auto-detect the account ID and save both to `~/.mailtrap.yaml`.

## Non-Interactive Mode

The CLI never prompts for input. All required values must be supplied as flags.
If a required flag is missing, the command exits with an error listing the missing flag.

## Output Formats

Use `-o` / `--output` to control output format:
- `table` (default) — human-readable table
- `json` — machine-readable JSON
- `text` — plain text

For scripting and piping, always use `--output json`.

## Key Conventions

- **Email address format**: `"email@example.com"` or `"Name <email@example.com>"`
- **Repeated flags**: flags like `--to`, `--cc`, `--bcc`, `--domain-ids` can be specified multiple times
- **Batch operations**: use `--file path/to/payload.json` with a JSON array of email objects
- **Config priority**: CLI flags > environment variables > config file
- **Exit codes**: 0 on success, 1 on error (with descriptive message)
- **API base**: all requests go to `https://mailtrap.io/api/accounts/{account-id}/...`

## Command Groups

| Group | Purpose | Reference |
|-------|---------|-----------|
| `send` | Transactional & bulk email sending | [sending.md](references/sending.md) |
| `domains` | Sending domain management | [domains.md](references/domains.md) |
| `templates` | Email template CRUD | [templates.md](references/templates.md) |
| `stats` | Aggregated sending statistics | [email-logs.md](references/email-logs.md) |
| `email-logs` | Individual email log lookup | [email-logs.md](references/email-logs.md) |
| `contacts` | Contact management & import/export | [contacts.md](references/contacts.md) |
| `contact-lists` | Contact list CRUD | [contacts.md](references/contacts.md) |
| `contact-fields` | Custom contact fields | [contacts.md](references/contacts.md) |
| `sandbox-send` | Sandbox email testing | [sandbox.md](references/sandbox.md) |
| `projects` | Sandbox project management | [sandbox.md](references/sandbox.md) |
| `sandboxes` | Sandbox management | [sandbox.md](references/sandbox.md) |
| `messages` | Sandbox message inspection | [sandbox.md](references/sandbox.md) |
| `attachments` | Sandbox attachment retrieval | [sandbox.md](references/sandbox.md) |
| `accounts` | Account listing | [accounts.md](references/accounts.md) |
| `account-access` | Account access management | [accounts.md](references/accounts.md) |
| `tokens` | API token management | [accounts.md](references/accounts.md) |
| `permissions` | Permission management | [accounts.md](references/accounts.md) |
| `billing` | Usage information | [accounts.md](references/accounts.md) |
| `organizations` | Sub-account management | [accounts.md](references/accounts.md) |
| `suppressions` | Suppression list management | [domains.md](references/domains.md) |

## Gotchas

- `sandbox-send` requires `--sandbox-id` — this is separate from the production `send` commands
- Batch commands require a `--file` flag pointing to a valid JSON file; they do not accept inline JSON
- Stats commands always require `--start-date` and `--end-date`
- Message inspection commands (`messages html`, `messages raw`, etc.) require both `--sandbox-id` and `--id`
- The `configure` command validates the token against the API before saving
- Domain verification is handled through the Mailtrap web dashboard, not the CLI
