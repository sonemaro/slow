# SLOW

Slow sits in your infrastructure and detects your slow database queries.
It parses slow queries and provides a set of APIs over http.

# Build
### ***Source Code***
To build `slow` run this command:
```shell
go get -u github.com/sonemaro/slow
```
### ***Release***
Release

## Config
I tried my best to write `slow` based on `Twelve Factors` (https://12factor.net) so there isn't any configuration file. All you need is to set a few environment variables and there you go!

## Environment Variables
***SLOW_APP_ADDRESS***
Web API will be served on this address. Example: `0.0.0.0`

***SLOW_APP_PORT***
Web API will be served on this port. Example: `3000`

***SLOW_DB_TYPE***
This variable specifies type of the db in which we will try to parse slow query logs.
Available values: `mysql`, `postgres`

***SLOW_DB_LOGFILE***
Database log file path. Example:
`/var/lib/postgresql/12/main/pg_log/postgresql-2022-01-21_043743.log`

## Run

To run `slow`, set environment variables and run the binary.
Example:
```shell
SLOW_APP_ADDRESS=0.0.0.0 SLOW_APP_PORT=8083 SLOW_DB_TYPE=postgres SLOW_DB_LOGFILE=/var/lib/postgresql/12/main/pg_log/postgresql-2022-01-21_043743.log slow
```

**Note:** Sometimes reading from database log files requires sudo permission. So,
consider that you should run `slow` with `sudo` permission. 

## Docker

Again, `slow` has build as a `12 Factor App`! So dealing with docker is an easy thing.
Current docker image is extremely small and secure. It contains the binary and nothing else.
To build the image, run this command:
```shell
docker build -t latest .
```
For configuration, set environment variables. No config files!

***Docker Compose***
To run `slow` with compose, create `.env` file and run docker-compose
```shell
cp .env.example .env
docker-compose up
```

## API

Here it is the list of `v1` endpoints:

**Filter**  
[GET] /api/v1/filter/:operation  
`:operation` example: `select`, `delete`, `insert`, `update`

**Sort**
[GET] /api/v1/sort  
  
***Pagination***  
All endpoints support pagination by adding this query string to them:  
`?page=1&items=10`