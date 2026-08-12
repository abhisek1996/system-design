## caching
- reduce latency
- achive fault tolerence


different types of cache:
- client side [browser cache]
- CDN -  to store static data
- Load balancer
- Server side acaching. - Redis


**Server side application caching**
    client ---> LB ---> app server ---> cache -----> DB

## What is disributed cahing
- common cache for all app sever 
    - scalability
    - single point of server

- app sever ---->. cache client 
                        --------> cache server 1
                        --------> cache server 1
                        --------> cache server 1
                        --------> cache server 1
- app sever is alloted a cache server based on CONSISTENT HASHING technique


## caching strategies

1. **cache aside**
    - C --> AS ---> C
            if cache hit --> AS
            if cache miss --> fetch from DB ---> write to cache ---> AS
    **pros**
    - good for READ HEAVY
    **cons**
    - not updating whilte write
    -  data in consitency 
        - we are not invalidating the cache during write
                DB : A = 10
                cache : = 10
                write ---> DB : A = 11
                cache := 10 (not updated)

2. **Read throguh cache**
    - in case of cache miss
    - the cache take the responsibilty fetching the data, caching it back.

    **pros**
    - good for read heavy

    **cons**
    - same as above
    - cache structure has to be same as DB

3. **write around cache**
    - write directly to DB
    - invalidate cache
        -  dirty flag = true
    - it has to to be used with any one of the previous 
    
     **pros**
    - good for read heavy

    **cons**
    - same as above
    - write operation is totaly dependent of DB
        - write is not fault tolerant

4. **write through cache**
    - write to cache
    - write to DB
    - if any one fails,we have tp fail the it.


     **pros**
    - consistent
    - increase cache hit

    **cons**
    - alone it is not useful
    - 2 phase commit needs to be implemented.
    - not fault tolerant

5. **write back (or behind) cache**
    - write ti cache
    - Push data to queue
        - DB is updated from queue
    - even if DB is down still it works


    **pros**
    - good for write heavy application.
    - write latency is down.
    - cache hit chance is high
    - better if used with erite through cache aside/read through cache

    **cons**
    - ttl is 3 hrs
    - db is down for 5 hrs
        - queue won;t be able to update DB
    - after 3.5 hrs data is not in cache not in DB