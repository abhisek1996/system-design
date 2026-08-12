## Question

- Design a tiny url kind of service, take a long url and return a short url.


## requirement analysis

- How short -> as much as we can
- traffic :
    - 10 M users per day  => 2650 M url/year
    -  it should last for 100 years => 365 billion URLs
- characters:
    - a - z | 0 - 9 | A - Z
    - 62 characters
    - _, _, _ ... 6 characters, each can be choosen in 62 ways.
    - pow(62, 7) = 3.5 trillion
    -  so we need atleast 7 characters.

-  how to generate the hash value:
    - hash function
        - MD5 => 128 bits -- 16 bytes 
        - SHA1
    - base64
