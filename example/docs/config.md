# Environment Variables

## general

| Env variable | Type | Constraints | Default value | Description |
| --- | --- | --- | --- | --- |
| `CHECK_BATCH_SIZE` | `int` | — | `1000` | Check batch size |
| `CHECK_INTERVAL` | `time.Duration` | — | `10m0s` | Check interval |
| `CHECK_SCHEDULE` | `string` | — | `"* * 2-6 * * * *"` | Check schedule |
