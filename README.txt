#SmartEdge Code Challenge Implementation

## Description
Given a string input of up to 250 characters, return a JSON response compliant to the schema defined below.
You are responsible for generating a public/private RSA or ECDSA keypair and persisting the keypair on the filesystem
Subsequent invocations of your application should read from the same files
Document your code, at a minimum defining parameter types and return values for any public methods
Include Unit Test(s) with instructions on how a Continuous Integration system can execute your test(s)
You may only use first order libraries, you may not use any third party libraries or packages.
For example, you may use the OpenSSL library, but you may not use any libraries built on top of OpenSSL.

JSON Schema for your application response:

```
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "title": "Signed Identifier",
  "description": "Schema for a signed identifier",
  "type": "object",
  "required": [ "message", "signature", "pubkey" ],
  "properties": {
      "message": {
          "type": "string",
          "description": "original string provided as the input to your app"
      },
      "signature": {
          "type": "string",
          "description": "RFC 4648 compliant Base64 encoded cryptographic signature of the input, calculated using the private key and the SHA256 digest of the input"
      },
      "pubkey": {
          "type": "string",
          "description": "Base64 encoded string (PEM format) of the public key generated from the private key used to create the digital signature"
      }
    }
}
```

Example:
```
>./your-awesome-app "your@email.com"
```
Returns:

{
  "message":"your@email.com",
  "signature":"MGUCMGrxqpS689zQEi5yoBElG41u6U7eKX7ZzaXmXr0C5HgNXlJbiiVQYUS0ZOBxsLU4UgIxAL9AAgkRBUQ7/3EKQag4MjRflAxbfpbGmxb6ar9d4bGZ8FDQkUe6cnCIRleaxFnu2A==",
  "pubkey":"-----BEGIN PUBLIC KEY-----\nMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEDUlT2XxqQAR3PBjeL2D8pQJdghFyBXWI\n/7RvD8Tsdv1YVFwqkJNEC3lNS4Gp7a19JfcrI/8fabLI+yPZBPZjtvuwRoauvGC6\nwdBrL2nzrZxZL4ZsUVNbWnG4SmqQ1f2k\n-----END PUBLIC KEY-----\n"
}
            


