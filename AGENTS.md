# AGENTS.md

Guidance for AI coding agents (and human contributors) working in this repository.

## Project overview

`go-env` generates type-safe environment-variable configuration code for Go applications.
It reads a YAML file describing env variables and generates a Go file with typed
`Get<Var>()` / `Lookup<Var>()` functions plus `*Default` constants.

The project has two public parts:

1. **Generator CLI** — `cmd/gen` (binary name `envgen`), typically invoked through `go generate`.
2. **Runtime package** — `pkg/env`, which the generated code imports.

Module: `github.com/kirill-zak/go-env` · Requires Go 1.25+.

## Repository layout

| Path | Purpose |
|------|---------|
| `cmd/gen/` | Cobra CLI: flag parsing and DI wiring (thin layer, no business logic). |
| `pkg/env/` | Public runtime: `Env` store, `Lookup` / `LookupWithErr`, `AsGoType`, `Validate`, URL helpers. |
| `pkg/generator/` | Public generator orchestration: `Generator`, `FileReader`, `FileWriter`, `TemplateEngine`. |
| `entity/` | Shared domain model used by all layers (`Variable`, `Section`, `GeneratorArgs`, template args). No external dependencies. |
| `internal/variables/` | YAML parsing into `entity.Variable` / `entity.Section`. |
| `internal/env/converter/` | String → Go type conversion (shared by the generator and the runtime). |
| `internal/generator/templating/` | `text/template` rendering + import fixing (`golang.org/x/tools/imports`). |
| `errors/` | Typed errors: `ErrEnvNotExists`, `ErrInvalidValue`, `ErrInvalidType`, `ErrInvalidRule`. |
| `example/` | Working example: `config.yaml` → `gen.go` (go:generate directive) → `config_gen.go` + `docs/config.md`. |
| `scripts/go.mk` | Make targets included by the root `Makefile`. |

## Development

```bash
make test   # go test ./... -race -cover
make lint   # golangci-lint v2 (config: .golangci.yml)
make fmt    # gofmt + goimports
make vet    # go vet ./...
```

Regenerate the example:

```bash
cd example && go generate ./...
```

CI runs golangci-lint v2 and the test suite on every pull request to `master`.

## Generating env getters from a YAML file

1. **Install the generator CLI:**

   ```bash
   go install github.com/kirill-zak/go-env/cmd/gen@latest
   ```

   Make sure `$(go env GOPATH)/bin` is on your `PATH`.

2. **Create a YAML configuration file** describing your variables (format below), e.g. `config.yaml`:

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

3. **Add a Go file to the target package** (convention in this repo: `gen.go`, see `example/`)
   with a `go:generate` directive:

   ```go
   package myapp

   //go:generate go run github.com/kirill-zak/go-env/cmd/gen@latest -p=myapp -o=config_gen.go -d ./docs/config.md config.yaml
   ```

4. **Run generation:**

   ```bash
   go generate ./...
   ```

   This writes the generated file (`config_gen.go` with the flags above; default is `env_gen.go`).
   For every variable it emits: a `<Name>Default` constant, `Get<Name>()` (returns the default
   when the lookup fails) and `Lookup<Name>()` (returns `(value, error)`).
   With `-d`, a Markdown reference is written as well.

   Do not edit generated files by hand — change the YAML and re-run `go generate`.

### CLI reference

```
envgen [flags] [file or glob pattern]
```

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--package` | `-p` | `config` | Package name of the generated file. |
| `--output` | `-o` | `env_gen.go` | Path of the generated Go file. |
| `--doc` | `-d` | *(empty)* | Path of the generated Markdown reference; empty disables it. |

The single positional argument is a config file path or a glob pattern; variables from all
matching files are merged, duplicates removed by variable name.

## Using the project (runtime)

Generated code is an ordinary Go package — import it and call the getters:

```go
n    := myapp.GetCheckBatchSize()      // int, falls back to CheckBatchSizeDefault
d, err := myapp.LookupCheckInterval()  // (time.Duration, error)
```

Runtime knobs, package `github.com/kirill-zak/go-env/pkg/env`:

- **Global prefix** — call once before the first lookup (e.g. in `init()`):

  ```go
  import "github.com/kirill-zak/go-env/pkg/env"

  func init() { env.InitGlobal("MYAPP") }
  ```

  Lookups then try `MYAPP_<VAR>` before `<VAR>`. `InitGlobal` is a no-op after the first
  call (`sync.Once`).
- **Custom env instances** — `env.NewEnvWithPrefix(prefix)` plus `(*Env).Lookup` /
  `LookupWithErr` when you do not want to use the global instance.
- **Middleware** — `(*Env).AddMiddleware(fn ...EnvMiddlewareFn)`; consulted after all env
  names, in registration order (useful for config files, secrets managers, tests).
- **Validation** — generated `Lookup` functions apply `validate.rules`; with
  `validate.critical: true` a violation panics, otherwise an error is returned.

**Lookup order** (`Env.addPrefixToNames`): all prefixed env names first, then all bare names.
For prefix `MYAPP` and names `[A, B]`: `MYAPP_A`, `MYAPP_B`, `A`, `B`.

## Using go-env in a third-party project

Standard setup (generator tool + runtime dependency):

```bash
# 1. Install the generator tool
go install github.com/kirill-zak/go-env/cmd/gen@latest

# 2. Add the runtime dependency (the generated code imports pkg/env)
cd your/module
go get github.com/kirill-zak/go-env
```

Then add a `config.yaml` and a `//go:generate` directive as described above, run
`go generate ./...`, and **commit the generated file**.

Alternative — use the generator as a library (advanced, avoids a separate CLI install):

```go
import (
    "github.com/kirill-zak/go-env/entity"
    "github.com/kirill-zak/go-env/pkg/generator"
)

logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
g := generator.NewGenerator(
    generator.NewFileReader(logger),
    generator.NewTemplateEngine(),
    generator.NewFileWriter(logger),
    logger,
)

err := g.Generate(entity.GeneratorArgs{
    PkgName: "myapp",
    OutName: "config_gen.go",
    DocName: "docs/config.md", // optional
    Files:   []string{"config*.yaml"},
})
```

Note: the `configReader` / `templateEngine` / `fileWriter` interfaces in
`pkg/generator/contract.go` are unexported, so the supported wiring is the concrete
implementations shown above (mirrors `cmd/gen/main.go`).

## YAML configuration format

The top level is a single `env:` mapping; each key is the variable name, whose uppercased
form is the primary env name.

| Field | Type | Description |
|-------|------|-------------|
| `type` | `string` | Required. One of `int`, `float`, `duration`, `string`, `string_slice`, `bool`, `url`, `[]url` (`strings` is an alias of `string_slice`). |
| `default` | `string` | Value used when the variable is unset; must parse to the declared type. |
| `description` | `string` | Human-readable description (rendered into the docs). |
| `validate.critical` | `bool` | `true` → validation failure panics; otherwise an error is returned. |
| `validate.rules` | `[]string` | `positive`, `negative` (applies to `int` / `float` / `duration`). |
| `alias` | `[]string` | Additional env names (uppercased), appended after the primary name in lookup order. |
| `comment` | `string` | Free-form note, rendered into the Markdown docs. |

### Type mapping

| YAML `type` | Go type | Parsing |
|-------------|---------|---------|
| `int` | `int` | `strconv.Atoi` |
| `float` | `float64` | `strconv.ParseFloat(_, 64)` |
| `duration` | `time.Duration` | `time.ParseDuration` |
| `string` (or empty) | `string` | as-is |
| `string_slice` / `strings` | `[]string` | comma-separated, trimmed |
| `bool` | `bool` | `false` / `0` → `false`, any other non-empty value → `true` |
| `url` | `*url.URL` | `url.Parse` |
| `[]url` | `[]*url.URL` | comma-separated `url.Parse` |

### Full example

```yaml
env:
  # Startup.
  init_timeout:
    type: duration
    default: "5s"
    description: How long startup may take
    validate:
      critical: true
      rules:
        - positive

  http_port:
    type: int
    default: "8080"
    description: HTTP listen port
    alias:
      - PORT

  upstreams:
    type: []url
    default: "http://localhost:8080"
    description: Upstream URLs
```

### Sections

A `# Section name.` comment directly above a variable marks a section in the generated
Markdown docs. Note: the current pipeline (`FileReader.ReadConfig` →
`variables.ReadFromYAML`) does not propagate section comments, so all variables land in a
single `general` section (simplified logic in `pkg/generator/generator.go`). Full section
parsing exists in `internal/variables.ReadSectionFromYAML` and is covered by tests.

## Code style and conventions

- **Dependency direction** (top → bottom): `cmd` → `pkg/generator` → `internal/*` + `entity`.
  `entity` is dependency-free. Generated code depends only on `pkg/env`.
- **Public API** is `pkg/env`, `pkg/generator`, `entity` and `errors` only. `internal/*` is
  private and must not be imported from outside (or by third-party projects).
- **Errors**: wrap with `fmt.Errorf("...: %w", err)` and use the typed errors from `errors/`
  (they implement `Unwrap`, so `errors.Is` / `errors.As` work).
- **Tests**: `testify`, table-driven, `_test.go` files next to the source. Keep
  `make test` green (`-race -cover`).
- **Logging**: `log/slog` via an explicit `*slog.Logger` dependency.
- **Naming**: snake_case YAML names become PascalCase identifiers via `snakeToCamel`
  (acronym table: DSN, HTTP, API, …); getters are `Get<Name>` / `Lookup<Name>`,
  defaults are `<Name>Default`.
- **Formatting**: keep `make fmt` (gofmt + goimports) clean.
