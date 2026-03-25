# domains

Detailed flag specifications for `mailtrap domains` and `mailtrap suppressions` commands.

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

No additional flags.

---

## suppressions delete

Remove an address from the suppression list.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Suppression ID |
