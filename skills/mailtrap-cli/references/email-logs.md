# email-logs

Detailed flag specifications for `mailtrap email-logs` and `mailtrap stats` commands.

---

## email-logs list

List email logs (sent email history).

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--cursor` | string | No | Pagination cursor: `next_page_cursor` from the previous response |
| `--sent-after` | string | No | Only logs sent after this ISO 8601 timestamp |
| `--sent-before` | string | No | Only logs sent before this ISO 8601 timestamp |
| `--to` | string | No | Filter by recipient email |
| `--to-operator` | string | No | Operator for `--to`: `ci_equal` (default), `ci_not_equal`, `ci_contain`, `ci_not_contain` |
| `--from` | string | No | Filter by sender email |
| `--from-operator` | string | No | Operator for `--from`, same values as `--to-operator` |
| `--subject` | string | No | Filter by subject |
| `--subject-operator` | string | No | Operator for `--subject`, same values as `--to-operator` |
| `--status` | string | No | Filter by status: `delivered`, `not_delivered`, `enqueued`, `opted_out` |
| `--event` | string | No | Filter by event: `delivery`, `open`, `click`, `bounce`, `spam`, `unsubscribe` |
| `--category` | string | No | Filter by category |

**Output:** Table/JSON of email logs with ID, to, subject, status, and timestamp. In table and text output the next-page cursor is printed as `--cursor <value>` when more logs are available.

**Example:**
```bash
mailtrap email-logs list \
  --sent-after 2024-01-01T00:00:00Z \
  --to "user@example.com" \
  --status delivered
```

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
