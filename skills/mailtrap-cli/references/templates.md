# templates

Detailed flag specifications for `mailtrap templates` commands.

---

## templates list

List all email templates for the account.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--per-page` | int | No | Number of templates per page (default 50, max 100) |
| `--token` | int | No | Page number to retrieve (page-token pagination) |

**Output:** Table of templates with ID, UUID, name, subject, category and creation time, followed by `Next page: --token N` when more pages exist. With `--output json`, the full response object is printed: `.data` holds the templates and `.pagination.next_token` the next page (`null` on the last page).

---

## templates get

Get a specific email template.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | int | Yes | Template ID |

---

## templates create

Create a new email template.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | Yes | Template name |
| `--subject` | string | Yes | Template subject line |
| `--body-html` | string | No | HTML body content |
| `--body-text` | string | No | Plain text body content |
| `--category` | string | No | Template category (default `General`) |

**Example:**
```bash
mailtrap templates create \
  --name "Welcome Email" \
  --subject "Welcome to {{company}}" \
  --body-html "<h1>Welcome, {{name}}!</h1>" \
  --body-text "Welcome, {{name}}!"
```

---

## templates update

Update an existing email template.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | int | Yes | Template ID |
| `--name` | string | No | New template name |
| `--subject` | string | No | New subject line |
| `--body-html` | string | No | New HTML body |
| `--body-text` | string | No | New text body |
| `--category` | string | No | New category |

Only supplied flags are updated; omitted fields remain unchanged.

---

## templates delete

Delete an email template.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | int | Yes | Template ID |
