## CAP Theorem

- desirable property of dist system with replicated data


## Basics

- Consistency
    - data should be same for all the nodes

- Availability 
    - all the noded should respond

- Partition Tolerence
    - user should be able to get a resp even if there is a partition.

-  all 3 cannot be used together. Why ?
    - A - 
    - B - a = 5
    - C - a = 5

        - if there is a partion b/w B & C for 5 minutes
        - A - a = 6
        - B - a = 6
        - C - a = 5 - replication was not possible because of partition.

## scenarios
 - A & P is possible -  we dropped C
 - P & C - we will make C node down, A is compromised.
 - C & A - if there is a partion then the entire system has to be down.


 ## In real world
  - P is always required
  - CP or AP







