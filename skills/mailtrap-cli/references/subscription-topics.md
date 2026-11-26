# subscription-topics

Detailed flag specifications for `mailtrap subscription-topics` commands, and the subscription topic flags of other groups.

Subscription topics let contacts unsubscribe from one kind of marketing email (for example "Product updates") instead of from everything. Topics are created and edited in the Mailtrap app; the API and the CLI only read them. A topic covers one domain (`domain_id`), or every domain when `domain_id` is `null`.

`subscription-topics` commands are token-scoped: they take no `--account-id`. On an account without subscription topics they fail with `API error 403`.

---

## subscription-topics list

List the account's subscription topics, one page at a time. Deleted topics are left out.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--per-page` | int | No | Number of topics per page (default 50, max 100) |
| `--token` | int | No | Page number to retrieve (page-token pagination) |

**Output:** Table of topics with ID, name, description, domain ID (empty for every domain) and creation time, followed by `Next page: --token N` when more pages exist (`Next page: --token N --per-page M` when `--per-page` was set). With `--output json`, the full response object is printed: `.data` holds the topics and `.pagination.next_token` the next page (`null` on the last page).

---

## subscription-topics get

Get a subscription topic. Deleted topics are returned too, with `deleted_at` set (DELETED AT column).

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | int | Yes | Subscription topic ID |

---

## Topics on campaigns

`email-campaigns create` and `email-campaigns update` take `--subscription-topic-id <id>`: recipients who unsubscribed from that topic don't get the campaign. The topic must cover the campaign's domain. Responses carry `subscription_topic_id` (TOPIC ID column), `null` when the campaign has none.

```bash
mailtrap email-campaigns update --id 4567 --subscription-topic-id 7
```

A campaign's topic can be changed but not removed. On an account without subscription topics the flag is ignored and the campaign reads back with `subscription_topic_id: null`.

---

## Topics on contacts

`contacts create` and `contacts update` take `--subscription-topics`, a JSON array with one entry per topic and domain. See [contacts.md](contacts.md).

```bash
mailtrap contacts update --id <CONTACT_ID> --email user@example.com \
  --subscription-topics '[{"subscription_topic_id":7,"domain_id":4321,"status":"unsubscribed"}]'
```

`contacts get`, `create` and `update` return the contact's preferences in `subscription_topics` (JSON output only): a list of `{subscription_topic_id, domain_id, status}`, `null` when they couldn't be read (retry), or no key at all on an account without subscription topics.

---

## Topic unsubscriptions

`suppressions list --include-topic-unsubscriptions` adds the contacts' unsubscriptions from single topics, with type `topic unsubscription` and the topic in `subscription_topic_id` (TOPIC_ID column). They are left out, with no error, unless the API token can read contacts. See [domains.md](domains.md).
