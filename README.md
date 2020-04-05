# Crypte
Simple tool for (de)compress, (verify)sign and (de)crypt a message using NaCl encryption.

## Usage
```
## Crypte ##
Tool for (de)compress and (de)crypt message


Usage:
- Generate Public/Private Keys:
./crypte -k
- Encrypt, sign and compress a message:
./crypte -e -p <PublicKeyFile> -s <PrivateKeyFile> -m <Message>
- Decrypt, verify sign and decompress a message:
./crypte -d -p <PublicKeyFile> -s <PrivateKeyFile> -m <Message>
```

