# VictoriaMetrics

## Command-line Flags

### `-keepOldData`
If set to `true`, this flag prevents the automatic deletion of old data when the `-retentionPeriod` is reduced. This is useful to avoid accidental data loss when switching to a smaller retention period.

If set, old data won't be removed on startup when a smaller retention period is provided.

Example:
```bash
./victoria-metrics -retentionPeriod=30d -keepOldData=true
```

- Default: `false`
- Type: `boolean`

`