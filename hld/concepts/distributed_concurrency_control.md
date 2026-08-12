## scenarios
- many concurrent are tring to book the same seat in movie theater.
- critical section :  trying to access a shared resource.

LLD solution: 
    - using syncronized for critical section.
    - same process and many thread threads - this will be handled.

but what will this work in a distrbuted system ?
    - distributed concurrency control


## distributed concurrency control
- Omtimistic con control
- pessimistic con control


1. what is the usage of transaction ?
    - it helps to achive integrity.
    - rollback if all operation is not done.
    - DB remains consistent.

2. DB locking 
    - no other transaction update the locked row.
    - 2 locks
        - shared locks
            - read locks
            - only read can happen
            - multiple transaction can take lock but onyl for read.
        - exclusive locks
            - if one transaction takes a lock.
            - one can read it or write to it.
            - no one can take shared or exclusive locks

3. Isolation level present
    - this tells us how much level of concurency is allowed.

    1. Dirty read: 
        - T1 reads a data written by T2 which is not yet commited.
        - if T2 rolls back, all data in T1 is dirty.
    
    2. Non repeatable read
        - T1 reads diffetent data in same transaction.
        - if same row is read multiple time and there is a chance of getting different data because of some other T2 have committed.

    3. Phantom Read:
        - if a transaction wxecutes same query bit gets different rows.
        




    1. read uncomiited
        read lock - yes
        write lock - yes

        dirty - YES
        NRP - YES
        PR - YES

        - good for only read scenario 

    2. read committed
        - read - shared lock on read and release as soon as read is done. [updation is possible as soon as read is done]
        - write - exclusive only after transaction is done.

        dirty - NO
        NRP - YES
        PR - YES

    3. repeatable read

        - read - shared lock on read and release only after transaction is done.  [update is not pssible as write lock by t2 is not allowed ]
        - write - exclusive  and release only after transaction is done.

        dirty - NO
        NRP - NO
        PR - YES

    4. serializable
        - same as RR + apply a range lock and release at end of transaction.
        - nothing can be added inside row, it locks nearby ranges as well.

        dirty - NO
        NRP - NO
        PR - NO

1. OCC
    - read C

- solved concurency using version
- validate verson before updating

- much higher concurency
- 


2. PCC -- 2PL is the way to implement it.
    - RR
    - SR

    - this can lead to dead locks
    - have to rollback and try again

    t1
        - read A
        - write B
    t2
        - read b
        - write a

