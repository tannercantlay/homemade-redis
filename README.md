# Homemade Redis

Building my own custom Redis server in Go.

## Running

Build and run with all base Go modules:

```sh
go run .
```

## Project Structure

- `main.go` — Entry point, TCP server logic
- `resp.go` — RESP protocol parsing and responding

## TODO

- [x] Reading RESP
- [x] Writing RESP
- [ ] Reading RESP commands
- [ ] Data persistence
- [ ] Logging improvements
- [ ] Containerize
