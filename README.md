# Lockify 🔐 

A developer-first CLI for managing encrypted environment variables across environments (dev, staging, prod). Safe to store encrypted in Git. CI/CD friendly. Built in Go.

---

## TL;DR

- **Encrypted environment variables** per environment (dev, staging, prod)
- **Encrypted vaults can be committed to Git** — passphrases kept separate
- **Optional OS keyring caching** for local development convenience
- **CI/CD export support** via environment variable (no interactive prompts)
- **Not a replacement for enterprise secret managers** like Vault, AWS Secrets Manager, or SOPS
- Optimized for developer ergonomics and small-team workflows

---

## Table of Contents

- [Why Lockify?](#why-lockify)
- [Key Features](#key-features)
- [Installation](#installation)
  - [Using Homebrew](#using-homebrew)
  - [Download Pre-built Binaries](#download-pre-built-binaries)
  - [Using Go](#using-go)
  - [From Source](#from-source)
- [Quick Start](#quick-start)
- [Vault Location & Configuration](#vault-location--configuration)
- [Passphrase Caching](#passphrase-caching)
- [Why Lockify vs. Alternatives?](#why-lockify-vs-alternatives)
- [What Lockify Is NOT](#what-lockify-is-not)
- [GitHub Actions Example](#github-actions-example)
- [Security Summary](#security-summary)
- [Contributing](#contributing)
- [License](#license)
- [Security Contact](#security-contact)

---

## Why Lockify?

Lockify solves a common problem for developers and small teams: how to keep environment variables secure, sharable, and scriptable without relying on third-party cloud vendors or complex infrastructure.

It focuses on an offline-first, CLI-first experience with strong cryptography and a minimal mental model.

- Store encrypted vaults per environment
- Export decrypted environment variables in CI/CD using a secret passphrase
- Designed for developers and small teams who prefer Git-native workflows
- Keep secrets versioned alongside code without exposing plaintext

---

## Key Features

- **AES-256-GCM Encryption** (authenticated encryption)
- **Argon2id KDF** for deriving encryption keys from passphrases
- **Passphrase caching** via OS keyring (optional, for local convenience)
- **Multi-environment vaults** (dev, staging, prod, or custom environments)
- **Import/export** `.env` and JSON formats
- **Key rotation** without losing data
- **Offline-first** — no network dependencies
- Clean, testable codebase using Domain-Driven Design and clean architecture principles

---

## Installation

### Using Homebrew

```sh
brew tap ahmed-abdelgawad92/lockify
brew install lockify
```

### Download Pre-built Binaries

Visit the [GitHub Releases](https://github.com/ahmed-abdelgawad92/lockify/releases/latest) page to download pre-built binaries for your platform.

**Linux:**

Download the appropriate `.tar.gz` file for your architecture (amd64 or arm64):

```sh
# Extract and install
tar -xzf lockify_1.2.3_linux_amd64.tar.gz
sudo mv lockify /usr/local/bin/
```

**macOS:**

Download the appropriate `.tar.gz` file:
- Intel: `lockify_*_darwin_amd64.tar.gz`
- Apple Silicon: `lockify_*_darwin_arm64.tar.gz`

```sh
# Extract and install
tar -xzf lockify_*_darwin_*.tar.gz
sudo mv lockify /usr/local/bin/
```

**Windows:**

Download the appropriate `.zip` file (amd64 or arm64):

```powershell
# Extract and add to PATH
Expand-Archive -Path lockify_*_windows_*.zip -DestinationPath .
# Move lockify.exe to a directory in your PATH
```

### Using Go

```sh
go install github.com/ahmed-abdelgawad92/lockify@latest
```

**Note:** After installing via Go, ensure `$GOPATH/bin` (or `$HOME/go/bin` if GOPATH is not set) is in your PATH:

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

After building, move the binary to a directory in your PATH:

```sh
sudo mv lockify /usr/local/bin/  # macOS/Linux
```

---

## Quick Start

### 1. Initialize a Vault

```sh
lockify init --env prod
```

**Example Output:**

```
Initializing Lockify vault
? Enter passphrase for environment "prod": ********
? Confirm passphrase: ********
? Cache passphrase in system keyring? No
✓ Lockify vault initialized at .lockify/prod.vault.enc
```

Use the `--cache` flag to skip the interactive caching prompt:

```sh
lockify init --env prod --cache
```

### 2. Add a Key-Value Entry

```sh
lockify add --env prod
```

**Example Output:**

```
Setting a new entry to the vault...
? Enter passphrase for environment "prod": ********
? Enter key: DATABASE_URL
? Enter value: postgresql://localhost:5432/mydb
✓ key DATABASE_URL is added successfully.
```

### 3. Add a Secret Entry

For sensitive values that should be masked during input:

```sh
lockify add --env prod --secret
```

**Example Output:**

```
Setting a new entry to the vault...
? Enter passphrase for environment "prod": ********
? Enter key: API_SECRET
? Enter secret: ********
✓ key API_SECRET is added successfully.
```

### 4. List All Keys

```sh
lockify list --env prod
```

**Example Output:**

```
Listing all secrets in the vault
✓ Found 3 key(s):
  - DATABASE_URL
  - API_SECRET
  - SMTP_PASSWORD
```

### 5. Get a Value

```sh
lockify get --env prod --key DATABASE_URL
```

**Example Output:**

```
getting an entry from the vault
✓ retrieved key's value successfully
postgresql://localhost:5432/mydb
```

### 6. Export to `.env` (CI/CD-friendly)

```sh
lockify export --env prod --format dotenv > .env
```

**Example Output:**

```
DATABASE_URL=postgresql://localhost:5432/mydb
API_SECRET=sk_live_1234567890
SMTP_PASSWORD=mypassword123
```

### 7. Import from `.env` File

```sh
lockify import .env --env prod --format dotenv
```

**Example Output:**

```
Importing env variables to the vault...
? Enter passphrase for environment "prod": ********
✓ Imported 5 key(s), skipped 0 key(s)
```

### 8. Delete an Entry

```sh
lockify delete --env prod --key DATABASE_URL
```

**Example Output:**

```
deleting entry from the vault
✓ key "DATABASE_URL" deleted successfully
```

### 9. Cache Management

Cache a passphrase explicitly:

```sh
lockify cache set --env prod
```

Clear cached passphrase for a specific environment:

```sh
lockify cache clear --env prod
```

Clear all cached passphrases:

```sh
lockify cache clear
```

---

## Vault Location & Configuration

### Where Are Vaults Stored?

By default, Lockify creates a `.lockify/` directory in your **current working directory** and stores encrypted vault files there:

```
your-project/
├── .lockify/
│   ├── prod.vault.enc
│   ├── staging.vault.enc
│   └── dev.vault.enc
├── .git/
└── src/
```

**Important:** Vault files are **fully encrypted** and safe to commit to Git. The `.lockify/` directory structure makes it easy to:
- Track encrypted secrets alongside your code
- Share vaults with team members via Git
- Deploy with your application

### Path Behavior

Vaults are relative to where you run the command:

```sh
cd /path/to/project
lockify init --env prod  # Creates /path/to/project/.lockify/prod.vault.enc
```

**Best Practice:** Run Lockify commands from your project root to keep vaults consistent.

### Tradeoff Note

This project-relative approach optimizes for simplicity and Git integration. The tradeoff is that vaults are not centralized across projects. Support for configurable or global vault paths may be added in future versions based on community feedback.

---

## Passphrase Caching

### How Caching Works

Lockify can optionally cache passphrases in your **OS keyring** (macOS Keychain, Windows Credential Manager, Linux Secret Service) to avoid repeated prompts during day-to-day development.

### Caching Methods

**1. Interactive Prompt (Default)**

```sh
lockify init --env prod
? Enter passphrase for environment "prod": ********
? Confirm passphrase: ********
? Cache passphrase in system keyring? (y/N)
```

**2. Explicit Flag (Non-Interactive)**

```sh
lockify init --env prod --cache
```

This automatically caches without prompting—ideal for scripts.

**3. Manual Caching**

```sh
lockify cache set --env prod
? Enter passphrase for environment "prod": ********
✓ Passphrase cached successfully for environment "prod"
```

### When Are Passphrases Required?

- **First Use:** When creating or accessing a vault for the first time
- **After Cache Clear:** After running `lockify cache clear`
- **Cached:** If passphrase is in keyring, no prompt (validates silently)
- **CI/CD:** Use `LOCKIFY_PASSPHRASE` environment variable (no caching needed)

### Security Considerations

- Caching uses OS-level encryption (Keychain/Credential Manager)
- Passphrases are **never** stored in plaintext
- Cache is per-environment (each vault has its own cached passphrase)
- Clear cache when changing devices or for added security

---

## Why Lockify vs. Alternatives?

### vs. OpenSSL / Manual Encryption

**OpenSSL:**

```sh
# Encrypt
openssl enc -aes-256-gcm -salt -in .env -out .env.enc

# Decrypt (no key management, no structure)
openssl enc -d -aes-256-gcm -in .env.enc -out .env
```

**Lockify:**

```sh
# Multi-environment support, key-value structure
lockify init --env prod
lockify add --env prod
lockify export --env prod --format dotenv > .env
```

**Advantages:**
- **Multi-Environment Management:** Separate vaults for dev/staging/prod
- **Structured Access:** Get individual keys without decrypting entire file
- **Git-Friendly:** Encrypted JSON format with metadata
- **Key Rotation:** Change passphrases without losing data
- **CI/CD Integration:** Export directly to `.env` format
- **Passphrase Caching:** Optional OS keyring integration

### vs. Cloud Secret Managers (AWS Secrets Manager, Vault, etc.)

Lockify is ideal when you want:
- **Offline-First:** No cloud dependencies
- **Git-Native:** Secrets versioned with code
- **Zero Cost:** No monthly fees
- **Simple Setup:** No infrastructure to maintain
- **Developer Workflows:** Works without internet

**When to use cloud managers:**
- Large-scale production systems
- Secrets rotation at scale
- Compliance requirements (audit logs, access control)
- Cross-service secret sharing
- High-risk or regulated production environments

### vs. `.env` Files in Git (Plaintext)

❌ **Never commit plaintext secrets to Git**

Lockify lets you:
- ✅ Commit encrypted vaults safely
- ✅ Share secrets with team via Git
- ✅ Version control your secrets
- ✅ No risk of accidental exposure

---

## What Lockify Is NOT

Lockify is designed for developer workflows and small teams. It is **not** intended as a replacement for enterprise-grade secret management systems.

**Lockify is NOT:**

- **A replacement for HashiCorp Vault, AWS Secrets Manager, or SOPS** — These tools offer features like dynamic secrets, fine-grained access control, audit logs, and integration with identity providers that Lockify does not provide.

- **Intended for high-risk or regulated production environments** — If you're operating in healthcare (HIPAA), finance (PCI-DSS), or other regulated industries, you likely need enterprise solutions with compliance certifications and audit trails.

- **A zero-trust secret delivery system** — Lockify relies on passphrase-based encryption and does not integrate with identity providers, certificate authorities, or policy engines.

- **Optimized for large-scale secret rotation** — While Lockify supports key rotation, it does not automate rotation schedules or integrate with external systems for dynamic credential generation.

**Lockify IS optimized for:**
- Developer ergonomics and day-to-day workflows
- Small teams that prefer Git-native secret management
- Projects where offline-first and zero-cost are priorities
- Environments where the risk model aligns with passphrase-based encryption

If you need enterprise features, consider using Lockify for developer machines and switching to a managed secret service for production.

---

## GitHub Actions Example

Store your passphrase as a GitHub Actions secret (`LOCKIFY_PASSPHRASE`), then export environment variables in your workflow:

```yaml
steps:
  - uses: actions/checkout@v4

  - name: Install Lockify
    run: go install github.com/ahmed-abdelgawad92/lockify@latest

  - name: Export environment variables
    env:
      LOCKIFY_PASSPHRASE: ${{ secrets.LOCKIFY_PASSPHRASE }}
    run: lockify export --env prod --format dotenv > .env

  - name: Use secrets in subsequent steps
    run: |
      source .env
      echo "Database URL: $DATABASE_URL"
```

**Security Note:** Ensure your GitHub repository secrets are properly configured and access is restricted to authorized workflows.

---

## Security Summary

Lockify is intended for developer workflows and small teams, not as a replacement for enterprise-grade secret management systems.

**Encryption & Key Derivation:**
- **AES-256-GCM** (authenticated encryption with associated data)
- **Argon2id** key derivation function (memory-hard, resistant to GPU attacks)
- Per-environment vaults with independent encryption keys

**Passphrase Handling:**
- Passphrases are **never** stored in plaintext
- Optional caching uses OS-level encrypted keyrings (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- In CI/CD, passphrases are provided via environment variables

**Vault Files:**
- Vault files are **fully encrypted ciphertext** and safe to commit to Git
- Non-sensitive metadata (environment name, creation time) is stored in plaintext
- Encrypted data includes all key-value pairs

**Key Rotation:**

Rotate passphrases without losing data:

```sh
lockify key rotate --env prod
```

This re-encrypts the vault with a new passphrase derived from a new key.

**Best Practices:**
- Use strong, unique passphrases for each environment
- Store production passphrases in CI/CD secret managers (GitHub Secrets, GitLab CI/CD variables, etc.)
- Clear cached passphrases when changing devices or handing off machines
- Regularly rotate passphrases for sensitive environments

---

## Contributing

Contributions are welcome! Please see `CONTRIBUTING.md` for guidelines on how to contribute to Lockify.

---

## License

MIT — see `LICENSE`.

---

## Security Contact

If you discover a security vulnerability, please report it responsibly:

**Email:** ahmed.elkayaty92@gmail.com

Please do not open public issues for security vulnerabilities.
