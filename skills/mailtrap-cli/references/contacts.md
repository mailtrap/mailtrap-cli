# contacts

Detailed flag specifications for `mailtrap contacts`, `mailtrap contact-lists`, and `mailtrap contact-fields` commands.

---

## contacts get

Get a specific contact.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact ID |

---

## contacts create

Create a new contact.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--email` | string | Yes | Contact email address |
| `--first-name` | string | No | First name |
| `--last-name` | string | No | Last name |
| `--list-ids` | int[] | No | List IDs to add the contact to |

---

## contacts update

Update an existing contact.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact ID |
| `--email` | string | No | New email address |
| `--first-name` | string | No | New first name |
| `--last-name` | string | No | New last name |

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
| `--list-id` | int | No | List ID to import contacts into |

---

## contacts export

Export contacts.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--list-id` | int | No | List ID to export from |

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
| `--type` | string | Yes | Event type name |
| `--data` | string | No | Event data as JSON string |

---

## contact-lists list

List all contact lists.

No additional flags.

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
| `--data-type` | string | No | New data type |
| `--merge-tag` | string | No | New merge tag |

---

## contact-fields delete

Delete a custom contact field.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Contact field ID |
