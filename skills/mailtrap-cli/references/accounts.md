# accounts

Detailed flag specifications for account management commands: `mailtrap accounts`,
`mailtrap account-access`, `mailtrap tokens`, `mailtrap permissions`,
`mailtrap billing`, and `mailtrap organizations`.

---

## accounts list

List all accounts accessible with the current API token.

No additional flags.

---

## account-access list

List all account access entries (users/tokens with access).

No additional flags.

---

## account-access remove

Remove access for a user or token.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Account access ID |

---

## tokens list

List all API tokens for the account.

No additional flags.

---

## tokens get

Get details of a specific API token.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Token ID |

---

## tokens create

Create a new API token.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | Yes | Token name |
| `--permissions` | string | Yes | Permissions JSON array, e.g. `'[{"resource_type":"account","resource_id":123,"access_level":100}]'` |

**Note:** The token value is shown only once in the response. Store it securely.

---

## tokens delete

Delete an API token.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Token ID |

---

## tokens reset

Reset (regenerate) an API token.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | string | Yes | Token ID |

**Note:** The new token value is shown only once. The old token stops working immediately.

---

## permissions bulk-update

Bulk update permissions for an account access entry.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--access-id` | string | Yes | Account access ID |
| `--file` | string | Yes | Path to JSON file with permissions data |

---

## permissions resources

List all available permission resources and their types.

No additional flags.

---

## billing usage

Get current billing usage for the account.

No additional flags.

---

## organizations list-sub-accounts

List sub-accounts for an organization.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--org-id` | string | Yes | Organization ID |

---

## organizations create-sub-account

Create a sub-account under an organization.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--org-id` | string | Yes | Organization ID |
| `--name` | string | Yes | Sub-account name |
