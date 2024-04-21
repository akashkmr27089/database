# Creating a basic database where crud operation can be done (SQL)
1. Creating a database.
2. Inserting a data to the meomory RAM
3. Accessing the data --> List all, get one, and where conditions 
4. Able to update one or many data 
5. Able to detete the database 


# Design.
1.  First Create Common command database where list of database is stored,
    1.a Each database will contain schema_meta, index_meta files --> These files will contain the structure for table.
2.  If yes, create an entry if the database exists. Each database is a single file 
3.  Writing to common command database is immidiately and remaining databse entry is batchified --> First comple incremental
    database entry to the database 
4.  Incremental: 
    a. If incremental, all the changes are directly written in the database 
    b. The files are encoded before sending it back to client,

    Batchify: 
    # TODO
5. All the operation will go to engine as an oplog code --> database, collection, command


# Task 
Create Basic CRUD operation (Done)
Create Batch and incremental flush Basic Structure (Done)
Crate Machanism for writing data to file --> Secodary Storage (Done)
Load and upload data from secondary storage to main memory and vice versa (Done)
Make serach more capable --> implement trie algorithm for searching a key 


1. Create a basic structure to store data in structure and store it as a file
2. A parser to understand the file and reach to a specific part of the file 
3. For starter work on key value pair. This will be synced directly

# Open question: 
1.  How Databases knows which part of file to access when knowing the row number. Do they store the pointer seperately ? 
    --> We know for each column how much data to take to store the data --> 
        int32 --> 32bit, float32 --> float32 etc 
        use this as hint for knowing the pointer position 
2.  Should the database be row based or column based ? --> If column based, then how the datas will be stored in files. internally, do we need to store the mapping ?
    --> Does column based database needs more amount of space then the row based database ? 
    --> Row based database. 
    --> How does the sharding works in cloumn based database (Question for lateron)
    --> Seek generally takes constant amout of time for ssd --> Thus seeking to specific part of file and then taking data 

