# domains

Detailed flag specifications for `mailtrap domains`, `mailtrap suppressions` and `mailtrap tracking-opt-outs` commands.

---

## domains list

List all sending domains for the account.

No additional flags. Uses `--account-id` from global config.

**Output:** Table/JSON of domains with ID, name, status, and DNS records.

---

## domains get

Retrieve details of a specific sending domain.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Domain ID |

---

## domains create

Register a new sending domain.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | Yes | Domain name (e.g. `yourdomain.com`) |

**Note:** After creation, configure DNS records shown in the response. Domain verification is done via the Mailtrap web dashboard.

---

## domains delete

Delete a sending domain.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Domain ID |

---

## suppressions list

List all suppressions (bounced/unsubscribed addresses).

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--email` | string | No | Filter by email address |
| `--start-time` | string | No | Filter by start time |
| `--end-time` | string | No | Filter by end time |
| `--last-id` | string | No | Pagination cursor: id of the last record from the previous response |

**Note:** The endpoint returns up to 1000 suppressions per request. Page through larger result sets with `--last-id`.

---

## suppressions create

Add an address to the suppression list.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--email` | string | Yes | Email address to suppress |
| `--domain-id` | int | Yes | Sending domain the suppression applies to |
| `--sending-stream` | string | Yes | `transactional` or `bulk` |
| `--type` | string | No | Suppression reason: `hard bounce`, `spam complaint`, `unsubscription`, `manual import`. Defaults to `manual import` |

---

## suppressions delete

Remove an address from the suppression list.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Suppression ID |

---

## tracking-opt-outs list

List addresses excluded from open and click tracking.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--email` | string | No | Filter by email address |
| `--start-time` | string | No | Filter by start time |
| `--end-time` | string | No | Filter by end time |
| `--last-id` | string | No | Pagination cursor: `last_id` from the previous response |

**Note:** Uses the API token's account; `--account-id` is not needed.

---

## tracking-opt-outs create

Exclude an address from open and click tracking.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--email` | string | Yes | Email address to opt out |
| `--domain-id` | int | Yes | Sending domain the opt-out applies to |

---

## tracking-opt-outs delete

Remove an address from the tracking opt-out list, so tracking applies again.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Tracking opt-out ID |
