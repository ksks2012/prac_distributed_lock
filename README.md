# Reference

<https://blog.csdn.net/jiandanokok/article/details/114296755#t4>

# Instruction

An implementation of Distributed Lock in golang

# structure

- etc: setting file
- docs: document
- global: global variables
- internal (internal module):
 <!-- TODO: -->
- dao: data access object
- middleware
- model: database model control
- routers: api routes
- service: process business logic
- pkg: package
- storage: temp file
- scripts: build, install, analysis scripts
- third_party: third_party tools

# Go generate

```sh
go generate github.com/distributed_lock/internal/dao/dbversion/mysql
```

# Build

```sh
go build github.com/distributed_lock/cmd/db_lock
```

# run

## run db_lock

```sh
./db_lock -config ./etc/config.yaml -log_level debug
```

### args

```sh
Usage of ./db_lock:
  -config string
        assgin the path of config file (default "etc/")
  -mode string
        running level (info, debug)
```

# Integration Testing

```sh
go build github.com/distributed_lock/testing/logger
```
