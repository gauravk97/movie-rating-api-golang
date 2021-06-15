- Start postgres
- Prepare environment, fill DB parameters:

``` bash
$ source env-sample
```

- Build and run:

```bash
$ export GO111MODULE=on
$ export GOFLAGS=-mod=vendor
$ go mod download
$ go build
$ ./movie-rating.exe
```

Server is listening on localhost:8010