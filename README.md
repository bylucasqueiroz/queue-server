# Queue Simulation using GRPC and Golang

This project implements a simplified, in-memory message queue system using gRPC in Go, inspired by AWS Simple Queue Service (SQS). The goal of the project is to create a basic task queue system where producers can send messages to a queue, consumers can retrieve those messages with a visibility timeout, and once a task is complete, the messages can be deleted from the queue.

### Generate GRPC code by Proto

``` bash
protoc --go_out=. --go-grpc_out=. --proto_path=. queue.proto
```

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