# Lockify

🔐 Lockify — a secure, developer-friendly CLI for managing encrypted environment variables across environments (dev, staging, prod).  
Safe to store in Git. CI/CD friendly. Built in Go.

---

## Table of Contents

- [Why Lockify?](#why-lockify)
- [Key Features](#key-features)
- [Install](#install)
  - [Using Homebrew](#using-homebrew)
  - [Download Pre-built Binaries](#download-pre-built-binaries)
  - [Using Go](#using-go)
  - [From Source](#from-source)
- [Quick Start](#quick-start)
- [GitHub Actions Example](#github-actions-example)
- [Security Summary](#security-summary)
- [Contributing](#contributing)
- [License](#license)
- [Security Contact](#security-contact)

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

### Using Homebrew

```sh
brew tap ahmed-abdelgawad92/lockify
brew install lockify
```

### Download Pre-built Binaries

Visit the [GitHub Releases](https://github.com/ahmed-abdelgawad92/lockify/releases/latest) page to download pre-built binaries for your platform.

**Linux:**
#### Download the appropriate .tar.gz file for your architecture (amd64 or arm64)
```sh
# Then extract and install:
tar -xzf lockify_1.2.3_linux_amd64.tar.gz
sudo mv lockify /usr/local/bin/
```

**macOS:**
#### Download the appropriate .tar.gz file:
- Intel: lockify_*_darwin_amd64.tar.gz
- Apple Silicon: lockify_*_darwin_arm64.tar.gz
```sh
# extract and install:
tar -xzf lockify_*_darwin_*.tar.gz
sudo mv lockify /usr/local/bin/
```

**Windows:**
#### Download the appropriate .zip file (amd64 or arm64)
```powershell
# Extract and add to PATH
Expand-Archive -Path lockify_*_windows_*.zip -DestinationPath .
# Move lockify.exe to a directory in your PATH
```

### Using Go

```sh
go install github.com/ahmed-abdelgawad92/lockify@latest
```

**Note:** After installing via Go, make sure `$GOPATH/bin` (or `$HOME/go/bin` if GOPATH is not set) is in your PATH:

```sh
# Add to ~/.zshrc or ~/.bashrc
export PATH="$PATH:$(go env GOPATH)/bin"
# Or if GOPATH is not set:
export PATH="$PATH:$HOME/go/bin"
```

### From Source

```sh
git clone https://github.com/ahmed-abdelgawad92/lockify.git
cd lockify
go build -o lockify .
```

After building, you may want to move the binary to a directory in your PATH:
```sh
sudo mv lockify /usr/local/bin/  # macOS/Linux
```

---

## Quick Start

### 1. Initialize a Vault

```sh
lockify init --env prod
```

### 2. Add a Key-Value entry

```sh
lockify add --env prod
```

### 3. Add a Secret entry

```sh
lockify add --env prod --secret
```

### 4. Export to `.env` (CI-friendly)

```sh
lockify export --env prod --format dotenv > .env
```

### 5. Import .env to a vault

```sh
lockify import .env --env prod --format dotenv
lockify import env.json --env staging --format json
```

### 6. Get a Value

```sh
lockify get --env prod --key DATABASE_URL
```

### 7. List all keys

```sh
lockify list --env prod
```

### 8. Delete an entry

```sh
lockify delete --env prod --key DATABASE_URL
```

### 9. Clear cached passphrase

```sh
lockify cache clear
```

---

## GitHub Actions Example

```yaml
steps:
  - uses: actions/checkout@v4

  - name: Install Lockify
    run: go install github.com/ahmed-abdelgawad92/lockify@latest

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
