## Indexing
1.**How data is actually stored ?**

- DBMS creates data pages and manages. Size ~ 8KB = 8 * 
- Data page:    
    - Header: 96 bytes
        - page no
        - free space
        - checksum
    - Data records
        - 8060 bytes
        -  if each row is of 64 KB
        - ~ 125 DB rows 
    - Offset Array - 36 bytes
        - it is an array, it points to each row in the data records

- **Data BLock**:
    - minimum amount of data which an i/o operation can fetch.
    - it is a physical memory.
    - managed by storage system, not by DBMS.
    - **pages are stored in inside these data blocks**
    - **Size : 4KB - 32 KB**

- DBMS:
    - maintains a mapping of data page and data block.

2. **What type of indexing present in RDBMS**
    1. Indexing
        - clustered indexing
            - order of the row inside the data pages match with the order of the indexing.
            - this is dones with use offset.
            - actual order of the entiries in the page may varry.
            - the offtset array is sorted as per the idex and it points to each row.

            - there can be one clustered index per table.
                - priority: primary key
            - if we have not primary key DBMS create an hidden auto increment column and make it as clustered index.

        - non clustered
            - used in case of non primary key or a conposite key.
            - there be many composite index.
            - so adding many secondary index is not a good idea, hude memory space.

    2. Indexing: 
        - increase perp of DB query.
        - wihtout it dbms has to go through the all **data pages and each rows in them.**
            - load the page to main memory.
            - read each row.
    
    3. B+ tree: 
        - TC: O(logn)
        - B Tree: Balanced Tree
            - maintained sorted data.
            - m order b tree: each node can have atmost m children. 
            - m pointers and (m-1) key.
        - B+ tree: all child rows are linked to each other.
            -  child nodes are sorted from left to right.
            


    search a value:
    - index pages are loaded.
    - b+ tree is loaded.
    - data pages is loaded.
    - data block.
    - read the data.
        
