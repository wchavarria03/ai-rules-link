# Canonical Rules Location

The CLI looks for rules in this order:
1. `${XDG_CONFIG_HOME}/ai-rules` (if set and exists)
2. `~/ai-rules` (in your home directory)
3. If neither exists, the CLI uses its own embedded rules (no setup needed)

You only need to create these directories if you want to override or customize the default rules.

For example, to override:
- Create `~/ai-rules` and add files like `gorules.mdc`, `pythonrules.mdc`, etc.

All rules must exist as markdown files in:

```
~/.sync-rules/rules/
```

For example:
- `gorules.mdc`
- `pythonrules.mdc`
- `baserules.mdc`
- etc. 