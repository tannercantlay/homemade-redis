# Homemade Redis

Building my own custom Redis server in Go.

## Running

Build and run with all base Go modules:

```sh
go run .
```

Build and run docker container:
'''sh
docker build -t homemade-redis .
docker run -p 6379:6379 homemade-redis
'''

## Project Structure

- `main.go` — Entry point, TCP server logic
- `resp.go` — RESP protocol parsing and responding

## TODO

- [x] Reading RESP
- [x] Writing RESP
- [x] Reading RESP commands
- [x] Data persistence
- [ ] Del function
- [ ] DBSize function https://redis.io/docs/latest/commands/dbsize/
- [ ] Logging improvements
- [x] Containerize

