# Queue Simulation using GRPC and Golang

### Start the server
``` bash
go run cmd/server/main.go
```
### Start the producer
``` bash
go run cmd/producer/main.go
```

### Start the consumer
``` bash
go run cmd/consumer/main.go
```

## SQLITE

### Queue SQLite 3

Create a database named test.db using the following command:

``` bash
sqlite3 db/queue.db
```

The command outputs the following:

``` bash
SQLite version 3.16.2 2017-01-06 16:32:41
Enter ".help" for usage hints.
sqlite>
```

### Check database list using the following command:

``` bash
sqlite> .databases
```

### Type .exit to exit the CLI.

``` bash
sqlite> .exit
```

### Run command ls and you will find the file database is created.

``` bash
$ ls /db
queue.db
```

## Table

### Create messages table
``` sql
CREATE TABLE messages (
   ID TEXT,
   Body TEXT,
   PRIMARY KEY(ID)
);

```

### Create queues message
``` sql
CREATE TABLE queues (
   ID TEXT,
   Name TEXT,
   PRIMARY KEY(ID)
);
```