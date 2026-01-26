# Git Hooks Configuration

## Overview

GRGN Stack uses Husky to enforce code quality through automated git hooks.

## Installed Hooks

### pre-commit

Runs before every commit:

- **lint-staged**: Automatically fixes linting issues on staged files only
  - TypeScript/React: ESLint with auto-fix
  - Go: gofmt formatting + go vet checks
- **Quick tests**: Runs unit tests (Go and TypeScript)

**Skip if needed** (not recommended):

```bash
git commit --no-verify -m "your message"
```

### commit-msg

Validates commit message format:

- **Enforces Conventional Commits** format
- Format: `type(scope): subject`
- Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build, revert

**Examples:**

```bash
✅ feat(auth): add Google OAuth support
✅ fix(api): handle null pointer in user service
✅ docs: update configuration guide
✅ test(backend): add coverage for handlers
❌ updated some files  # Invalid!
```

### pre-push

Runs before pushing to remote:

- **Full test suite**: Runs all Go and TypeScript tests
- **Blocks push if tests fail**

**Skip if needed** (emergency only):

```bash
git push --no-verify
```

## Configuration Files

### package.json

```json
{
  "lint-staged": {
    "web/**/*.{ts,tsx}": ["cd web && eslint --fix"],
    "backend/**/*.go": ["gofmt -w", "cd backend && go vet ./..."]
  }
}
```

### .husky/

- `pre-commit` - Linting and quick tests
- `commit-msg` - Message format validation
- `pre-push` - Full test suite

## Bypassing Hooks

### Individual Commit

```bash
git commit --no-verify -m "emergency fix"
```

### Individual Push

```bash
git push --no-verify
```

### Disable Temporarily

```bash
# Disable hooks
git config core.hooksPath /dev/null

# Re-enable
git config --unset core.hooksPath
```

## Best Practices

1. **Don't bypass hooks** - They catch issues early
2. **Write meaningful commits** - Follow Conventional Commits
3. **Fix issues locally** - Don't push broken code
4. **Small commits** - Hooks run faster
5. **Trust the automation** - Hooks maintain quality

## Troubleshooting

### Hooks not running

```bash
# Reinstall hooks
npm run prepare
```

### Lint-staged fails

```bash
# Check what files are staged
git status

# Fix manually
npm run lint:fix

# Try again
git commit
```

### Tests fail on commit

```bash
# Run tests locally
npm run test

# Fix issues
# Commit again
```

## CI/CD Integration

These hooks mirror CI/CD checks:

- Pre-commit = CI linting stage
- Pre-push = CI test stage
- Ensures code is CI-ready before pushing

## Updating Hooks

Edit files in `.husky/` directory:

- `.husky/pre-commit` - Pre-commit checks
- `.husky/commit-msg` - Message validation
- `.husky/pre-push` - Pre-push checks

Changes take effect immediately for all developers after pulling.
