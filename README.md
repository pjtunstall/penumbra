# PENUMBRA

- [Overview](#overview)
- [Project status](#project-status)

## Overview

PENUMBRA (Planning & Execution Nexus for Urgent Management, Briefing & Recording App) is a new web application for caseworkers to manage their tasks. Use of PENUMBRA is not mandatory, but may become so for any caseworker who elects not to use it.

## Project status

This project should be considered inclomplete if any `todo` remains in the actual code. Any CI/Cd pipeline should enforce that.

## Usage

To initialize a database, compile the dbinit binary with `go build -o dbinit cmd/dbinit/main.go` and run it `./dbinit` (or the equivalent command for your operating system). This will initialize a database called `dev.db` in a newly created `data` directory in the root of this project.

Then to build and run the app, run `go run cmd/webapp/main.go` (assuming your working directory is the project root).

To run all tests, run `go test ./...`.
