Goals
- scalability
- de-centralisation
- eventual consistency 

steps
1. partition
    - consistent hashing

2. de-centralisation
    - replication
        - co-ordinator(server where data is stored) will put date in N-1 servers
        - these is some logic to chose those N-1 servers, moving in clockwise direction in consistent hashing circle.
            - different data centers.
        - N is some value
        - A preference List is maintained by the co-ordinator.
            - for range [1- 45]
                - s1 - co-ordinator
                - s2
                - s3
            - for range [46 - 80]
                - s2 - co-ordinator
                - s3
                - s4
            - for range [81 - 100]
                - s3 - co-ordinator
                - s4
                - s5
            - every node knows about the preference list of each range.

3. Get and put operation 
    - PUT
        - R + W > N
            - if W =1, after putting in co-ordinator is 1 server is replicated we send success response to client. 
    - Load balancer
        - generic
            - can go to any server, but they will re-direct it based on the preference list.
            - high latency
        - partition aware
    - R = 2
        - for a GET operation, wait for atleast 2 replicas to respond.
        - then we considers success and send to client.
4. Data versioning
    - in get of GET, why S1 doesn't directly respond to client ?
        -  No system is perfect, there are failures.
        - S1 - 45->car
        - then a update operation comes with 45-cart
        - S1 is down, so it doesn't get updated.
        - S2 - 45->cart, and then it went down and couldn't replicate in s3.
        - s3 - 45 -> carm
        - S1->car, S2->cart, S3->carm :  different versions of same key.
    - **vector clock**
        - list of [server, counter] for each key. 
        
    - dynamo DB is eventually consistent
        - A, P. C is sacrifised. 

5. GOssip protocol
    - each server maintains a membders list --> last seen, range they cater.
    - each sever sends heart beat to a few servers.
    -  if more than 1 servers say that s1 heart beat is not received then it is down.
6. Merkle Tree
    - how to check if the data is same in all the replicas ?
