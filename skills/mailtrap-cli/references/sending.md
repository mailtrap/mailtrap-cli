# sending

Detailed flag specifications for `mailtrap send` commands.

---

## send transactional

Send a single transactional email via the Mailtrap Sending API.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--from` | string | Yes | Sender address (`"email"` or `"Name <email>"`) — must be on a verified domain |
| `--to` | string[] | Yes | Recipient(s), can be repeated |
| `--subject` | string | Yes | Email subject line |
| `--text` | string | No | Plain-text body |
| `--html` | string | No | HTML body |
| `--cc` | string[] | No | CC recipients, can be repeated |
| `--bcc` | string[] | No | BCC recipients, can be repeated |
| `--category` | string | No | Email category for filtering |
| `--template-uuid` | string | No | Template UUID — uses a pre-built template |
| `--reply-to` | string | No | Reply-to address |

**API endpoint:** `POST https://send.api.mailtrap.io/api/send`

**Output:** JSON response with message ID on success.

**Example:**
```bash
mailtrap send transactional \
  --from "App <noreply@yourdomain.com>" \
  --to user@example.com \
  --subject "Welcome!" \
  --html "<h1>Hello</h1>"
```

---

## send bulk

Send a single bulk/marketing email via the Mailtrap Bulk Sending API.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--from` | string | Yes | Sender address |
| `--to` | string[] | Yes | Recipient(s), can be repeated |
| `--subject` | string | Yes | Email subject line |
| `--text` | string | No | Plain-text body |
| `--html` | string | No | HTML body |
| `--cc` | string[] | No | CC recipients |
| `--bcc` | string[] | No | BCC recipients |
| `--category` | string | No | Email category |
| `--template-uuid` | string | No | Template UUID |
| `--reply-to` | string | No | Reply-to address |

**API endpoint:** `POST https://bulk.api.mailtrap.io/api/send`

---

## send batch-transactional

Send a batch of transactional emails from a JSON file.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--file` | string | Yes | Path to JSON file with email array |

**JSON file format:**
```json
[
  {
    "from": {"email": "noreply@yourdomain.com", "name": "App"},
    "to": [{"email": "user1@example.com"}],
    "subject": "Hello 1",
    "text": "Body 1"
  },
  {
    "from": {"email": "noreply@yourdomain.com"},
    "to": [{"email": "user2@example.com"}],
    "subject": "Hello 2",
    "html": "<b>Body 2</b>"
  }
]
```

---

## send batch-bulk

Send a batch of bulk emails from a JSON file.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--file` | string | Yes | Path to JSON file with email array |

Same JSON format as `send batch-transactional`, sent via the bulk API endpoint.
