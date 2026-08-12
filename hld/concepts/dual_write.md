# Dual write problem
 - when a compnent needs to persist a change in 2 different systems.
    - DB & queue
    - maestro & LSS queue

- happens in SAGA pattern.

**Why not 2 phase commit**
-  heterogenous system
    - a message broker doesn't have this.
- 2PC is slow, as it provides strong consustency.
- If 2PC co-ordinator may crash.



# Event driven architecture

## transactional outbox
- in a single transaction DB chnange is inserted in the DB and the event is alos inserted in a outbox table.
- poller server polls on the outbox table and publishes the unpublished event.
- dis adv - 
    - 

## Listen to yourself pattern
- same as above.
- event is put in outbox DB.
- then both consumer and the emitter DB listens from same event.

-  all dis advantage as above
- + one more: if GET/ before write to the emiiter DB
-  use cache

## Transactional log tailing pattern

- CDC from the DB logs.

- dis:
    - latent, as DB may process logs in batch



**common dis:**
- duplicate event can bepicked by poller
     -  enable.idemptency = true
- ordering of events in KAFKA
    - kafka streams in a solution.
    - one publisher and one partion in kafaka can solve this. [not practical]
    -  idempotency needs to be handled at consumers end.
    
-  outbox table can grow
    - delete publised row
    - batch clean up job needed to be triggered.
    
- poller not able to publish
    - retry
        -  have some status
            - success, pending, failed etc.

