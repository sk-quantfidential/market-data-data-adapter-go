# Pull Request: TSE-0001.12.0 - Multi-Instance Infrastructure Foundation

**Epic:** TSE-0001 - Foundation Services & Infrastructure
**Milestone:** TSE-0001.12.0 - Multi-Instance Infrastructure Foundation
**Branch:** `feature/TSE-0001.12.0-named-components-foundation`
**Status:** ✅ Ready for Merge

## Summary

This PR implements the **Phase 0** foundation for multi-instance infrastructure support in the market-data-adapter-go, providing:

1. **Instance-Aware Configuration**: `ServiceName` and `ServiceInstanceName` fields in Config
2. **Automatic Schema Derivation**: Smart schema naming based on service instance patterns
3. **Automatic Redis Namespace Derivation**: Instance-specific Redis namespace isolation
4. **Singleton and Multi-Instance Support**: Unified derivation logic for both patterns
5. **Comprehensive Test Coverage**: 19 unit tests covering all derivation scenarios

This is the **foundational layer** that enables the market data adapter to support multi-instance deployment with proper database and cache isolation for different market data sources (Coinmetrics, Bloomberg, etc.).

## Market-Data-Adapter-Go Repository Changes

### Commit Summary

**Total Commits**: 1 (Phase 0 foundation)

#### Phase 0: Multi-Instance Infrastructure Foundation
**Commit:** `c833bd9`
**Files Changed:** 5 files, 349 insertions(+), 6 deletions(-)

**Modified:**
- `internal/config/config.go` - Added instance-aware configuration fields
- `pkg/adapters/factory.go` - Added derivation functions and integration

**New:**
- `pkg/adapters/factory_test.go` - 19 comprehensive unit tests

**Changes:**
- Added `ServiceInstanceName` field to Config
- Added `SchemaName` field to Config (with automatic derivation)
- Added `RedisNamespace` field to Config (with automatic derivation)
- Implemented `deriveSchemaName(serviceName, instanceName)` function
- Implemented `deriveRedisNamespace(serviceName, instanceName)` function
- Integrated derivation into `NewMarketDataAdapter` factory

**Config Structure:**
```go
// internal/config/config.go
type Config struct {
    ServiceName         string
    ServiceInstanceName string
    SchemaName          string  // Auto-derived if empty
    RedisNamespace      string  // Auto-derived if empty
    PostgresURL         string
    RedisURL            string
    // ... other fields
}
```

### Derivation Logic

#### PostgreSQL Schema Name Derivation

**Singleton Services** (`serviceName == instanceName`):
- Pattern: Extract first two parts from service name, join with underscore
- Examples:
  - `market-data-simulator` → Schema: `market_data`
  - `market-data-adapter` → Schema: `market_data`

**Multi-Instance Services** (`serviceName != instanceName`):
- Pattern: Extract parts from instance name, join with underscore, lowercase entity identifier
- Examples:
  - `market-data-Coinmetrics` → Schema: `market_data_coinmetrics`
  - `market-data-Bloomberg` → Schema: `market_data_bloomberg`
  - `market-data-Binance` → Schema: `market_data_binance`

**Implementation:**
```go
func deriveSchemaName(serviceName, instanceName string) string {
    if serviceName == instanceName {
        // Singleton: "market-data-simulator" → "market_data"
        parts := strings.Split(serviceName, "-")
        if len(parts) >= 2 {
            return parts[0] + "_" + parts[1]
        }
        return serviceName
    }

    // Multi-instance: "market-data-Coinmetrics" → "market_data_coinmetrics"
    parts := strings.Split(instanceName, "-")
    if len(parts) >= 3 {
        return strings.ToLower(parts[0] + "_" + parts[1] + "_" + parts[2])
    } else if len(parts) >= 2 {
        return strings.ToLower(parts[0] + "_" + parts[1])
    }
    return strings.ToLower(instanceName)
}
```

#### Redis Namespace Derivation

**Singleton Services** (`serviceName == instanceName`):
- Pattern: Extract first two parts from service name, join with underscore
- Examples:
  - `market-data-simulator` → Namespace: `market_data`
  - Keys: `market_data:*`

**Multi-Instance Services** (`serviceName != instanceName`):
- Pattern: Extract service prefix, use entity identifier after colon
- Examples:
  - `market-data-Coinmetrics` → Namespace: `market_data:Coinmetrics`
  - Keys: `market_data:Coinmetrics:*`
  - `market-data-Bloomberg` → Namespace: `market_data:Bloomberg`
  - Keys: `market_data:Bloomberg:*`

**Implementation:**
```go
func deriveRedisNamespace(serviceName, instanceName string) string {
    if serviceName == instanceName {
        // Singleton: "market-data-simulator" → "market_data"
        parts := strings.Split(serviceName, "-")
        if len(parts) >= 2 {
            return parts[0] + "_" + parts[1]
        }
        return serviceName
    }

    // Multi-instance: "market-data-Coinmetrics" → "market_data:Coinmetrics"
    parts := strings.Split(instanceName, "-")
    if len(parts) >= 2 {
        return parts[0] + "_" + parts[1] + ":" + strings.Join(parts[2:], "-")
    }
    return instanceName
}
```

## Test Coverage

### Unit Tests (19 tests, all passing)

**Schema Derivation Tests** (8 tests):
- Singleton service pattern (market-data-simulator)
- Multi-instance patterns (market-data-Coinmetrics, market-data-Bloomberg)
- Edge cases (missing parts, special characters)

**Redis Namespace Derivation Tests** (8 tests):
- Singleton service pattern
- Multi-instance patterns with proper colon separation
- Entity identifier preservation (Coinmetrics, Bloomberg)

**Factory Integration Tests** (3 tests):
- Schema derivation integration in factory
- Redis namespace derivation integration in factory
- Backward compatibility (empty fields trigger derivation)

**Test File:** `pkg/adapters/factory_test.go` (256 lines)

## Environment Variables

**New Environment Variables:**
- `SERVICE_INSTANCE_NAME` - Instance identifier (optional, defaults to SERVICE_NAME)
- `SCHEMA_NAME` - Override for PostgreSQL schema (optional, auto-derived if empty)
- `REDIS_NAMESPACE` - Override for Redis namespace (optional, auto-derived if empty)

**Example Configurations:**

**Singleton Service:**
```bash
SERVICE_NAME=market-data-simulator
SERVICE_INSTANCE_NAME=market-data-simulator  # Optional, defaults to SERVICE_NAME
# Derived: schema=market_data, namespace=market_data
```

**Multi-Instance Service (Coinmetrics):**
```bash
SERVICE_NAME=market-data-adapter
SERVICE_INSTANCE_NAME=market-data-Coinmetrics
# Derived: schema=market_data_coinmetrics, namespace=market_data:Coinmetrics
```

**Multi-Instance Service (Bloomberg):**
```bash
SERVICE_NAME=market-data-adapter
SERVICE_INSTANCE_NAME=market-data-Bloomberg
# Derived: schema=market_data_bloomberg, namespace=market_data:Bloomberg
```

## Backward Compatibility

✅ **Fully backward compatible:**
- Existing deployments without `SERVICE_INSTANCE_NAME` continue to work
- `SERVICE_INSTANCE_NAME` defaults to `SERVICE_NAME` (singleton pattern)
- Explicit `SCHEMA_NAME` and `REDIS_NAMESPACE` override derivation
- No breaking changes to existing code

## Multi-Instance Use Cases

This foundation enables:

1. **Multiple Market Data Sources:**
   - `market-data-Coinmetrics` with dedicated schema `market_data_coinmetrics`
   - `market-data-Bloomberg` with dedicated schema `market_data_bloomberg`
   - `market-data-Binance` with dedicated schema `market_data_binance`

2. **Database Isolation:**
   - Each instance has its own PostgreSQL schema
   - No data collisions between different data sources
   - Independent schema migrations per instance

3. **Cache Isolation:**
   - Each instance has its own Redis namespace
   - `market_data:Coinmetrics:symbols:*` vs `market_data:Bloomberg:symbols:*`
   - No cache key collisions between instances

4. **Independent Scaling:**
   - Each data source can scale independently
   - Different resource allocations per source
   - Source-specific monitoring and alerting

## Files Changed

```
go.mod                       |   5 +-
go.sum                       |   1 +
internal/config/config.go    |  25 ++++-
pkg/adapters/factory.go      |  68 ++++++++++++
pkg/adapters/factory_test.go | 256 +++++++++++++++++++++++++++++++++++++++++++
5 files changed, 349 insertions(+), 6 deletions(-)
```

## Testing Instructions

### Run Unit Tests
```bash
# Run all tests
go test -v ./pkg/adapters/...

# Run specific test suite
go test -v ./pkg/adapters -run TestDeriveSchemaName
go test -v ./pkg/adapters -run TestDeriveRedisNamespace
go test -v ./pkg/adapters -run TestNewMarketDataAdapter
```

### Test Singleton Deployment
```bash
export SERVICE_NAME=market-data-simulator
export SERVICE_INSTANCE_NAME=market-data-simulator
go run cmd/example/main.go
# Expected: schema=market_data, namespace=market_data
```

### Test Multi-Instance Deployment
```bash
export SERVICE_NAME=market-data-adapter
export SERVICE_INSTANCE_NAME=market-data-Coinmetrics
go run cmd/example/main.go
# Expected: schema=market_data_coinmetrics, namespace=market_data:Coinmetrics
```

## Deployment Considerations

1. **Schema Creation:** PostgreSQL schemas must be created before deployment
2. **Permissions:** Service accounts need access to their derived schemas
3. **Redis Namespacing:** Existing keys unaffected due to namespace isolation
4. **Monitoring:** Log output includes derived schema and namespace for verification

## Next Steps (Future Work)

After merge:
- Phase 1: Update market-data-simulator-go to use instance-aware configuration
- Phase 2: Implement PostgreSQL schema creation in migrations
- Phase 3: Add Redis namespace support to repository implementations
- Phase 4: Deploy multi-instance market data adapters (Coinmetrics, Bloomberg, etc.)

## BDD Acceptance Criteria

✅ **All acceptance criteria met:**
- Service can be configured with instance name
- PostgreSQL schema is automatically derived based on instance pattern
- Redis namespace is automatically derived based on instance pattern
- Singleton and multi-instance patterns both supported
- Backward compatibility maintained
- Comprehensive unit tests passing (19/19)

---

**Ready for Merge:** Yes ✅
**Breaking Changes:** None
**Migration Required:** None
**Documentation Updated:** Yes (this PR document)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
