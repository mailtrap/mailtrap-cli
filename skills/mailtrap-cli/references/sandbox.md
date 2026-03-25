# sandbox

Detailed flag specifications for sandbox commands: `mailtrap sandbox-send`, `mailtrap projects`,
`mailtrap sandboxes`, `mailtrap messages`, and `mailtrap attachments`.

---

## sandbox-send single

Send a single email to a sandbox for testing.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |
| `--from` | string | Yes | Sender address |
| `--to` | string[] | Yes | Recipient(s), can be repeated |
| `--subject` | string | Yes | Email subject |
| `--text` | string | No | Plain-text body |
| `--html` | string | No | HTML body |
| `--cc` | string[] | No | CC recipients |
| `--bcc` | string[] | No | BCC recipients |
| `--category` | string | No | Email category |
| `--template-uuid` | string | No | Template UUID |
| `--reply-to` | string | No | Reply-to address |

**Note:** Sandbox emails are not delivered to real inboxes. They appear in the Mailtrap sandbox for inspection.

---

## sandbox-send batch

Send a batch of emails to a sandbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |
| `--file` | string | Yes | Path to JSON file with email array |

---

## projects list

List all sandbox projects.

No additional flags.

---

## projects get

Get a specific sandbox project.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Project ID |

---

## projects create

Create a new sandbox project.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | Yes | Project name |

---

## projects update

Update a sandbox project.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Project ID |
| `--name` | string | Yes | New project name |

---

## projects delete

Delete a sandbox project.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Project ID |

---

## sandboxes list

List all sandboxes.

No additional flags.

---

## sandboxes get

Get a specific sandbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Sandbox ID |

---

## sandboxes create

Create a new sandbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--project-id` | string | Yes | Project ID to create the sandbox in |
| `--name` | string | Yes | Sandbox name |

---

## sandboxes update

Update a sandbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Sandbox ID |
| `--name` | string | Yes | New sandbox name |

---

## sandboxes delete

Delete a sandbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Sandbox ID |

---

## sandboxes clean

Delete all messages in a sandbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Sandbox ID |

---

## sandboxes mark-read

Mark all messages in a sandbox as read.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Sandbox ID |

---

## sandboxes reset-credentials

Reset SMTP credentials of a sandbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Sandbox ID |

---

## sandboxes toggle-email

Toggle email forwarding for a sandbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Sandbox ID |
| `--email` | string | Yes | Email address to toggle forwarding for |

---

## sandboxes reset-email

Reset the email username of a sandbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Sandbox ID |

---

## messages list

List all messages in a sandbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |

---

## messages get

Get a specific sandbox message.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |
| `--id` | string | Yes | Message ID |

---

## messages update

Update a sandbox message (mark as read/unread).

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |
| `--id` | string | Yes | Message ID |
| `--is-read` | bool | No | Mark as read (default: false) |

---

## messages delete

Delete a sandbox message.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |
| `--id` | string | Yes | Message ID |

---

## messages forward

Forward a sandbox message to a real email address.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |
| `--id` | string | Yes | Message ID |
| `--email` | string | Yes | Destination email address |

---

## messages html

Get the rendered HTML body of a sandbox message.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |
| `--id` | string | Yes | Message ID |

---

## messages text

Get the plain text body of a sandbox message.

Same flags as `messages html`.

---

## messages raw

Get the raw email content of a sandbox message.

Same flags as `messages html`.

---

## messages source

Get the raw HTML source of a sandbox message.

Same flags as `messages html`.

---

## messages eml

Download the EML file of a sandbox message.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |
| `--id` | string | Yes | Message ID |
| `--output-file` | string | No | File path to save EML (default: stdout) |

---

## messages headers

Get the mail headers of a sandbox message.

Same flags as `messages html`.

---

## messages spam-score

Get the spam score analysis of a sandbox message.

Same flags as `messages html`.

---

## messages html-analysis

Get the HTML analysis report of a sandbox message.

Same flags as `messages html`.

---

## attachments list

List attachments for a sandbox message.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |
| `--message-id` | string | Yes | Message ID |

---

## attachments get

Get a specific attachment from a sandbox message.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sandbox-id` | string | Yes | Sandbox ID |
| `--message-id` | string | Yes | Message ID |
| `--id` | string | Yes | Attachment ID |
