# cofiswarm-kvpool

KV cache policy: auto-clear thresholds, proactive eviction, token budget (ported from coordinator).

- Migration: Sprint 6 in [MIGRATION-SPRINTS](https://github.com/keepdevops/cofiswarmdev/blob/main/docs/MIGRATION-SPRINTS.md)
- Legacy C++ reference: `legacy/cpp/` (`kv_auto_clear.h`, `token_ledger`, etc.)

## Policy defaults (from `kv_auto_clear.h`)

| Field | Default |
|-------|---------|
| `proactive_threshold` | 0.60 |
| `pressure_threshold` | 0.75 |
| `proactive_fraction` | 0.30 |
| `divergence_threshold` | 0.6 |

## HTTP

| Route | Description |
|-------|-------------|
| `GET /healthz` | Liveness |
| `GET /v1/policy` | Current policy config |
| `POST /v1/evaluate` | `{"kv_pressure":0.65}` → clear/evict decision |

Default listen: `:8014`.

## Build & run

```bash
make build
./bin/cofiswarm-kvpool
```

## FHS

| Path | Purpose |
|------|---------|
| `/etc/cofiswarm/kvpool/kvpool.yaml` | policy |
| `/var/lib/cofiswarm/kvpool/` | ledger / cache state |
