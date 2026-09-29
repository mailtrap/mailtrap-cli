# inbound

Detailed flag specifications for `mailtrap inbound` commands (folders, inboxes, messages, threads).

`inbound` commands go to `https://mailtrap.io/api/inbound/...` and do **not** use `--account-id`.

---

## inbound folders list

List all inbound folders. No additional flags.

---

## inbound folders get

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Folder ID |

---

## inbound folders create

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | Yes | Folder name |

---

## inbound folders update

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Folder ID |
| `--name` | string | No | New folder name |

---

## inbound folders delete

Removes the folder and all of its inboxes.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Folder ID |

---

## inbound inboxes list

Inboxes are managed within a folder.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--folder-id` | string | Yes | Folder ID |

---

## inbound inboxes get

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--folder-id` | string | Yes | Folder ID |
| `--id` | string | Yes | Inbox ID |

---

## inbound inboxes create

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--folder-id` | string | Yes | Folder ID |
| `--name` | string | Yes | Inbox name |
| `--domain-id` | int | No | Sending domain ID for a custom-domain (catch-all) inbox; omit for a Mailtrap-hosted inbox |

---

## inbound inboxes update

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--folder-id` | string | Yes | Folder ID |
| `--id` | string | Yes | Inbox ID |
| `--name` | string | No | New inbox name |

---

## inbound inboxes delete

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--folder-id` | string | Yes | Folder ID |
| `--id` | string | Yes | Inbox ID |

---

## inbound messages list

Messages and threads are accessed via the top-level inbox route (`/api/inbound/inboxes/{inbox-id}/...`).

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--last-id` | string | No | Pagination cursor (`last_id` from the previous response) |

**Output:** In JSON, the full response: `{"data": [...], "total_count": N, "last_id": "..."}`. Pass `last_id` as `--last-id` for the next page; it is `null` on the last page. Table and text print the total and a `Next page: --last-id <value>` footer.

---

## inbound messages get

Returns the message with its body and attachment download URLs. JSON output is the API response as returned, including `attachments` (with `download_url`), `headers`, `references`, `bcc`, `reply_to`, `raw_message_url` and body sizes.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--id` | string | Yes | Message ID |

---

## inbound messages delete

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--id` | string | Yes | Message ID |

---

## inbound messages reply / reply-all / forward

Each sends a **real email** and returns the sent message IDs.

- `reply` — sends to the original sender.
- `reply-all` — sends to the original sender and copies the other recipients.
- `forward` — sends to new recipients; `--to` is required.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--id` | string | Yes | Message ID |
| `--from` | string | No | Sender address, `Name <email>` or `email` (custom-domain inboxes only) |
| `--to` | string | forward only | Recipient address, `Name <email>` or `email` (repeatable) |
| `--cc` | string | No | CC recipient (repeatable) |
| `--bcc` | string | No | BCC recipient (repeatable) |
| `--reply-to` | string | No | Reply-To address |
| `--text` | string | No | Plain-text body |
| `--html` | string | No | HTML body |
| `--category` | string | No | Email API category for the sent message |

**Notes:**

- Addresses accept `Name <email>` or a bare `email`, the same as the `send` command.
- `reply`/`reply-all` require a body (`--text` and/or `--html`); the subject is derived from the original message (prefixed `Re:`) — there is no `--subject` flag.
- `forward` may omit the body (the original is quoted automatically) and requires at least one `--to`.

---

## inbound threads list

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--last-id` | string | No | Pagination cursor (`last_id` from the previous response) |

**Output:** Same shape as `inbound messages list`: `{"data": [...], "total_count": N, "last_id": "..."}` in JSON.

---

## inbound threads get

Returns the thread with its messages embedded (oldest first).

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--id` | string | Yes | Thread ID |

---

## inbound threads delete

Inbound messages in the thread are removed; sent messages are preserved.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--inbox-id` | string | Yes | Inbox ID |
| `--id` | string | Yes | Thread ID |
