## Why rate limiter ?
- DDOS attack makes server down.
- actual user are obstructed.

Algorithm
-
1. Token bucket
    - capacity = 4
    - more token ->  overflow
    - refiller -> refils token very 1 min
    - whena requets comes a token is used.

    e.g implememted through counter

    - rule, 3 req per user per minute
    - post - UI:{counter: 2, timestamp: 10.1.00 }
    - post - UI:{counter: 1, timestamp: 10.1.20 }
    - post - UI:{counter: 0, timestamp: 10.1.35 }
    - post - UI:{counter: 0, timestamp: 10.1.35 } - dropped ---> 429


2. leaking bucket
    - fixed capacity
    - implememted through queue

    - request is processed in a constant rate.
    - dist: when a constant burst comes this won't work.
3. fixed window counter
    - divide into a fixed time interval and it has a counter for each section.
    - edge scenario issue - request can be double.
4. sliding window log
    -  the widow slides.
    - 3 req/min
        - window is 1 minute
        - post - UI:{size: 1, timestamp: 10.1.00 }
        - post - UI:{size: 2, timestamp: 10.1.20 }
        - post - UI:{size: 3, timestamp: 10.1.35 }
        - post - UI:{size: 3, timestamp: 10.1.40 } -- log is stored - but not proceed.
        - 
    - disv: lot of logs are stored.
5. sliding window counter
    - take adv of 3 & 4
    - 