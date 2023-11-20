# pandafs

Distributed file system in Golang but it's useless for all intents and purposes
and only exists for my personal entertainment.

## How it works

The project is split into 4 parts:

1. Client: The client-side CLI that communicates to the master.
2. Master: Acts as a name node; handles metadata and manages the mules.
3. Mule: Acts as a worker node; where the data is actually stored.
4. Core: Contains compiled protobuf files, and node-handler functions.

## How to run it

    # Get all the module dependencies
    go get ./..

    # Run the master node
    go run main.go master

In another terminal(s),

    # Run the mule
    go run main.go mule


## How to send commands

~~Client is still a WIP lmao~~

##### This project is maintained by [Abhinav Chennubhotla](https://github.com/PhoenixFlame101).
