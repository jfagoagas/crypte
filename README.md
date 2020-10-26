# Crypte
Simple tool for (de)compress, (verify)sign and (de)crypt a message using lz4 and NaCl encryption.

## Usage
```
## Crypte ##
Tool for (de)compress and (de)crypt message using NaCl anf lz4.

Usage:
- Generate Public/Private Keys:
crypte -k

- Encrypt, sign and compress with lz4 a message
crypte -e -p <PublicKeyFile> -s <PrivateKeyFile> -m <File>

- Decrypt, verify sign and decompress with lz4  a message:
crypte -d -p <PublicKeyFile> -s <PrivateKeyFile> -m <File:w>
```

