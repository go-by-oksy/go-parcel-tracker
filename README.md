# Go Parcel Tracker

A parcel tracking application built with Go and SQLite.

The project implements a storage layer for registering parcels, updating their addresses and statuses, retrieving parcel information, and deleting parcels according to their current lifecycle state.

## Features

- Register new parcels
- Retrieve a parcel by tracking number
- Retrieve all parcels belonging to a client
- Update the delivery address of registered parcels
- Change parcel status
- Delete parcels only while they are registered
- Store data in SQLite
- Create the database schema automatically
- Use isolated temporary databases in tests
- Run automated formatting, testing and build checks with GitHub Actions

## Parcel Lifecycle

A parcel can have one of the following statuses:

- `registered`
- `sent`
- `delivered`

The delivery address can be changed only while the parcel has the `registered` status.

A parcel can also be deleted only while it is still registered.

## Tech Stack

- Go
- SQLite
- `database/sql`
- `modernc.org/sqlite`
- Unit testing
- Testify
- GitHub Actions

## Project Structure

```text
.
├── .github/workflows   # CI configuration
├── main.go             # Example application workflow
├── parcel.go           # Parcel model and SQLite storage
├── parcel_test.go      # Storage tests
├── go.mod
└── README.md
```

## Running Locally

### Requirements

- Go installed
- Git

Clone the repository:

```bash
git clone https://github.com/go-by-oksy/go-parcel-tracker.git
cd go-parcel-tracker
```

Download dependencies:

```bash
go mod download
```

Run the application:

```bash
go run .
```

The application automatically creates a local SQLite database named `tracker.db`.

The database file is excluded from version control.

## Testing

Run all tests:

```bash
go test ./...
```

Each test uses its own temporary SQLite database, so tests do not depend on previously stored data.

Check that the application builds:

```bash
go build ./...
```

## Project Background

This project was completed as part of the Yandex Practicum Go development course and is based on the provided project requirements.

The SQLite storage implementation, database initialization, isolated tests, query optimization and subsequent portfolio improvements were completed by [Oksana](https://github.com/go-by-oksy).