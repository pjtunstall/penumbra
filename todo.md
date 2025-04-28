# Todo

## Security

- Sanitize all inputs and restrict their size.
  - In particular, prevent passwords larget than 72 bytes, bycrypt's limit. (See the comment in `SubmitRegister` in `handlers.go`).
- Switch to gorilla/mux for a simple way to do more secure route parsing rather than just using `TrimPrefix` to extract ids. I'm parsing the suffix to an int; that's some validation, but consider risks associated with malicious routes.
- Implement rate limiting.
- Limit number of users and number of tasks per user.
- Set cookie's `Secure` field value to `true` in production.
- Implement password reset in production.
- Likewise 2FA.

## Tests

- Finish writing tests for handlers, router, and database store methods. Consider what else can be tested.
- Resolve any inconsistencies in testing style, e.g. when to use testify. Less dependencies is good, but testify might make the code clearer due to its more declarative style.

## Error handling

- Ensure that error handling is consistent.
- Consider when to panic and what to log, and in what format.
- Notify user if their username or password is incorrect.
- Have an error page template to gracefully display error messages that the user in the name.
- Add more checking around task status to ensure done is converted correctly and never set to an anomalous value.

## Naming

- Be more consistent about names, e.g. `task` versus `tasks`. Consider especially at names of routes (and look up conventions regarding these) and names of handler functions: be more consistent, e.g. about when "Task" is included in the name.

## Refactor

- See what can be simplified.
- At the moment, handlers are wrapped in `HandleProtected` only if they don't have another argument besides the response writer and request; otherwise, they check authorization themselves. This could be done more consistently or at least without the repetition.
- Split up large files.
