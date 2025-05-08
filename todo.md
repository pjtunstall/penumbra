# Todo

## Fix

- New tasks are now being created at the start of the day, It should be the end.

## Features

- Add an `openapi.yaml` file.
- Add ability to uncheck a task in case it was accidentally marked as done.
- Have client notify the user when a deadline is close.
- Cache fetched tasks in memory and update status when a deadline passes. Let client be the one to derive whether a task is overdue; currently the server computes it from due data and current time before sending task data to the client. Be sure to synchronize time between server and client.

## Security

- Switch to using UUIDs for task and user ids, rather than incrementing an int.
- Install Tailwind and DaisyUI locally. Write a small custom CSS entry file. Run the Tailwind compiler (with DaisyUI plugin) to generate a single `.css` file. Serve that file from your server (covered by 'self' in the CSP).
- Check that all inputs are sanitized and restrict their size. They go safely into the databse, so just make sure no input is inserted directly into the HTML.
  - Limit the size of task names and descriptions, user names, and emails. Verify email format. In production, emails would also need verifying by sending a confirmation code.
- Switch to gorilla/mux (or chi?) for a simple way to do more secure route parsing rather than just using `TrimPrefix` to extract ids. I'm parsing the suffix to an int; that's some validation, but consider risks associated with malicious routes.
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
- Notify user if their username or password is incorrect, or too long on registration.
- Have an error page template to gracefully display error messages that the user in the name.
- Add more checking around task status to ensure done is converted correctly and never set to an anomalous value.

## Naming

- Use a prefix like `/v1/` to indicate versioning.
- Review names of handlers for consistency and clarity.

## Refactor

- See what can be simplified.
- At the moment, handlers are wrapped in `HandleProtected` only if they don't have another argument besides the response writer and request; otherwise, they check authorization themselves. This could be done more consistently or at least without the repetition.
- Split up large files.

## Cosmetic

- Put a little bit more vertical gap between form fields.
