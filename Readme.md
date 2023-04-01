# Rewards Microservice

A simple service for managing user points which allows receiving point based on the activity type i.e. add review, 
invite friend and these points can be spent or shared with other users as well.

## Used packages
- mysql
- gorm
- google.golang.org/grpc
- github.com/kelseyhightower/envconfig


## Directory structure

- bin: contains executable file(s)
- cmd: server and client main functions
- docker: contains all docker files
- internal:
  - api: implement grpc server methods and include environment struct 
  - models: database models (user, activity)
  - pb: generated code from proto using the following command:
    ```shell
    make generate
    ```
- pkg: emtpy reserved for reusable packages
- proto: contains protobuf files

## Run the service
Build and run docker containers using this command:
```shell
make up
```

## Stop the service
Stop and remove docker containers using this command:
```shell
make down
```

## Test using client application

To run the test client with default parameters:
```shell
make client
```
it should print something like that:
```
2023/04/01 12:31:02 balance: 219.500000
2023/04/01 12:31:02 successfully shared points: 5.500000
2023/04/01 12:31:02 user points: 214.000000
2023/04/01 12:31:02 recipient points: 66.000000
```

## How to authenticate users (not implemented)
We can use JWT token to authenticate user by adding interceptor function
to grpc and get token metadata. This way we can get current user's data.

