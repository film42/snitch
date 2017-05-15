Snitch
======

An aggregate health checking service. Point snitch at http servers or processes and it will return a `200` status code
if more than 70% of the checks are passing. Snitch was built to be used with AWS Cloud Watch for monitoring docker
containers by triggering if a snitch process falls out of a private ELB.

### Installing

To build from source you can do:

```
$ go get github.com/film42/turbulence
$ ./bin/turbulence
```

Or you can grab the pre-built docker container:

```
$ docker run -p 9999:9999 -d film42/snitch:latest --check localhost:3000 --check localhost:3001
```

## Docs

```
Must provide at least one process or host to check:
  -check value
        List of host:port combos to check. Example: --check localhost:1234 --check localhost:3422
  -port int
        Port to listen on (default 9999)
  -process value
        List of process substring to check. Example: --process sidekiq --process puma
```
