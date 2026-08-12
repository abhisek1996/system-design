
1. Single server
2. Application and DB server separation
    - separate both servers. cleint -> appl server --> DB
3.  LB + multiple app servers
    - LB  and app server - user pricate IP - it is secure
    - client doesn't call the app server directly but calls the LB
4. Database replication
    - master --- slave
    - write on master and read on slave
    - if there is a failure i have replcias
    - if master fails, slave will become master
5. Use of cache
    - app server to DB --->  use a cache
    - TTL
6. CDN
    - does caching of static data across the region.
    - is data is not present it will ask the nearest CDN if not then to DB server.
    - cleint --  CDN - even before load balancer
7. Data centers
8. Messaging queues
    -  for async opertation
    - exchange
        - direct
            - one to one mapping to subscriber
        - fanout
            - sends to all queue - subscriber read or ignored
        - topic
            - can send to > 1 queue : by using some kind of wild card matching.
9. Database scaling
    **a. Vertical**
        - increase the capability
        - there is always a limit
    **b. horizontal**
    - add more nodes
        - sharding
            - horizontal - this is considered better
                - divide table by rows
                - or based on some logic, say by name first letter
            - vertical
                - divide by column
        - draw back 
            - if all names comes with A, 1st shard will be filled quickly ad we have to shard it further, it will form a tree - consistent hashing
            - if rows are divided how to use join - de normalise