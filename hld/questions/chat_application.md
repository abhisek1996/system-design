1. Rrquirement gathering
    1. functional (MVP)
        - 1-1 send and recieve - text, image, files
        - group message support.
        - message status - sent, delivered, read
        - last seen
        - user login - auth

    2. non functional (good to have)
        - scalability - huge traffic
        - availability - 99.99%
        - latency - low latency
        - fault tolerance


2. back of the envelope
    - total users - 2 billion
    - DAU - 50 milliom
    - 1 user 10 mssage to 4 people per day
    - 50 * 10 * 4 = 2 billion messages per day
    - 1 message = 100 charaters = 100 bytes
    - 2 * 10^9 * 100 = 2 * 10^11 bytes = 200 GB per day
    - 200 * 365 = 73 TB per year
    -  for 10 years = 730 TB

3. High level design
    -  it NOT a peer to peer - not scalable.
    - it must be done via a chat server.

    - HTTP :  it works on a request response basis.
        - client A ---> LB ---> chat server ---> DB
        - client B ---> LB ---> chat server ---> DB

        - sending works fine, but how will B receive the message?
            - solution 
                - polling
                    - client is asking the server for new messages over http.
                    - connection is created and closed.
                    - non scalable
                - long polling
                    - client is asking the server for new messages.
                    - connection is kept open for a threshold time. And it send NO if no message
                    - then connection is closed.
                    - non scalable
                - web socket
                    - bi-directinal peristent connection.
                    - connection is closed only when the user logs out or network issue :  **persistent**
                    - both client and sever and send and revieve: **bi-directional** 
 


    - 1-1 message
        - client A --->  chat server 100
        - client B ---->  chat sever 200

        - how to send a message from A to B?
            - user mapping service
                - which chat service maps to which user. 
                - zookeeper can we used here.
                - it also takes responsiblity of assigning a char server to a user based on the geography. record is added to zookeeper.
        - DB
            - read: 
                - search a user
                - search a group
                - search message history
            - write:
                -  send message
                - update profile

            SQL or No SQL
                - do we need heavy joins ?
                - we need heavy searching.

                