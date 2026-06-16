# cofiswarm-kvpool

Cofiswarm component: `kvpool`.

- Layout: [REPO-STANDARD-LAYOUT](https://github.com/keepdevops/cofiswarmdev/blob/main/docs/REPO-STANDARD-LAYOUT.md)
- Migration: [MIGRATION-SPRINTS](https://github.com/keepdevops/cofiswarmdev/blob/main/docs/MIGRATION-SPRINTS.md)

## FHS paths

| Path | Purpose |
|------|---------|
| `/etc/cofiswarm/kvpool/` | config |
| `/var/lib/cofiswarm/kvpool/` | state |
| `/var/log/cofiswarm/kvpool/` | logs |

## Test

```bash
./test/scripts/assert-layout.sh kvpool
```
