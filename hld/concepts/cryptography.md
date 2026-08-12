## Encryption
    - convert a readble a text to a chiper.
    - cyrptographic key.


1. **Symmetric Encryption**
    - fast :  shorter key
    - same key used for encrpt and decrpt.
    -  algorithm: 
        - DES (not recomended) - 56 bits
        - AES (recomended)
            - 126, 192, bits
    Advantage:
        - fast
        - good for bulk data
            - used in chat applications.
    disadvantage:
        - key distribution is a problem between client and server
        - key should be different for all the clients, server has to manage those number of keys and distribution.

2. **Asymmetric Encryption**
    - slow
    - different key is used for encrpt and decrpt.
        - private.
        - public key.
    -  algorithm:
        -  RSA
        -  DSA
        -  diffie-hellman
        -  ECC
    - adavantage
        - key distribution is not a problem.
        - good for authenti