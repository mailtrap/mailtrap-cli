# email-campaigns

Detailed flag specifications for `mailtrap email-campaigns` commands.

`email-campaigns` commands are token-scoped: they take no `--account-id` (requests go to `/api/email_campaigns`).

A campaign starts as a draft. Only a draft can be updated, deleted, started or scheduled; a scheduled campaign can be cancelled or reset to draft, and a sending one terminated.

**Output:** Campaign commands print a table with ID, name, state, domain, recipient count, topic ID and creation time. With `--output json`, the campaign as the API returns it.

---

## email-campaigns list

List the account's email campaigns, one page at a time.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--per-page` | int | No | Number of campaigns per page (default 50, max 100) |
| `--search` | string | No | Filter campaigns by name |
| `--token` | int | No | Page number to retrieve (page-token pagination) |

**Output:** Table of campaigns, followed by `Next page: --token N` when more pages exist; repeat `--search` and `--per-page` with it. With `--output json`, the full response object is printed: `.data` holds the campaigns and `.pagination.next_token` the next page (`null` on the last page).

---

## email-campaigns get

Get an email campaign.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Email campaign ID |

---

## email-campaigns create

Create a draft email campaign.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | Yes | Campaign name |
| `--domain-id` | int | Yes | ID of the verified domain to send from (see `domains list`) |
| `--from-local-part` | string | Yes | Local part (before the @) of the From address |
| `--subject` | string | Yes | Email subject line |
| `--from-display-name` | string | No | Display name shown in the From header |
| `--reply-to-display-name` | string | No | Reply-To display name |
| `--reply-to-local-part` | string | No | Reply-To local part (before the @) |
| `--reply-to-domain` | string | No | Reply-To domain |
| `--body-html` | string | No | HTML body of the email (the design) |
| `--body-text` | string | No | Plain-text alternative of the email body |
| `--merge-tags` | string[] | No | Merge tag names used in the content (comma-separated), without `{{ }}` |
| `--delivery-mode` | string | No | Delivery mode: `rapid`, `gradual` |
| `--emails-per-hour` | int | No | Emails per hour when the delivery mode is `gradual` |
| `--contact-list-ids` | int[] | No | Contact list IDs to send to (comma-separated) |
| `--contact-segment-ids` | int[] | No | Contact segment IDs to send to (comma-separated) |
| `--subscription-topic-id` | int | No | Subscription topic ID; see [subscription-topics.md](subscription-topics.md) |

**Example:**
```bash
mailtrap email-campaigns create --name "Spring Sale" --domain-id 4321 \
  --from-local-part news --subject "Spring is here" \
  --body-html '<h1>Hi {{first_name}}!</h1><a href="__unsubscribe_url__">Unsubscribe</a>' \
  --merge-tags first_name --contact-list-ids 55,56
```

---

## email-campaigns update

Update a draft email campaign. Takes the same flags as `create`, all optional, plus:

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Email campaign ID |
| `--clear-contact-lists` | bool | No | Remove all contact lists from the audience (not with `--contact-list-ids`) |
| `--clear-contact-segments` | bool | No | Remove all contact segments from the audience (not with `--contact-segment-ids`) |

Only supplied flags are updated; at least one is required. `--contact-list-ids` and `--contact-segment-ids` replace the whole audience list rather than adding to it.

```bash
mailtrap email-campaigns update --id 4567 --subject "Spring is finally here" --contact-list-ids 55,56
```

---

## email-campaigns delete

Delete a draft email campaign.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Email campaign ID |

**Output:** `Email campaign deleted successfully.`, also with `--output json`.

---

## email-campaigns start

Start sending a draft email campaign now.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Email campaign ID |

---

## email-campaigns schedule

Schedule a draft email campaign.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Email campaign ID |
| `--datetime` | string | Yes | When to send the campaign, ISO 8601 (e.g. `2030-01-01T09:00:00Z`) |

---

## email-campaigns cancel

Cancel a scheduled email campaign.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Email campaign ID |

---

## email-campaigns reset

Reset a scheduled email campaign to draft.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Email campaign ID |

---

## email-campaigns terminate

Terminate a campaign that is sending.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Email campaign ID |

---

## email-campaigns stats

Get an email campaign's statistics: sent, delivered, opens, clicks, bounces, spam complaints and unsubscribes, with their rates.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Email campaign ID |
| `--start-date` | string | No | Start of the aggregation window (`YYYY-MM-DD`) |
| `--end-date` | string | No | End of the aggregation window (`YYYY-MM-DD`) |

**Output:** The statistics, not the campaign: counts and rates in a table, or the response's `data` object with `--output json`.
