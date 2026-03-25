# templates

Detailed flag specifications for `mailtrap templates` commands.

---

## templates list

List all email templates for the account.

No additional flags.

**Output:** Table/JSON of templates with ID, name, subject, and category.

---

## templates get

Get a specific email template.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Template ID |

---

## templates create

Create a new email template.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | Yes | Template name |
| `--subject` | string | Yes | Template subject line |
| `--body-html` | string | No | HTML body content |
| `--body-text` | string | No | Plain text body content |
| `--category` | string | No | Template category |

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
| `--id` | string | Yes | Template ID |
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
| `--id` | string | Yes | Template ID |
