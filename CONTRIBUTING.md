# Contributing to Lockify

Thank you for your interest in contributing to Lockify!  
This document explains how to set up your environment, submit changes, and follow project conventions.

---

## Requirements

- Go 1.24+
- Just (optional: https://github.com/casey/just)
- Git
- A GitHub account

---

## Getting Started

### 1. Fork the Repository

```
https://github.com/apixify/lockify
```

### 2. Clone Your Fork

```sh
git clone https://github.com/<your-username>/lockify.git
cd lockify
```

### 3. Install Dependencies

```sh
go mod tidy
```

---

## Branching Model

Use the following naming conventions:

- `feature/lock-<issue-no>-<feature-name>` — new features  
- `fix/lock-<issue-no>-<bug-name>` — bug fixes  
- `chore/lock-<issue-no>-<task>` — maintenance work  
- `docs/lock-<issue-no>-<topic>` — documentation  

Example:

```sh
git checkout -b feature/lock-1-add-keyring-support
```

---

## Running Tests

```sh
go test ./... -v
```

---

## Code Style

- Follow idiomatic Go practices.
- Error handling should use wrapped errors:

```go
return fmt.Errorf("failed to read file: %w", err)
```

- Avoid long functions; keep logic in domain or use case layers.
- All exported functions must have comments.

---

## Commit Messages

Use conventional commits:

- `[lock-<issue-no>] add passphrase rotation`
- `e.g. [lock-1] Nice descriptive commit`

---

## Submitting Pull Requests

1. Ensure tests pass with `go test ./...`.
2. Rebase onto latest `main`.
3. Open a PR with:
   - Clear description  
   - Screenshots (if relevant)  
   - Steps to test  

PRs will be reviewed for:

- Correctness  
- Security impact  
- Code quality  
- Architecture alignment  

---

## Reporting Issues

Open an issue for:

- Bugs
- Feature requests
- Security concerns (or email the security contact listed in README)

Thanks for helping improve Lockify!
