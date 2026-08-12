## Why retry ?
- it save a lot of resources, there are a tons of things which are done before a downstream call.


## When retry ?
- don't retry on permanent errors.
    - validation errors
    - 4XX
    - 429 - too many request -  retry after a gap.
- don't if API is not idemptent.

- 5xx -  retriable errors


## When retry ?

- fixed interval retry.
    - adv:
        - easy
        - easy to config
    - dis:
        - thurndering heard effect.
            - all will retry at same time.
            -  new requests will starve.
- expotential backoff.
    - delay increases exponentially.
    - delay = base*2^no.failed retry attempt.
    - adv:
        - reduce load
        - lesser thurndering heard problem.
    - disadv:
        - still chances for thundering heard problem.
        - if outage period is long user may have to wait for a longger time.

- exponential + jitter.
    - jitter adds randomness to avoid storms. 
    - delay = random(0, min(maxDelay, eponential backoff))







