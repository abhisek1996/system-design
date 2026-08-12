## Microservies and Patterns
- Monolythic ( legacy service )
- microservices

**Monolythic** - everything in same place
- Overload IDE
- scaling is very hard
    - CI is fast
    - if I want to scale even one functionality of the service, entire server needs to be scaled.

**Microservice**
- All dis of monolyth is an adv of micro
- scale is easy 
- each component can be scaled separatedlty.

- Disadvantage
    - Breaking into small services is critical.
    - Not loosely coupled - a lot of API calls.
    - increase in latency
    - monitoring all the services is difficult.
    - debugging is difficult -  change is resp of once can break others.
    - transaction management is difficult.
        - each server has its own DB, we cannot have once common transaction.


## Patterns

**Phases**
- decomposition
    - bussiness capability - BC
        - order management
        - payment
    - sub domain -  DDD - domain driven design
- database
    - common DB
    - DB per service
- communication
    - API
    - Events
- integration ( with client / UI)
    - API gateways 
- deployment
- observability
....


## strangler
- when we refactor monolyth to microservice
- slowly we move few servies to microservice
- add a controler to divert traffic from monolyth to microservice
- slowly we degrade the monolyth

## data management in microservice**
- DB for each sevrice 
    - X -  manage transaction is diffcult we have to have once transaction. - solved by **saga**
    - X - joining accross the table in different DB - solved by **CQRS**
- shared DB
    -  s1, s2, s3
    - X - if s1 needs more space, i have scale the entire DB
    - X - if s3 wants to modify the table it has to check if it impacts s1 or s2
    - transaction is easy


## saga
- orchestrator
    -  s1, s2 , s3 publish events to a success and failure queue, if there is failure all rollback.
-  choriography 
    - ochestrator calls S1, s2, s3, if one fails rollbacks all

**question**
    - P1 makes a 100 payment 
        - deduct balance -- passes
        - record -- fails
        - use SAGA - publish a failure event and send a compensation event to revert the balance

## CQRS

- command query request segregation
- command - CRUD
- Query for joins.
- each DB do their CRUD
- have a common view for query
- how the common view is updated 
    - via event 
    - via DB trigger



