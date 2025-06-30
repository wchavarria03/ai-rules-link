# Usage & CLI Options

## Rules Lookup Order

When you run any command, the CLI looks for rule files in this order:
1. `${XDG_CONFIG_HOME}/ai-rules` (if set and exists)
2. `~/ai-rules` (in your home directory)
3. If neither exists, the CLI uses its own embedded rules (no setup needed)

You do **not** need to copy rule files manually. If you want to override or customize rules, create one of the above directories and add your own rule files.

## Basic Usage

```bash
ai-rules-link rules --rule=go --rule=python
```
- Symlinks selected rules into your project's `.cursor/rules/` directory.

## Using the --global Flag

By default, rules are created in the folder where you run the CLI (the current working directory). If you want to create rules in your home directory instead, use the `--global` flag:

```bash
ai-rules-link rules --rule=go --rule=python --global
```
- This will create symlinks or consolidated rules in `~/` (your home directory) instead of the current directory.

You can combine `--global` with other flags, such as `--consolidate`:

```bash
ai-rules-link rules --rule=go --rule=python --consolidate --global
```
- This will create a single `consolidatedrules.mdc` file in `~/.cursor/rules/`.

## Symlink Rules for Cursor or Consolidate into One File

You can symlink any set of rules into your project's `.cursor/rules/` directory using the `rules` command and the `--rule` flag:

```bash
ai-rules-link rules --rule=go --rule=base --rule=nextjs --rule=python --rule=personalcommits --rule=workcommits
```

- This will create symlinks in `.cursor/rules/` for each specified rule (e.g., `gorules.mdc`, `baserules.mdc`, `nextjsrules.mdc`, `pythonrules.mdc`, `personalcommitsrules.mdc`, `workcommitsrules.mdc`).
- The canonical rules files must exist in `~/.sync-rules/rules/` (e.g., `~/.sync-rules/rules/gorules.mdc`).
- If a canonical file does not exist, the command will print an error and skip that rule.
- The command is fully dynamic: you can use any rule name as long as the corresponding file exists.
- **To create rules in your home directory instead, add `--global` to the command.**

### Consolidate All Rules into One File

You can merge all selected rules into a single file using the `--consolidate` flag:

```bash
ai-rules-link rules --rule=go --rule=python --rule=personalcommits --consolidate
```
- This will create a single file named `consolidatedrules.mdc` in your destination rules directory (default: `.cursor/rules/`).
- The file will contain the merged content of all selected rules, in the order specified.
- No symlinks will be created when using `--consolidate`.

## Example

```bash
ai-rules-link rules --rule=go --rule=nextjs --rule=python --rule=personalcommits --rule=workcommits
```
- This will attempt to symlink `gorules.mdc`, `nextjsrules.mdc`, `pythonrules.mdc`, `personalcommitsrules.mdc`, and `workcommitsrules.mdc` from `~/.sync-rules/rules/` into your project's `.cursor/rules/` directory.

## Listing Symlinks

To see which rules are currently symlinked in your project:

```bash
ai-rules-link status
```
- This will list all symlinks in `.cursor/rules/` and their targets.

## --force Flag

If you use the `--force` flag, the CLI will always overwrite destination files with embedded rules, even if those files have been modified by the user. Use this with caution if you want to reset rules to the embedded defaults. 