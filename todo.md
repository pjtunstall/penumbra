# Todo

## Features

- Have client notify the user when a deadline is close.
- Cache fetched tasks in memory and update status when a deadline passes. Let client be the one to derive whether a task is overdue; currently the server computes it from due data and current time before sending task data to the client. Be sure to synchronize time between server and client.
- Create an `openapi.yaml` file. If asking AI help with this, make sure it does correcpond to the code; in particular, make sure that it reflects the fact that the server returns HTML, not JSON. (At some point, experiment with using `openapi-generator` to generate client libraries.)
- Add ability to uncheck a task in case it was accidentally marked as done.

## Error handling

- Have an error page template to gracefully display error messages that the user in the name.
- Ensure that error handling is consistent.
- Consider when to panic and what to log, and in what format.
- Notify user if their username or password is incorrect or too long on registration.
- Add more checking around task status to ensure done is converted correctly and never set to an anomalous value.

## Security

- There should be security tests as well as basic functionality tests.
- Check that all inputs are sanitized and restrict their size. SQL-injection is prevented, so just make sure no input is inserted directly into the HTML.
  - Limit the size of task names and descriptions, user names, and emails. Verify email format. In production, emails would also need verifying by sending a confirmation code.
- Switch to gorilla/mux (or chi?) for a simple way to do more secure route parsing rather than just using `TrimPrefix` to extract ids. I'm parsing the suffix to an int; that's some validation, but consider risks associated with malicious routes.
- Implement rate limiting.
- Limit number of users and number of tasks per user.
- Set cookie's `Secure` field value to `true` in production.
- Implement password reset in production.
- Likewise 2FA.

## Tests

- Finish writing tests for handlers, router, and database store methods. Consider what else can be tested.
- Watch out for inconsistencies in testing style, e.g. when to use testify. Less dependencies is good, but testify might make the code clearer due to its more declarative style.

## Naming

- Use a prefix like `/v1/` to indicate versioning.
- Capitalize ID in Go; there are no instances of UUID. But be careful not to corrupt existing UUIDs (in tests) and hashes (of dependencies).

## Refactor

- See what can be simplified.
- Split up large files.

## Cosmetic

- Put a little bit more vertical gap between form fields.
