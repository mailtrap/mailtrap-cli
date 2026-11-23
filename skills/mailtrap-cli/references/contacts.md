# contacts

Detailed flag specifications for `mailtrap contacts`, `mailtrap contact-lists`, and `mailtrap contact-fields` commands.

---

## contacts get

Get a specific contact.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact ID |

**Output:** With `--output json`, the contact as the API returns it. `subscription_topics` holds the contact's subscription topic preferences, one `{subscription_topic_id, domain_id, status}` entry per topic and domain; it is `null` when they couldn't be read (retry), and absent on accounts without subscription topics. `create` and `update` return the same shape.

---

## contacts create

Create a new contact.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--email` | string | Yes | Contact email address |
| `--first-name` | string | No | First name |
| `--last-name` | string | No | Last name |
| `--list-ids` | int[] | No | List IDs to add the contact to |
| `--subscription-topics` | string | No | Subscription topic preferences as a JSON array (see below) |

`--subscription-topics` takes one entry per topic and domain: `subscription_topic_id` and `domain_id` (both required, from `subscription-topics list` and `domains list`) and `status` (`subscribed` or `unsubscribed`). Topic and domain pairs left out keep their status. On accounts without subscription topics the API rejects the flag with `422` ("are not available for this account"). If the contact is saved but its preferences are not, the API returns `503`; send them again with `contacts update`.

```bash
mailtrap contacts create --email user@example.com \
  --subscription-topics '[{"subscription_topic_id":7,"domain_id":4321,"status":"unsubscribed"}]'
```

---

## contacts update

Update an existing contact.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact ID |
| `--email` | string | No | New email address |
| `--first-name` | string | No | New first name |
| `--last-name` | string | No | New last name |
| `--list-ids-included` | int[] | No | List IDs to add the contact to |
| `--list-ids-excluded` | int[] | No | List IDs to remove the contact from |
| `--unsubscribed` | bool | No | Unsubscribe the contact (`--unsubscribed=false` subscribes it again) |
| `--subscription-topics` | string | No | Subscription topic preferences as a JSON array, as for `contacts create`; pairs left out keep their status |

Only supplied flags are updated; omitted fields remain unchanged. At least one flag besides `--id` is required, so a topic-only update needs just `--id` and `--subscription-topics`.

---

## contacts delete

Delete a contact.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact ID |

---

## contacts import

Import contacts from a JSON file.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--file` | string | Yes | Path to JSON file with contact data |

The file is posted as is. Lists are set per contact in the file: `list_ids_included` adds the contact to lists and `list_ids_excluded` removes it from them. Each contact may also carry `subscription_topics`, with the same entries as `contacts create --subscription-topics`:

```json
{"contacts": [{"email": "user@example.com", "list_ids_included": [55], "subscription_topics": [{"subscription_topic_id": 7, "domain_id": 4321, "status": "unsubscribed"}]}]}
```

---

## contacts export

Export contacts.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--list-ids` | int[] | No | Export only contacts in these lists |
| `--status` | string | No | Export only contacts with this subscription status: `subscribed`, `unsubscribed` |

---

## contacts import-status

Check the status of a contact import job.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Import job ID |

---

## contacts export-status

Check the status of a contact export job.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Export job ID |

---

## contacts create-event

Create a custom event for a contact (used for automation triggers).

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact ID |
| `--name` | string | Yes | Event name (max 255 characters) |
| `--params` | string | No | Event params as a JSON object |

---

## contact-lists list

List all contact lists.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--search` | string | No | Filter by name (case-insensitive prefix match) |

---

## contact-lists get

Get a specific contact list.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact list ID |

---

## contact-lists create

Create a new contact list.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | Yes | List name |

---

## contact-lists update

Update a contact list.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact list ID |
| `--name` | string | Yes | New list name |

---

## contact-lists delete

Delete a contact list.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact list ID |

---

## contact-fields list

List all custom contact fields.

No additional flags.

---

## contact-fields get

Get a specific contact field.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact field ID |

---

## contact-fields create

Create a custom contact field.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | Yes | Field name |
| `--data-type` | string | Yes | Data type: `text`, `integer`, `float`, `boolean`, `date` |
| `--merge-tag` | string | Yes | Merge tag for the field (e.g. `{{company}}`) |

---

## contact-fields update

Update a custom contact field.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact field ID |
| `--name` | string | No | New field name |
| `--merge-tag` | string | No | New merge tag |

The data type can't be changed.

---

## contact-fields delete

Delete a custom contact field.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact field ID |
