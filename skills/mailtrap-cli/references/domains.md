# domains

Detailed flag specifications for `mailtrap domains`, `mailtrap company-info`, `mailtrap suppressions` and `mailtrap tracking-opt-outs` commands.

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

## domains update

Update the tracking and inbound settings of a sending domain.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Domain ID |
| `--open-tracking` | bool | No | Track opens on emails sent from this domain |
| `--click-tracking` | bool | No | Track clicks on links in emails sent from this domain |
| `--tracking-opt-out` | bool | No | Add the tracking opt-out link to tracked emails; requires open or click tracking |
| `--auto-unsubscribe-link` | bool | No | Automatically add an unsubscribe link to emails |
| `--inbound-enabled` | bool | No | Allow the domain to be attached to an inbound inbox as a catch-all |

Only the flags actually passed are sent, so a partial update leaves the other settings alone. Pass `--flag=false` to turn a setting off.

---

## domains delete

Delete a sending domain.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Domain ID |

---

## company-info get

Retrieve the company info of a sending domain, used for domain compliance verification.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--domain-id` | string | Yes | Sending domain ID |

**Note:** Uses the API token's account; `--account-id` is not needed.

---

## company-info create

Set the company info of a sending domain.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--domain-id` | string | Yes | Sending domain ID |
| `--name` | string | Yes | Company or individual name |
| `--address` | string | Yes | Street address |
| `--city` | string | Yes | City |
| `--country` | string | Yes | Country |
| `--zip-code` | string | Yes | ZIP or postal code |
| `--website-url` | string | Yes | Company website URL |
| `--phone` | string | No | Phone number |
| `--privacy-policy-url` | string | No | URL to the privacy policy page |
| `--terms-of-service-url` | string | No | URL to the terms of service page |
| `--info-level` | string | No | Whether the sender is a `business` or an `individual` |

---

## company-info update

Change the company info of a sending domain.

Takes the same flags as `company-info create`, all optional except `--domain-id`. Only the flags actually passed are sent, so a partial update leaves the other fields alone. An update with no attribute flags is rejected rather than sent as an empty payload.

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
