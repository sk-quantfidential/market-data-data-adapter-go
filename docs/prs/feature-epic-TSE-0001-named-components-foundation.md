# Pull Request: Multi-Instance Infrastructure Foundation

**Epic:** TSE-0001 - Foundation Services & Infrastructure
**Milestone:** TSE-0001.12.0 - Multi-Instance Infrastructure Foundation
**Branch:** `feature/epic-TSE-0001-named-components-foundation`
**Status:** ✅ Ready for Review

## Summary

This PR implements the multi-instance infrastructure foundation for the market-data-adapter-go service, establishing the groundwork for multi-instance deployment patterns across the trading ecosystem.

### Key Changes

1. **Named Service Instances**: Explicit instance identification via `SERVICE_INSTANCE_NAME`
2. **Config-Level Data Adapter**: Centralized DataAdapter initialization and lifecycle management
3. **Instance-Aware Service Discovery**: Enhanced registration with instance metadata
4. **Multi-Instance Ready**: Foundation for horizontal scaling

The market-data-adapter implements the multi-instance foundation pattern to support the broader ecosystem's deployment model.

## What Changed

### Phase 0: Repository Initialization
- Created minimal market-data-adapter-go repository structure
- Added protocol buffer integration (submodule)
- Configured Go modules and dependencies
- Added gitignore exclusions for configuration files

### Phase 1: Multi-Instance Infrastructure
- Added `ServiceInstanceName` field to Config struct
- Implemented `SERVICE_INSTANCE_NAME` environment variable with backward compatibility
- Integrated with market-data-data-adapter-go package
- Added DataAdapter lifecycle management at config level

### Code Structure
```
market-data-adapter-go/
├── internal/
│   ├── cache/              # Redis cache integration
│   ├── config/             # Configuration with multi-instance support
│   └── database/           # PostgreSQL database integration
├── pkg/
│   ├── adapters/           # Repository implementations
│   ├── interfaces/         # Repository interfaces
│   └── models/             # Domain models
└── proto/                  # Protocol buffer definitions (submodule)
```

## Testing

All validation checks pass:
- ✅ `scripts/validate-all.sh` - All 7 checks passing
- ✅ Repository structure validated
- ✅ Protocol buffer integration verified
- ✅ Configuration loading tested
- ✅ Multi-instance configuration validated

## Migration Notes

**Environment Variables:**
- `SERVICE_INSTANCE_NAME` - Optional, defaults to `SERVICE_NAME`
- For singleton deployments: No changes required
- For multi-instance deployments: Set `SERVICE_INSTANCE_NAME` to unique instance identifier

**Backward Compatibility:**
- Existing deployments continue to work without changes
- Instance name defaults to service name for singleton services

## Related PRs

Part of Epic TSE-0001: Foundation Services & Infrastructure
- Implements multi-instance pattern across ecosystem
- Aligns with audit-correlator-go, custodian-simulator-go, exchange-simulator-go patterns

## Checklist

- [x] Code follows repository conventions
- [x] All validation checks pass
- [x] PR documentation complete
- [x] Multi-instance foundation implemented
- [x] Backward compatibility maintained
- [x] Configuration validated
