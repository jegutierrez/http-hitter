# http-hitter

## Usage

### Options
```bash
  -conc int
        How many instances to run in parallel. (default 4)
  -log
        If you like logging every request.
  -n int
        Number of requests
  -url string
        The URL you wish to hit. (default "http://localhost")
```
### Example

```bash
  go run request.go -url http://localhost:4000 -n 1000 -conc 4 -log true
```