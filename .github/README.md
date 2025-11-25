# GitHub Actions CI

This directory contains continuous integration workflows for the learning-redis project.

## Workflows

### `ci.yml` - Build and Test

**Triggers:**
- Push to `main` or `master` branch
- Pull requests to `main` or `master` branch

**What it does:**
1. ✅ Builds all Go modules and examples
2. ✅ Verifies all code compiles without errors
3. ✅ Runs tests (if any)
4. ✅ Checks Go code formatting

**Components tested:**
- Main module (`go.mod`)
- Basic examples (strings, lists, sets, hashes, sorted sets)
- Interview scenarios:
  - Caching (cache-aside pattern)
  - Leaderboard (sorted sets)
  - Rate limiter (token bucket)
- Mini-Redis simulator

## Local Testing

To test what CI does locally:

```bash
# Build everything
make test  # or go test ./...

# Check formatting
gofmt -s -l .

# Build each component manually
go build ./...
cd examples/basic/strings && go build
cd examples/interview-scenarios/01-caching && go build
# ... etc
```

## Viewing CI Results

When you push code or open a PR, GitHub Actions will automatically:
- Run all builds
- Report results in the PR
- Show green checkmark ✅ if all builds pass
- Show red X ❌ if any build fails

Check the "Actions" tab in GitHub to see detailed logs.

