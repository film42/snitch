Snitch
======

An aggregate health checking service. Point snitch at http servers or processes and it will return a `200` status code
if more than 70% of the checks are passing. Snitch was built to be used with AWS Cloud Watch for monitoring docker
containers by triggering if a snitch process falls out of a private ELB.

## Running

```
$ go get github.com/film42/snitch
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
