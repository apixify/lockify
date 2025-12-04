# Lockify

> 🔐 Lockify — a secure, developer-friendly CLI for managing encrypted environment variables across environments (dev, staging, prod).  
> Safe to store in Git. CI/CD friendly. Built in Go.

---

## Why Lockify?

Lockify solves a simple, common problem: how to keep environment variables secure, sharable, and scriptable without relying on third-party cloud vendors.  
It focuses on an offline-first, CLI-first experience with strong cryptography and a minimal mental model.

- Store encrypted vaults per environment  
- Export decrypted envs in CI using a secret passphrase  
- Designed for developers and small teams who prefer Git-native workflows  

---

## Key Features

- **AES-256-GCM Encryption** (authenticated encryption)  
- **Argon2id KDF** for deriving encryption keys  
- **Passphrase caching** via OS keyring (optional)  
- **Multi-environment vaults** (dev, staging, prod, …)  
- **Import/export** `.env` and JSON formats  
- **Key rotation** without losing data  
- Clean, testable codebase using DDD and clean architecture  

---

## Install

### Using Go

```sh
go install github.com/apixify/lockify@latest
```

### From Source

```sh
git clone https://github.com/apixify/lockify.git
cd lockify
go build -o lockify .
```

---

## Quick Start

### 1. Initialize a Vault

```sh
lockify init --env prod
```

### 2. Add a Secret

```sh
lockify add --env prod --secret
```

### 3. Export to `.env` (CI-friendly)

```sh
lockify export --env prod --format dotenv > .env
```

### 4. Get a Value

```sh
lockify get --env prod --key DATABASE_URL
```

---

## GitHub Actions Example

```yaml
steps:
  - uses: actions/checkout@v4

  - name: Install Lockify
    run: go install github.com/apixify/lockify@latest

  - name: Export env vars
    env:
      LOCKIFY_PASSPHRASE: ${{ secrets.LOCKIFY_PASSPHRASE }}
    run: lockify export --env prod --format dotenv > .env
```

---

## Security Summary

- Vault files **can be committed to Git** (fully encrypted).  
- Passphrases are **never** stored in plaintext.  
- Optional passphrase caching uses the **OS keyring**.  
- Rotate passphrases using:

```sh
lockify rotate-key --env <env>
```

---

## Project Layout (High Level)

```
cmd/          # CLI commands (Cobra)
internal/
  domain/     # Entities and interfaces (pure)
  app/        # Use cases
  infra/      # Crypto, filesystem, keyring implementations
  ui/         # Prompts and output formatting
  di/         # Command-scoped dependency injection
docs/         # Documentation
```

---

## Contributing

Contributions are welcome!  
Please see `CONTRIBUTING.md` for guidelines.

---

## License

MIT — see `LICENSE`.

---

## Security Contact

If you discover a vulnerability, contact:  
**ahmed.elkayaty92@gmail.com**
