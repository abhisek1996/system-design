## JWT
- secured way of transmitting information between client and server.
- it is digitaly signed JSON object, using RSA.

**where is it used**
- authentication
-  authorisation
- SSO


**JsessionID**
- used before
- session details are stored in DB
- it is stateful, as it hits DB
- in case of a distributed server it can cause a problem.


**JWT**
- it has all details
- it is stateless
-  **structure**
    - JWT/JWSJWS
    aaaaaaa.bbbbbbb.ccccccc
    - header.payload.signature
    - header: 
        - type: JWT
        - alg: RSA/HMAC
        - kid: key id
    - payload: (registered claims, public claims, private claims - not understood by multiple parties, only authentication sever can use it. )
        - iss: issuer
        - sub: subject
        - aud: audience
        - exp: expiration time
        - nbf: not before
        - iat: issued at
        - jti: jwt id
    - signature: 
        - encode header and payload seprately using base64.
        - concatinate both using . encodedHeader.encodedPayload
        - sign the encodedHeader.encodedPayload using RSA private key.
        - concatinate signature using . encodedHeader.encodedPayload.signature 


## challenges with JWT
- how to invalidate a token incase of a fraud or blacklisted users.
    - server needs to keep the list to blacklisted tokens (jti), can be caches for quick look up.
        - same problem for JsessionID
    - change the keys
        - all tokens are invalidated.
    - token should be short lived. **[this is popular]**
    - token should be used only once. **[this is popular]**

- JWT tokens are ecoded not encrypted.
    - anyone can decode the token and read the payload.
    - so we should not store any sensitive information in the payload.
    - JWE - JSON Web Encryption - this is used.
- uncesure JWT with "alg": none, such JWT should be rejected.
- JWK exploit: 
    - public key is shared in this
    {
        "jwk": {
            "kty": "RSA",
            "public: "key"
            "kid": "key id" **use this kid to get a whitelist public key from a listed of keys** **it attacker added this public key into the list that is a issue, so this listing of kid should be done properly**
        }
        
    }

    - never use the public key from the JWK . to verify the JWT token.


## The Full S2S Flow with RS256
```
┌─────────────────────────────────────────────────────┐
│                  KEY SETUP (one time)               │
│                                                     │
│  Order Service generates:                           │
│    privateKey (kept secret, never leaves service)   │
│    publicKey  (published to a JWKS endpoint)        │
│                                                     │
│  https://order-svc/.well-known/jwks.json            │
│  → { "keys": [{ "kid": "v2", "n": "...", "e": "AQAB" }] }  │
└─────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────┐
│              RUNTIME (every request)                │
│                                                     │
│  1. Order Service builds JWT payload                │
│  2. Signs with privateKey (RS256)                   │
│  3. Sends JWT in HTTP header:                       │
│     Authorization: Bearer eyJ...                    │
│                                                     │
│  4. Payment Service receives request                │
│  5. Decodes header → finds kid: "v2"                │
│  6. Fetches public key from JWKS (or cache)         │
│  7. Verifies signature → ✅                         │
│  8. Checks: exp not passed, aud = "payment-service" │
│  9. Grants access                                   │
|_____________________________________________________|

```

## The goal of signing is NOT to hide the token.
- The goal is to prove "only the private key owner could have produced this."
- So yes — everyone can "decrypt" (verify) the signature. That's intentional. It's a proof of origin, not a secret.
- Private Key  →  only ONE entity has this  →  signs
- Public Key   →  EVERYONE has this         →  verifies

"If this signature checks out with A's public key,
 then A must have signed it. No one else could have."

What's Actually Secure Here
The security question isn't "who can read it" — it's "who can FAKE it."