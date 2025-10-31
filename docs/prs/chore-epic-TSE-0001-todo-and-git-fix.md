# chore(epic-TSE-0001): update protobuf schemas submodule

## Summary

This PR updates the `proto` submodule to the latest version, incorporating recent schema changes from the protobuf-schemas repository.

**Submodule Updates:**
- Updated from `b9193ea` to `ae535ab`
- Includes PR #15: Vendor sales proto schemas (epic VIE-0002 phase 1)
- Includes PR #14: Git quality standards foundation (epic TSE-0001)

## What Changed

### Proto Submodule Update

**New Commits Included:**
1. Merge pull request #15 - Vendor sales proto schemas for VIE-0002 phase 1
2. Merge pull request #14 - Git quality standards foundation for TSE-0001

**Impact:**
- Market data adapter now has access to latest vendor sales schema definitions
- Aligned with ecosystem-wide git quality standards
- No breaking changes to existing market data service contracts

### Repository Context

The protobuf-schemas repository is a submodule shared across all services in the trading ecosystem. This update ensures market-data-adapter-go has the latest schema definitions for proper type safety and service integration.

## Test Plan

- [x] Verify submodule points to valid commit
- [x] Run `bash scripts/validate-all.sh` - passes all checks
- [x] Confirm no breaking changes to market data services
- [x] Verify TODO.md exists and is valid
- [x] Test markdown linting passes

## Related Issues

- Epic TSE-0001: Foundation - Git Quality Standards
- Epic VIE-0002: Vendor Integration Engine - Phase 1

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
