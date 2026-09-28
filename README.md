# go-env

Simple, type-safe configuration management for Go applications based on environment variables.

`go-env` generates typed getter functions from a YAML configuration file, so you access
environment variables with compile-time–checked types, defaults, and validation instead of
string parsing scattered around your code.

## Features

- **Type-safe getters** — generates `Get<Var>()` / `Lookup<Var>()` functions with checked Go types.
- **Rich type support** — `int`, `float`, `duration`, `string`, `string_slice`, `bool`, `url`, `[]url`.
- **Validation rules** — critical and non-critical rules validated against the parsed value.
- **Aliases & fallbacks** — multiple env names per variable and an ordered lookup.
- **Markdown documentation** — optionally generates a Markdown reference for your variables.
- **Sections** — group variables with `# Section name.` comment headers.
- **Glob & multi-file** — accept one or more files, or a glob pattern.

## Requirements

- Go 1.25+

## Install

Install the generator CLI (`gen`):

```bash
go install github.com/kirill-zak/go-env/cmd/gen@latest
```

Then make sure your `PATH` includes `$(go env GOPATH)/bin`.

## Quick start

1. **Create a YAML configuration file** describing your environment variables:

   ```yaml
   env:
     check_batch_size:
       type: int
       default: 1000
       description: Check batch size
     check_interval:
       type: duration
       default: "10m"
       description: Check interval
   ```

2. **Add a `go:generate` directive** to a Go file in your package:

   ```go
   //go:generate go run github.com/kirill-zak/go-env/cmd/gen -p=example -o=config_gen.go -d ./docs/config.md config.yaml
   ```

   Run generation with:

   ```bash
   go generate ./...
   ```

   This produces a `config_gen.go` file with typed getters and constants, e.g.:

   ```go
   const CheckBatchSizeDefault = 1000

   // GetCheckBatchSize returns the value, falling back to the default on lookup failure.
   func GetCheckBatchSize() int { ... }
   ```

   See the complete example in [`example/`](./example).

3. **Use the generated getters** in your code:

   ```go
   cfg := example.GetCheckBatchSize()   // int
   d := example.GetCheckInterval()      // time.Duration
   ```

## CLI

```
envgen [flags] [file or glob pattern]
```

| Flag | Shorthand | Default      | Description                                   |
|------|-----------|--------------|-----------------------------------------------|
| `--package` | `-p` | `config` | Generated package name.                        |
| `--output`  | `-o` | `env_gen.go` | Path to write the generated getter file.       |
| `--doc`     | `-d` | *(empty)*     | Path to write generated documentation; empty disables it. |

The positional argument accepts one or more config files or a glob pattern.

## Configuration file format

The top-level `env:` mapping holds variables. Each variable supports the following fields:

| Field         | Type            | Description                                             |
|---------------|-----------------|---------------------------------------------------------|
| `type`        | `string`        | Value type: `int`, `float`, `duration`, `string`, `string_slice`, `bool`, `url`, `[]url`. |
| `default`     | `string`        | Default value used when the variable is not set.        |
| `description` | `string`        | Human-readable description (included in generated docs). |
| `validate.critical` | `bool`    | If `true`, validation failures become errors.           |
| `validate.rules`    | `[]string` | Validation rules applied to the value.                  |
| `alias`       | `[]string`      | Additional env names to look up (uppercased).           |
| `comment`     | `string`        | Free-form comment attached to the variable.             |

### Environment variable lookup

Each variable is looked up by its uppercased name plus any aliases, in order. If the global
environment prefix is configured (e.g. via `env.InitGlobal("MYAPP")`), the prefix is tried first,
then the bare name — as documented in the generated getters.

### Sections

Group variables into sections for the generated documentation by placing a `# Section name.`
comment header before a variable:

```yaml
env:
  # Startup.
  init_timeout:
    type: duration
    default: "5s"
```

## Runtime package

Generated code depends on the runtime package [`pkg/env`](./pkg/env), which provides the env
store, type conversion and validation. Import it in your module:

```go
import "github.com/kirill-zak/go-env/pkg/env"

func init() {
    env.InitGlobal("MYAPP") // optional prefix for all lookups
}
```

Key APIs: `GetGlobal`, `InitGlobal`, `NewEnvWithPrefix`, `Env.Lookup`, `Env.LookupWithErr`,
`AsGoType`, `URL`/`URLSlice` helpers, and the generic `Validate` helper.

## Development

Clone the repository and use the provided targets:

```bash
make lint   # run golangci-lint
make test   # run the test suite
```

## License

[BSD 2-Clause](./LICENSE)