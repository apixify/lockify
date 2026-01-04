# Security Policy

Lockify prioritizes security above all else.  
This document explains how to report vulnerabilities and how we approach security.

---

## Supported Versions

The latest stable version is supported with security updates.

| Version | Supported |
|--------|-----------|
| 1.x.x  | ✔ Yes     |

---

## Reporting a Vulnerability

If you find a vulnerability, **do not open a public GitHub issue.**

Instead, contact:

```
ahmed.elkayaty92@gmail.com
```

We will respond within 48 hours.

---

## Security Architecture Overview

- Vault files use **AES-256-GCM** authenticated encryption.
- Encryption keys derived using **Argon2id**, protecting against brute-force attacks.
- Passphrases **never stored** in plaintext.
- Optional passphrase caching uses **OS keyring** secure storage.
- Vault format includes versioning for safe future migrations.
- CI systems must provide passphrases explicitly via secret manager.

---

## Recommendations for Users

- Use strong passphrases (10+ random words).
- Use a different passphrase per environment.
- Rotate keys periodically:

```sh
lockify key rotate --env prod
```

- Never commit exported `.env` files.
- Only commit encrypted vaults.
