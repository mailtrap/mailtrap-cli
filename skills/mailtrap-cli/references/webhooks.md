# webhooks

Detailed flag specifications for `mailtrap webhooks` commands.

A webhook sends events to a URL. Its type decides which events and flags apply:

| Type | Events | Flags |
|------|--------|-------|
| `email_sending` | the `--event-types` you list | `--sending-stream` and `--event-types` required; `--domain-id` optional |
| `campaigns` | the `--event-types` you list | `--event-types` required; `--domain-id` optional |
| `audit_log` | account activity | no `--sending-stream`, `--domain-id` or `--event-types` |
| `inbound_receiving` | Email Inbound messages received | no `--sending-stream`, `--domain-id` or `--event-types`; `--inbound-inbox-id` optional |

Event types: `delivery`, `soft_bounce`, `bounce`, `suspension`, `unsubscribe`, `open`, `spam_complaint`, `click`, `reject`, and `subscription_preferences_update` (bulk `email_sending` and `campaigns` webhooks only). The CLI passes them to the API unchecked.

**Output:** Webhook commands print a table with ID, URL, type, active, payload format, sending stream, domain ID, inbound inbox ID and event types.

---

## webhooks list

List the account's webhooks.

No additional flags.

---

## webhooks get

Get a webhook.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Webhook ID |

---

## webhooks create

Create a webhook.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--url` | string | Yes | URL the events are sent to |
| `--type` | string | Yes | Webhook type: `email_sending`, `campaigns`, `audit_log`, `inbound_receiving` |
| `--active` | bool | No | Whether the webhook is active (default `true`) |
| `--payload-format` | string | No | Payload format: `json`, `jsonlines` |
| `--sending-stream` | string | No | Sending stream: `transactional`, `bulk` (`email_sending` only) |
| `--event-types` | string[] | No | Event types (comma-separated) |
| `--domain-id` | int | No | Domain ID to scope the webhook to (`email_sending` and `campaigns`) |
| `--inbound-inbox-id` | int | No | Email Inbound inbox ID to scope the webhook to (`inbound_receiving` only; omit for all inboxes) |

**Output:** The webhook, with a SIGNING SECRET column. The signing secret is returned only here; store it to verify the events' signatures.

**Example:**
```bash
mailtrap webhooks create --url "https://example.com/hooks" --type email_sending \
  --sending-stream transactional --event-types delivery,bounce
```

---

## webhooks update

Update a webhook. The type, sending stream and domain can't be changed.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Webhook ID |
| `--url` | string | No | New URL |
| `--active` | bool | No | Whether the webhook is active (`--active=false` pauses it) |
| `--payload-format` | string | No | Payload format: `json`, `jsonlines` |
| `--event-types` | string[] | No | Event types (comma-separated); replaces the current list |
| `--inbound-inbox-id` | int | No | Email Inbound inbox ID to scope the webhook to (`inbound_receiving` only) |

Only supplied flags are updated; omitted fields remain unchanged. Pass at least one besides `--id`: the API answers an empty update with `API error 400: Bad Request`.

```bash
mailtrap webhooks update --id 1 --active=false
```

---

## webhooks delete

Delete a webhook.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Webhook ID |
