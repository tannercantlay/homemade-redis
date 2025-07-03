# Homemade Redis

Building my own custom Redis server in Go.

## Running

Build and run with all base Go modules:

```sh
go run .
```

## Project Structure

- `main.go` — Entry point, TCP server logic
- `resp.go` — RESP protocol parsing

## TODO

- [x] Reading RESP
- [ ] Writing RESP
- [ ] Data Persistence
- [ ] Logging improvements