# go-env
Simple configuration management of a Go application based on environment variables

## Install
Install `go-env` by command
```bash
go install github.com/kirill-zak/go-env/gen/gen@latest
```

## How to use

1. Install package.
2. Create file with env.
3. Run cmd:
```bash
//go:generate go run gen -p=example -o=config_gen.go -d ./docs/config.md config.yaml
```

## Example
See `./example/config.yaml` with example env.