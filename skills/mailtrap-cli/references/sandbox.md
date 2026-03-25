# sandbox

Detailed flag specifications for sandbox commands: `mailtrap sandbox-send`, `mailtrap projects`,
`mailtrap inboxes`, `mailtrap messages`, and `mailtrap attachments`.

---

## sandbox-send single

Send a single email to a sandbox inbox for testing.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Sandbox inbox ID |
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

Send a batch of emails to a sandbox inbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Sandbox inbox ID |
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

## inboxes list

List all sandbox inboxes.

No additional flags.

---

## inboxes get

Get a specific sandbox inbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Inbox ID |

---

## inboxes create

Create a new sandbox inbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--project-id` | string | Yes | Project ID to create the inbox in |
| `--name` | string | Yes | Inbox name |

---

## inboxes update

Update a sandbox inbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Inbox ID |
| `--name` | string | Yes | New inbox name |

---

## inboxes delete

Delete a sandbox inbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Inbox ID |

---

## inboxes clean

Delete all messages in a sandbox inbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Inbox ID |

---

## inboxes mark-read

Mark all messages in a sandbox inbox as read.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Inbox ID |

---

## inboxes reset-credentials

Reset SMTP credentials of a sandbox inbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Inbox ID |

---

## inboxes toggle-email

Toggle email forwarding for a sandbox inbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Inbox ID |
| `--email` | string | Yes | Email address to toggle forwarding for |

---

## inboxes reset-email

Reset the email username of a sandbox inbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Inbox ID |

---

## messages list

List all messages in a sandbox inbox.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |

---

## messages get

Get a specific sandbox message.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--id` | string | Yes | Message ID |

---

## messages update

Update a sandbox message (mark as read/unread).

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--id` | string | Yes | Message ID |
| `--is-read` | bool | No | Mark as read (default: false) |

---

## messages delete

Delete a sandbox message.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--id` | string | Yes | Message ID |

---

## messages forward

Forward a sandbox message to a real email address.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--id` | string | Yes | Message ID |
| `--email` | string | Yes | Destination email address |

---

## messages html

Get the rendered HTML body of a sandbox message.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
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
| `--inbox-id` | string | Yes | Inbox ID |
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
| `--inbox-id` | string | Yes | Inbox ID |
| `--message-id` | string | Yes | Message ID |

---

## attachments get

Get a specific attachment from a sandbox message.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--message-id` | string | Yes | Message ID |
| `--id` | string | Yes | Attachment ID |
