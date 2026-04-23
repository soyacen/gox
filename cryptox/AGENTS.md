# cryptox AGENTS.md

## OVERVIEW

Cryptography utilities cluster: symmetric/asymmetric encryption, hashing, HMAC, TLS, and X.509 certificate helpers.

## STRUCTURE

```
cryptox/
├── x509.go          # Certificate PEM encoding/decoding
├── aesx/            # AES symmetric encryption (ECB/CBC/CFB/OFB/CTR)
├── hmacx/           # HMAC signatures (MD5, SHA1, SHA224, SHA256, SHA384, SHA512, SHA512_224, SHA512_256)
├── knuthx/          # Knuth hash with architecture detection
├── md4x/            # MD4 hash wrappers
├── md5x/            # MD5 hash wrappers
├── rsax/            # RSA key gen, sign/verify, encrypt/decrypt
├── shax/            # SHA family (1, 224, 256, 384, 512, 512_224, 512_256)
└── tlsx/            # TLS server/client config builders
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add AES mode | `aesx/aes.go` | Implements `Cipher` interface; add to switch |
| Add hash algorithm | `shax/` | Mirror pattern: `Hash` + `HashHex` |
| Add HMAC variant | `hmacx/` | Follows sha naming convention |
| RSA key format | `rsax/` | Supports Hex and Base64 encodings |
| TLS quick config | `tlsx/` | Auto-random cert fallback for testing |
| Certificate PEM | `cryptox/x509.go` | Top-level helper, not in subpackage |

## CONVENTIONS

- **Dual API**: Every hash/HMAC provides raw bytes + Hex string variants (`Sha256` / `Sha256Hex`)
- **Naming**: Function names match algorithm exactly (`Sha256`, `HmacSha256`, `Md5Hex`)
- **Cipher interface**: `aesx.Cipher` abstracts 5 block modes with `Encrypt`/`Decrypt` methods
- **RSA formats**: Key generation returns both Hex and Base64 string representations
- **Unsafe arch detection**: `knuthx` uses `unsafe.Sizeof(uintptr(0))` for 32/64-bit multiplier selection
- **Cross-dependency**: `rsax/sign.go` imports `cryptox/shax` for hashing during sign/verify

## ANTI-PATTERNS

- Do NOT add streaming/chunked hash APIs; all functions accept `[]byte` and return complete digest
- Do NOT use `md4x` or `md5x` for security-sensitive operations; they exist for legacy compatibility
- Do NOT ignore RSA key size parameter; default is 2048, minimum should be enforced
- Avoid duplicating hash logic between `shax` and `hmacx`; HMAC wraps hash via `hmac.New`
- Do NOT modify `knuthx` multiplier constants; they are tied to Knuth's multiplicative hash
