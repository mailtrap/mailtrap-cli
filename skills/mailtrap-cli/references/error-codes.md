# error-codes

Common API errors and their meanings when using the Mailtrap CLI.

---

## HTTP Status Codes

| Code | Meaning | Common Cause |
|------|---------|--------------|
| 400 | Bad Request | Invalid flag values, malformed JSON in `--file`, missing required fields |
| 401 | Unauthorized | Invalid or expired API token |
| 403 | Forbidden | Token lacks permissions for this action |
| 404 | Not Found | Invalid resource ID (domain, template, inbox, message, etc.) |
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

### Resource not found
```
Error: API error 404: Not Found
```
**Fix:** Verify the resource ID exists. Use the corresponding `list` command to find valid IDs.

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
