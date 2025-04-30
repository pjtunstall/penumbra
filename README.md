# PENUMBRA

- [Overview](#overview)
- [Project status](#project-status)

## Overview

PENUMBRA (Planning & Execution Nexus for Urgent Management, Briefing & Recording App) is a new web application for caseworkers to manage their tasks. Use of PENUMBRA is not mandatory, but may become so for any caseworker who elects not to use it.

## Project status

This project should be considered inclomplete if any `todo` remains in the actual code. Any CI/Cd pipeline should enforce that.

## Usage

Dowload and install the [Go programming language](https://go.dev/doc/install) if you haven't already.

To initialize a database, compile the dbinit binary with `go build -o dbinit cmd/dbinit/main.go` and run it `./dbinit` (or the equivalent command for your operating system). This will initialize a database called `dev.db` in a newly created `data` directory in the root of this project.

Then to build and run the app in one step, run `go run cmd/webapp/main.go` (assuming your working directory is the project root). Open a web browser and navigate to `http://localhost:8080`.

To run all tests, run `go test ./...`.

### Routes

- `GET  /` – redirect to `/login` if not logged in, otherwise redirect to `/dashboard`
- `GET  /login` – show login form
- `POST /login` – submit login form
- `GET /register` - show register form
- `POST /register` - submit register form
- `GET /dashboard` - show dashboard, listing any task titles, due datss, and status, with the option to mark them as done
- `GET /about` - show about page
- `GET /logout` - log out and redirect to `/login`
- `GET /tasks` - list all tasks for the current user, including descriptions, due dates, and status
- `GET /tasks/create` - show form to create new task
- `POST /tasks/create` - submit form to create new task
- `GET /task/{id}` - show task details, including description, due date, and status in a form that allows the task to be deleted, or edited and updated.
- `POST /task/delete/{id}` - delete task
- `POST /task/done/{id}` - mark task as done
- `POST /task/update/{id}` - submit form to update task

Chat remarks:

> You're following a classic HTML form-based pattern, and it's perfectly fine for a traditional server-rendered app. REST purists would nudge you toward resource-based paths and HTTP verbs, but in practice for web apps, what you're doing is conventional and user-friendly.
>
> That said, here's how you could tweak your route names to better align with REST conventions while still supporting HTML forms:
>
> | Purpose                    | Current Route           | Suggested REST-style Route                                | Notes                                                             |
> | -------------------------- | ----------------------- | --------------------------------------------------------- | ----------------------------------------------------------------- |
> | Show empty "new task" form | GET /tasks/create       | GET /tasks/new                                            | `new` is a common REST-y way to signal a form for creation        |
> | Submit new task form       | POST /tasks/create      | POST /tasks                                               | REST uses POST to the collection to create                        |
> | Show edit task form        | GET /tasks/{id}         | GET /tasks/{id}/edit                                      | `edit` is a conventional form path in REST-style HTML apps        |
> | Submit task update         | POST /tasks/update      | POST /tasks/{id} or /{id}/update                          | Use POST with clear naming since browsers don’t support PUT/PATCH |
> | Delete/update UI           | Part of GET /tasks/{id} | Consider GET /tasks/{id}/edit and POST /tasks/{id}/delete | Split views/actions for clarity                                   |
