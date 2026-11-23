# error-codes

Common API errors and their meanings when using the Mailtrap CLI.

---

## HTTP Status Codes

| Code | Meaning | Common Cause |
|------|---------|--------------|
| 400 | Bad Request | Invalid flag values, malformed JSON in `--file`, missing required fields |
| 401 | Unauthorized | Invalid or expired API token |
| 403 | Forbidden | Token lacks permissions for this action |
| 404 | Not Found | Invalid resource ID (domain, template, sandbox, message, etc.) |
| 422 | Unprocessable Entity | Validation error (e.g. duplicate domain, invalid email format) |
| 429 | Too Many Requests | Rate limit exceeded — wait and retry |
| 500 | Internal Server Error | Mailtrap server issue — retry later |

---

## CLI Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (API error, missing flags, invalid input) |

---

## API Error Format

A failed request prints one line:

```
Error: API error <status>: <message>
```

The message comes from the response body:

| Body | Message |
|------|---------|
| `{"error": "Not Found"}` | `Not Found` |
| `{"errors": "Access forbidden"}` | `Access forbidden` |
| `{"errors": ["'from' is required", "'subject' is required"]}` | `'from' is required; 'subject' is required` |
| `{"errors": {"name": ["can't be blank"], "email": ["is invalid"]}}` | `email is invalid; name can't be blank` (one `field message` pair per message, sorted by field) |
| `{"errors": {"base": ["Contacts limit reached"]}}` | `Contacts limit reached` (`base` errors belong to no field) |
| empty | nothing: `API error 403` |
| not JSON | the body as is |

---

## Common Error Patterns

### Missing required flag
```
Error: required flag(s) "--from", "--to", "--subject" not set
```
**Fix:** Supply all required flags for the command.

### Invalid API token
```
Error: API error 401: Unauthorized
```
**Fix:** Run `mailtrap configure --api-token <valid-token>` or set `MAILTRAP_API_TOKEN`.

### Missing permission
```
Error: API error 403: Access forbidden
```
**Fix:** Use a token with access to the resource, or ask an account admin to grant it.

### Resource not found
```
Error: API error 404: Not Found
```
**Fix:** Verify the resource ID exists. Use the corresponding `list` command to find valid IDs.

### Validation failed
```
Error: API error 422: email is invalid
```
**Fix:** Correct the named field and retry. Several invalid fields are joined with `; `.

### Nothing to update
```
Error: at least one attribute flag is required
```
**Fix:** Commands such as `templates update` and `contacts update` send only the flags you pass; pass at least one besides `--id`.

### Domain not verified
```
Error: API error 422: Domain is not verified
```
**Fix:** Complete DNS verification in the Mailtrap web dashboard before sending.

### Rate limited
```
Error: API error 429: Too Many Requests
```
**Fix:** Wait before retrying. Consider using batch endpoints for high-volume sends.

### Invalid batch JSON
```
Error: failed to read file: invalid JSON
```
**Fix:** Validate the JSON file format. Use `jq . < batch.json` to check syntax.
