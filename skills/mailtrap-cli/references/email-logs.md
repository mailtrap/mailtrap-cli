# email-logs

Detailed flag specifications for `mailtrap email-logs` and `mailtrap stats` commands.

---

## email-logs list

List email logs (sent email history).

No additional flags. Returns recent email logs for the account.

**Output:** Table/JSON of email logs with ID, to, subject, status, and timestamp.

---

## email-logs get

Get a specific email log entry.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Email log ID |

---

## stats get

Get aggregated email sending statistics.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--start-date` | string | Yes | Start date (e.g. `2024-01-01`) |
| `--end-date` | string | Yes | End date (e.g. `2024-01-31`) |
| `--domain-ids` | string[] | No | Filter by domain IDs, can be repeated |
| `--streams` | string[] | No | Filter by streams, can be repeated |
| `--categories` | string[] | No | Filter by categories, can be repeated |

**Example:**
```bash
mailtrap stats get \
  --start-date 2024-01-01 \
  --end-date 2024-01-31 \
  --output json
```

---

## stats by-domain

Get statistics grouped by sending domain.

Same flags as `stats get`.

---

## stats by-category

Get statistics grouped by email category.

Same flags as `stats get`.

---

## stats by-esp

Get statistics grouped by email service provider.

Same flags as `stats get`.

---

## stats by-date

Get statistics grouped by date.

Same flags as `stats get`.
