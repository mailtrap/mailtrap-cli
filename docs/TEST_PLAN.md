# Mailtrap CLI — Integration Test Plan

## Overview

This document defines a comprehensive integration test plan that tests every CLI command against the **real Mailtrap API**. Tests verify actual HTTP responses, correct output formatting, and end-to-end workflows.

## Prerequisites

### 1. Configuration

Tests read credentials from `~/.mailtrap.yaml`:

```yaml
api-token: <your-mailtrap-api-token>
```

The `.mailtrap.yaml` file is in `.gitignore` — **never commit credentials**.

You can also set an environment variable instead:
```bash
export MAILTRAP_API_TOKEN=your-token
```

### 2. Sandbox Setup

Before running integration tests, ensure you have:
- A Mailtrap account with API access
- At least one **sandbox project** with an **inbox**
- Note down your **inbox ID** for sandbox send tests
- A verified **sending domain** for transactional/bulk send tests

### 3. Build the CLI

```bash
make build
# or
go build -o mailtrap .
```

---

## Known API Response Shapes

These were discovered during integration testing and are important for correct implementation:

| Endpoint | Response Shape | Notes |
|----------|---------------|-------|
| `GET /sending_domains` | `{"data": [...]}` | Wrapped in `data` key — **not** a flat array |
| `GET /email_logs` | `{"messages": [...], "total_count": N, "next_page_cursor": "..."}` | Paginated with cursor |
| `GET /billing/usage` | `{"billing": {...}, "testing": {...}, "sending": {...}}` | Nested object, not array |
| `GET /stats` | Requires `start_date` param | Returns error without it |
| `GET /api_tokens` | `{"errors": "Access forbidden"}` | May require admin-level token |
| `GET /contacts` | 404 (HTML page) | Endpoint may not exist or requires different path |
| `GET /suppressions` | `[...]` | Flat array |
| `GET /email_templates` | `[...]` | Flat array |
| `GET /contacts/lists` | `[...]` | Flat array |
| `GET /contacts/fields` | `[...]` | Flat array |
| `GET /account_accesses` | `[...]` | Flat array |
| `GET /permissions/resources` | `[...]` | Flat array |
| `GET /projects` | `[...]` | Flat array |
| `GET /inboxes` | `[...]` | Flat array |

---

## Test Execution

### Running All Integration Tests

```bash
# Build and run all commands, checking exit codes and output
./docs/run_integration_tests.sh

# Or test individual commands manually (see sections below)
```

### Test Categories

Tests are organized by endpoint group. Each test specifies:
- **Command** — the exact CLI invocation
- **Expected** — what to check in the output
- **Cleanup** — any resources to delete after

---

## 1. Accounts

| # | Test | Command | Expected |
|---|------|---------|----------|
| 1.1 | List accounts | `mailtrap accounts list` | Table with ID, NAME columns; contains account name |
| 1.2 | List accounts (JSON) | `mailtrap accounts list --output json` | Valid JSON array with `id`, `name`, `access_levels` fields |

## 2. Projects (Sandbox)

| # | Test | Command | Expected |
|---|------|---------|----------|
| 2.1 | List projects | `mailtrap projects list` | Table with existing projects |
| 2.2 | Get project | `mailtrap projects get --id <PROJECT_ID>` | Single project details |
| 2.3 | Create project | `mailtrap projects create --name "integration-test"` | New project in output |
| 2.4 | Update project | `mailtrap projects update --id <NEW_ID> --name "integration-test-updated"` | Updated name in output |
| 2.5 | Delete project | `mailtrap projects delete --id <NEW_ID>` | Success message |
| 2.6 | Get missing ID | `mailtrap projects get` | Error: `--id is required` |

**Cleanup:** Delete any created test project.

## 3. Inboxes (Sandbox)

| # | Test | Command | Expected |
|---|------|---------|----------|
| 3.1 | List inboxes | `mailtrap inboxes list` | Table with inbox entries |
| 3.2 | Get inbox | `mailtrap inboxes get --id <INBOX_ID>` | Single inbox details |
| 3.3 | Create inbox | `mailtrap inboxes create --project-id <PROJECT_ID> --name "test-inbox"` | New inbox in output |
| 3.4 | Update inbox | `mailtrap inboxes update --id <NEW_ID> --name "test-inbox-updated"` | Updated name |
| 3.5 | Mark read | `mailtrap inboxes mark-read --id <INBOX_ID>` | Success message |
| 3.6 | Reset credentials | `mailtrap inboxes reset-credentials --id <NEW_ID>` | Success/new credentials |
| 3.7 | Toggle email username | `mailtrap inboxes toggle-email-username --id <NEW_ID>` | Toggled response |
| 3.8 | Reset email username | `mailtrap inboxes reset-email-username --id <NEW_ID>` | Reset response |
| 3.9 | Clean inbox | `mailtrap inboxes clean --id <NEW_ID>` | Success message |
| 3.10 | Delete inbox | `mailtrap inboxes delete --id <NEW_ID>` | Success message |
| 3.11 | Get missing ID | `mailtrap inboxes get` | Error: `--id is required` |

**Cleanup:** Delete test inbox.

## 4. Sandbox Send

| # | Test | Command | Expected |
|---|------|---------|----------|
| 4.1 | Send single | `mailtrap sandbox-send single --inbox-id <INBOX_ID> --from test@example.com --to recipient@example.com --subject "Integration Test" --text "Hello from CLI"` | Success response with message ID |
| 4.2 | Send batch | `mailtrap sandbox-send batch --inbox-id <INBOX_ID> --from test@example.com --to recipient@example.com --subject "Batch Test" --text "Batch hello"` | Success response |
| 4.3 | Missing flags | `mailtrap sandbox-send single` | Error: required flags |

**Verification:** After 4.1, run `mailtrap messages list --inbox-id <INBOX_ID>` to confirm the message appears.

## 5. Messages (Sandbox)

Prerequisite: Send a test email to the sandbox inbox first (test 4.1).

| # | Test | Command | Expected |
|---|------|---------|----------|
| 5.1 | List messages | `mailtrap messages list --inbox-id <INBOX_ID>` | Table with message entries |
| 5.2 | List messages (JSON) | `mailtrap messages list --inbox-id <INBOX_ID> --output json` | Valid JSON array |
| 5.3 | Get message | `mailtrap messages get --inbox-id <INBOX_ID> --id <MSG_ID>` | Single message details |
| 5.4 | Update message | `mailtrap messages update --inbox-id <INBOX_ID> --id <MSG_ID> --is-read true` | Updated message |
| 5.5 | Spam score | `mailtrap messages spam-score --inbox-id <INBOX_ID> --id <MSG_ID>` | Spam score data |
| 5.6 | HTML analysis | `mailtrap messages html-analysis --inbox-id <INBOX_ID> --id <MSG_ID>` | Analysis data |
| 5.7 | Headers | `mailtrap messages headers --inbox-id <INBOX_ID> --id <MSG_ID>` | Email headers |
| 5.8 | HTML body | `mailtrap messages html --inbox-id <INBOX_ID> --id <MSG_ID>` | Raw HTML content |
| 5.9 | Text body | `mailtrap messages text --inbox-id <INBOX_ID> --id <MSG_ID>` | Raw text content |
| 5.10 | Source | `mailtrap messages source --inbox-id <INBOX_ID> --id <MSG_ID>` | Raw source |
| 5.11 | Raw | `mailtrap messages raw --inbox-id <INBOX_ID> --id <MSG_ID>` | Raw message |
| 5.12 | EML | `mailtrap messages eml --inbox-id <INBOX_ID> --id <MSG_ID>` | EML format |
| 5.13 | Forward | `mailtrap messages forward --inbox-id <INBOX_ID> --id <MSG_ID> --email forward@example.com` | Success |
| 5.14 | Delete message | `mailtrap messages delete --inbox-id <INBOX_ID> --id <MSG_ID>` | Success message |
| 5.15 | Missing inbox-id | `mailtrap messages list` | Error: `--inbox-id is required` |
| 5.16 | Missing id | `mailtrap messages get --inbox-id <INBOX_ID>` | Error: `--id is required` |

## 6. Attachments (Sandbox)

Prerequisite: Send an email with an attachment to the sandbox inbox.

| # | Test | Command | Expected |
|---|------|---------|----------|
| 6.1 | List attachments | `mailtrap attachments list --inbox-id <INBOX_ID> --message-id <MSG_ID>` | Table with attachments |
| 6.2 | Get attachment | `mailtrap attachments get --inbox-id <INBOX_ID> --message-id <MSG_ID> --id <ATT_ID>` | Attachment details |
| 6.3 | Missing flags | `mailtrap attachments list` | Error: required flags |

## 7. Email Sending (Transactional/Bulk)

**Note:** Requires a verified sending domain.

| # | Test | Command | Expected |
|---|------|---------|----------|
| 7.1 | Send transactional | `mailtrap send transactional --from you@yourdomain.com --to recipient@example.com --subject "CLI Test" --text "Hello"` | Success with message IDs |
| 7.2 | Send bulk | `mailtrap send bulk --from you@yourdomain.com --to recipient@example.com --subject "Bulk CLI Test" --text "Hello"` | Success |
| 7.3 | Send with named from | `mailtrap send transactional --from "CLI Test <you@yourdomain.com>" --to recipient@example.com --subject "Named" --text "Hi"` | Success, `from.name` populated |
| 7.4 | Batch transactional | `mailtrap send batch-transactional --from you@yourdomain.com --to recipient@example.com --subject "Batch" --text "Hi"` | Success |
| 7.5 | Batch bulk | `mailtrap send batch-bulk --from you@yourdomain.com --to recipient@example.com --subject "Batch Bulk" --text "Hi"` | Success |
| 7.6 | Missing flags | `mailtrap send transactional` | Error: required flags |
| 7.7 | JSON output | `mailtrap send transactional --from you@yourdomain.com --to recipient@example.com --subject "JSON" --text "Hi" --output json` | Valid JSON |

## 8. Domains (Sending)

**Note:** The API returns `{"data": [...]}` — the CLI may need a fix to unwrap this.

| # | Test | Command | Expected |
|---|------|---------|----------|
| 8.1 | List domains | `mailtrap domains list` | Table with domain entries |
| 8.2 | Get domain | `mailtrap domains get --id <DOMAIN_ID>` | Single domain details |
| 8.3 | Create domain | `mailtrap domains create --name "test-integration.example.com"` | New domain in output |
| 8.4 | Delete domain | `mailtrap domains delete --id <NEW_ID>` | Success message |
| 8.5 | Get missing ID | `mailtrap domains get` | Error: `--id is required` |

## 9. Templates

| # | Test | Command | Expected |
|---|------|---------|----------|
| 9.1 | List templates | `mailtrap templates list` | Table with template entries |
| 9.2 | List templates (JSON) | `mailtrap templates list --output json` | Valid JSON array |
| 9.3 | Get template | `mailtrap templates get --id <TEMPLATE_ID>` | Single template details |
| 9.4 | Create template | `mailtrap templates create --name "test-tpl" --subject "Test" --text "body"` | New template in output |
| 9.5 | Update template | `mailtrap templates update --id <NEW_ID> --name "test-tpl-updated"` | Updated template |
| 9.6 | Delete template | `mailtrap templates delete --id <NEW_ID>` | Success message |
| 9.7 | Get missing ID | `mailtrap templates get` | Error: `--id is required` |

**Cleanup:** Delete created template.

## 10. Suppressions

| # | Test | Command | Expected |
|---|------|---------|----------|
| 10.1 | List suppressions | `mailtrap suppressions list` | Table with suppression entries (may be empty) |
| 10.2 | Delete suppression | `mailtrap suppressions delete --id <SUPP_ID>` | Success message |
| 10.3 | Delete missing ID | `mailtrap suppressions delete` | Error: `--id is required` |

## 11. Stats

**Note:** Stats endpoint requires `--start-date` parameter.

| # | Test | Command | Expected |
|---|------|---------|----------|
| 11.1 | Get stats | `mailtrap stats get --start-date 2025-01-01` | Stats data |
| 11.2 | Stats by domain | `mailtrap stats by-domain --start-date 2025-01-01` | Domain-grouped stats |
| 11.3 | Stats by category | `mailtrap stats by-category --start-date 2025-01-01` | Category-grouped stats |
| 11.4 | Stats by ESP | `mailtrap stats by-esp --start-date 2025-01-01` | ESP-grouped stats |
| 11.5 | Stats by date | `mailtrap stats by-date --start-date 2025-01-01` | Date-grouped stats |
| 11.6 | Stats JSON | `mailtrap stats get --start-date 2025-01-01 --output json` | Valid JSON |

## 12. Email Logs

**Note:** API returns `{"messages": [...], "total_count": N, "next_page_cursor": "..."}` — CLI correctly unwraps the `messages` key.

| # | Test | Command | Expected |
|---|------|---------|----------|
| 12.1 | List email logs | `mailtrap email-logs list` | Table with log entries |
| 12.2 | Get email log | `mailtrap email-logs get --id <LOG_ID>` | Single log details |
| 12.3 | List JSON | `mailtrap email-logs list --output json` | Valid JSON |
| 12.4 | Get missing ID | `mailtrap email-logs get` | Error: `--id is required` |

## 13. Contact Lists

| # | Test | Command | Expected |
|---|------|---------|----------|
| 13.1 | List contact lists | `mailtrap contact-lists list` | Table with list entries |
| 13.2 | Get contact list | `mailtrap contact-lists get --id <LIST_ID>` | Single list details |
| 13.3 | Create contact list | `mailtrap contact-lists create --name "test-list"` | New list in output |
| 13.4 | Update contact list | `mailtrap contact-lists update --id <NEW_ID> --name "test-list-updated"` | Updated list |
| 13.5 | Delete contact list | `mailtrap contact-lists delete --id <NEW_ID>` | Success message |
| 13.6 | Get missing ID | `mailtrap contact-lists get` | Error: `--id is required` |

**Cleanup:** Delete created list.

## 14. Contact Fields

| # | Test | Command | Expected |
|---|------|---------|----------|
| 14.1 | List contact fields | `mailtrap contact-fields list` | Table with field entries |
| 14.2 | Get contact field | `mailtrap contact-fields get --id <FIELD_ID>` | Single field details |
| 14.3 | Create contact field | `mailtrap contact-fields create --name "test-field" --type "string"` | New field in output |
| 14.4 | Update contact field | `mailtrap contact-fields update --id <NEW_ID> --name "test-field-updated"` | Updated field |
| 14.5 | Delete contact field | `mailtrap contact-fields delete --id <NEW_ID>` | Success message |
| 14.6 | Get missing ID | `mailtrap contact-fields get` | Error: `--id is required` |

**Cleanup:** Delete created field.

## 15. Contacts

**Note:** The `GET /contacts` endpoint returned 404 during testing. May require a different API path or feature flag.

| # | Test | Command | Expected |
|---|------|---------|----------|
| 15.1 | Create contact | `mailtrap contacts create --email "test@integration.com"` | New contact in output |
| 15.2 | Get contact | `mailtrap contacts get --id <CONTACT_ID>` | Single contact details |
| 15.3 | Update contact | `mailtrap contacts update --id <CONTACT_ID> --first-name "Test"` | Updated contact |
| 15.4 | List contacts | `mailtrap contacts list` | Table with contacts (may 404 — investigate) |
| 15.5 | Import contacts | `mailtrap contacts import --file contacts.csv` | Import initiated |
| 15.6 | Import status | `mailtrap contacts import-status --id <IMPORT_ID>` | Status details |
| 15.7 | Export contacts | `mailtrap contacts export` | Export initiated |
| 15.8 | Export status | `mailtrap contacts export-status --id <EXPORT_ID>` | Status details |
| 15.9 | Create event | `mailtrap contacts create-event --id <CONTACT_ID> --type "purchase"` | Event created |
| 15.10 | Delete contact | `mailtrap contacts delete --id <CONTACT_ID>` | Success message |

**Cleanup:** Delete created contact.

## 16. Account Access

| # | Test | Command | Expected |
|---|------|---------|----------|
| 16.1 | List account accesses | `mailtrap account-access list` | Table with access entries (may be empty) |
| 16.2 | Remove access | `mailtrap account-access remove --id <ACCESS_ID>` | Success message |
| 16.3 | Remove missing ID | `mailtrap account-access remove` | Error: `--id is required` |

## 17. Permissions

| # | Test | Command | Expected |
|---|------|---------|----------|
| 17.1 | List resources | `mailtrap permissions resources` | Table with resource entries |
| 17.2 | Bulk update | `mailtrap permissions bulk-update --access-id <ID> --permissions '...'` | Updated permissions |
| 17.3 | Missing flags | `mailtrap permissions bulk-update` | Error: required flags |

## 18. Tokens

**Note:** Token management may require admin-level API token (got "Access forbidden" during testing).

| # | Test | Command | Expected |
|---|------|---------|----------|
| 18.1 | List tokens | `mailtrap tokens list` | Table with tokens (may require admin token) |
| 18.2 | Get token | `mailtrap tokens get --id <TOKEN_ID>` | Token details |
| 18.3 | Create token | `mailtrap tokens create --name "test-token"` | New token |
| 18.4 | Reset token | `mailtrap tokens reset --id <NEW_ID>` | New token value |
| 18.5 | Delete token | `mailtrap tokens delete --id <NEW_ID>` | Success message |
| 18.6 | Get missing ID | `mailtrap tokens get` | Error: `--id is required` |

**Cleanup:** Delete created token.

## 19. Billing

| # | Test | Command | Expected |
|---|------|---------|----------|
| 19.1 | Get usage | `mailtrap billing usage` | Usage data (billing, testing, sending sections) |
| 19.2 | Get usage (JSON) | `mailtrap billing usage --output json` | Valid JSON with nested objects |

## 20. Organizations

| # | Test | Command | Expected |
|---|------|---------|----------|
| 20.1 | List sub-accounts | `mailtrap organizations list-sub-accounts` | Table with sub-accounts |
| 20.2 | Create sub-account | `mailtrap organizations create-sub-account --name "test-sub"` | New sub-account |
| 20.3 | Missing name | `mailtrap organizations create-sub-account` | Error: `--name is required` |

**Caution:** Creating sub-accounts may have billing implications.

## 21. Configure

| # | Test | Command | Expected |
|---|------|---------|----------|
| 21.1 | Configure with token | `mailtrap configure --api-token test-token-123` | Config saved message |
| 21.2 | Configure without token | `mailtrap configure` | Error or prompts for token |

---

## Discovered Bugs / Issues

| # | Severity | Status | Description |
|---|----------|--------|-------------|
| B1 | **High** | FIXED | `domains list` — API returns `{"data": [...]}`, added `domainListResponse` wrapper struct. Also updated `Domain` fields to match API (`domain_name`, `dns_verified`, `compliance_status`). |
| B2 | **Medium** | FIXED | `email-logs list` — API returns `{"messages": [...], "total_count": N}`, added `emailLogListResponse` wrapper struct. Also updated `EmailLog` fields to match API (`message_id`, `from`, `to`). |
| B3 | **Medium** | NOT A BUG | `contacts list` returns 404 — the Mailtrap API has no `GET /contacts` endpoint. Contacts can only be managed individually (get/create/update/delete). No `list` subcommand exists in the CLI. |
| B4 | **Low** | NOT A BUG | `tokens list` returns "Access forbidden" — this is an API permission issue requiring an admin-level token, not a CLI bug. The error message is surfaced correctly. |
| B5 | **Low** | NOT A BUG | `stats get` requires `--start-date` — already enforced via `cobra.MarkFlagRequired("start-date")`. The CLI shows a clear error when omitted. |

---

## Test Execution Order (Recommended)

Run tests in dependency order so earlier tests create resources needed by later ones:

1. **Accounts** (read-only, validates API access)
2. **Projects** (CRUD — creates project for inbox tests)
3. **Inboxes** (CRUD — creates inbox for send/message tests)
4. **Sandbox Send** (sends email to inbox for message tests)
5. **Messages** (all subcommands on the sent message)
6. **Attachments** (requires message with attachment)
7. **Sending** (transactional/bulk — requires verified domain)
8. **Domains** (list/get — fix bug B1 first)
9. **Templates** (CRUD)
10. **Suppressions** (list/delete)
11. **Stats** (requires start-date)
12. **Email Logs** (fix bug B2 first)
13. **Contact Lists** (CRUD)
14. **Contact Fields** (CRUD)
15. **Contacts** (CRUD — investigate bug B3)
16. **Account Access** (read-only)
17. **Permissions** (read/update)
18. **Tokens** (CRUD — may need admin token)
19. **Billing** (read-only)
20. **Organizations** (read-only, skip create unless safe)
21. **Configure** (local config only)

---

## Summary

| Category | Endpoints | Test Cases |
|----------|-----------|------------|
| Accounts | 1 | 2 |
| Projects | 5 | 6 |
| Inboxes | 10 | 11 |
| Sandbox Send | 2 | 3 |
| Messages | 14 | 16 |
| Attachments | 2 | 3 |
| Sending | 5 | 7 |
| Domains | 4 | 5 |
| Templates | 5 | 7 |
| Suppressions | 2 | 3 |
| Stats | 5 | 6 |
| Email Logs | 2 | 4 |
| Contact Lists | 5 | 6 |
| Contact Fields | 5 | 6 |
| Contacts | 10 | 10 |
| Account Access | 2 | 3 |
| Permissions | 2 | 3 |
| Tokens | 5 | 6 |
| Billing | 1 | 2 |
| Organizations | 2 | 3 |
| Configure | 1 | 2 |
| **Total** | **~90** | **~107** |
