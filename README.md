# cofiswarm-kvpool

KV cache policy: auto-clear thresholds, proactive eviction, token budget (ported from coordinator).

- Migration: Sprint 6 in [MIGRATION-SPRINTS](https://github.com/keepdevops/cofiswarm-docs/blob/main/MIGRATION-SPRINTS.md)
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
| `POST /v1/admit` | `{"group":"llama8b","tokens":4000}` → KV token-budget gate |

Default listen: `:8014`.

## Arming the policy

Config is loaded at startup from `--config`, `$COFISWARM_KVPOOL_CONFIG`, or `configs/kvpool.yaml`.
A missing file leaves the policy **disabled** (safe default). The shipped `configs/kvpool.yaml` is
**armed** — `enabled: true`, proactive eviction at `0.60`, and per-`server_group` token budgets:

```yaml
enabled: true
proactive_threshold: 0.60
budgets:           # ≈ ctx_cap ÷ parallel slots; 0/absent = unbudgeted
  llama8b: 3072
  coder7b: 3072
```

`/v1/admit` enforces the budgets (the `kv_token_budget` gate); `/v1/evaluate` returns the
auto-clear / proactive-evict decision for a given KV pressure.

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
