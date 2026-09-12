# hypercube-cache-valkey

Valkey driver for [go-hypercube](https://github.com/go-hypercube/go-hypercube)'s `cache.Cache` interface, backed by [`valkey-io/valkey-go`](https://github.com/valkey-io/valkey-go).

## Install

```bash
go get github.com/go-hypercube/hypercube-cache-valkey
```

## Usage

```go
import (
	"github.com/valkey-io/valkey-go"
	valkeycache "github.com/go-hypercube/hypercube-cache-valkey"
)

client, err := valkey.NewClient(valkey.ClientOption{
	InitAddress: []string{"localhost:6379"},
})
if err != nil {
	panic(err)
}
cache := valkeycache.New(client)

app := hypercube.New(cfg, db, cache)
```

## Notes

- Implements the full `github.com/go-hypercube/go-hypercube/cache` interface: `Get`, `Set`, `Delete`, `Has`, `Increment`, `Expire`.
- Misses map to `cache.ErrNotFound` (checked via `valkey.IsValkeyNil`); negative TTLs return `cache.ErrInvalidTTL`.
- `Increment` uses Valkey's native `INCRBY`, which creates the key at 0 if absent — matches the interface contract exactly, no emulation needed.
- `Expire` with a zero TTL calls `PERSIST` rather than `EXPIRE 0`, since the `Cache` contract treats a zero TTL as "clear any existing expiration" rather than "expire immediately."
- Need raw Valkey data structures (hashes, lists, sets, sorted sets, pub/sub)? Bind the underlying `valkey.Client` into your app's container directly and `Resolve` it from plugins — this driver only covers the plain-cache subset.

## License

MIT
