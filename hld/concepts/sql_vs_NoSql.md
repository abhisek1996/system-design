
## SQL
- used to query RDMS.
- tables, rows, colums.
- relation between multiple tables.


**structure** 
- table, row, column
- pre dtermined schema
- relations between tables

**Nature**:
- concentrated / centralised. 
    - all data for a server in same server.
    - half of table in this and other is other server [this is not well supported]

scalability
- vertical scalability
- horizontal - not well supported.

property 
    - ACID


## NO SQL


**structure** 
- unstructured data
    - key value
        - e.g dynamo DB
        - key - value (it can be anything)
        - value is opaque:  we cannot query or search based on value.
        - very fast.
    - document
        - key : value[json, XML]
        - value is not opaque. we can query and search based on value.
        - e.g mongo DB
    - columnwise
        - key: value[column1:value1, column2:value2, column3:value3, ...]
        - each key list of  [ column: value ]
        - number of columns can be dynamic
    - graph db
        - nodes and edges
        - relation on the edges
        - social media. 

**Nature**:
- distributed
- all data for a server in same server.

**scalability**
- horizontal scalability

**property** 
- BASE
    - B,A: Basically available
        - highly available
    - S: safe state
        - state is data can be changes without user interaction.
    - E: Eventually consistent
        - initially we may stale data
        - later on we can get latest data, once sync is done between the nodes.



SQL vs NoSQL

| Feature | SQL | NoSQL |
| --- | --- | --- |
| Query | complex, dynamic | basic, or know columns|
| Data | relational nature | not tightly coupled |
| Data integrity | strong, we cannot loose a single transaction, financial data, consistency | weak, deals with volume,  |
| Availablity | low | high, with some inconsistency, searching capability is very fast |
| **Best For** | Structured data, complex queries | Unstructured/semi-structured data, high scalability |

