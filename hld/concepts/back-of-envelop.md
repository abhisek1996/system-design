## Why BOE ?
drives our descision for system design

Consderation:
1. Rough or t-shit size estimation.
2. don't spent more time.
3. keep the assumption value simple
    - take whole number -  100 million, 10 million etc for easy computation.


## cheat sheet

3 zero
thousand
kb

6 zero
million
mb

9 zero
billion
gb

12 zero
trillion
tb

15 zero
quatrallion
pb


*3 basics things*
- no of servers
- RAM
- storage 

**then go for trade off =>  CAP**


X million * Y MB => xy TB
    6, 0s * 6, 0s --> TB

5 million users * 2 K =>  5*2    6 + 3 => GB ==> 10GB




## Estimation of facebook

total user - 1 billion
DAU - 25% of total - 250 million

every user
    -  5 read + 2 writes = 7 query

(250 million * 7) / 60*60*24 ===> 18K querries for second


**storage**
- every user 
    - 2 posts - 250 char - 1 char 2 byte ==> 500 *2  byte ==> 1KB
        - 250 million * 1 KB ===> 250 GB
    - 10 % user upload 1 image - 3 KB
        - 25 million users
        - 300KB * 25 million ==> 7500 GB ==> ~ 8TB

- for 5 years - 2000 days
    - 2000 * 250 GB ~ 5 TB
    - 2000 * 8 TB ~ 20 PB


**ram estimation**

for each DAU: 
- last 5 post for each user
- 1 post --> 500 bytes
- 5 * 500 --->   ~ 3KB
- total ---> 3KB *250 million ===> 750GB

- if one machine 75 GB, so we need 10 machines
- latency - P95 => 500ms


**server estimation**

- 18K / sec
- 1 server => 50 thread
- 500 ms each request for P95 --- so one sever can serve 100 RPS


-  total server
    - 180 severs




