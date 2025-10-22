# Market Data Adapter (Go)

**Component**: market-data-adapter-go
**Domain**: Market data, pricing, scenarios
**Epic**: TSE-0001.4 (Market Data Foundation)
**Tech Stack**: Go, PostgreSQL, Redis
**Schema Namespace**: `market_data`

## Purpose

The Market Data Adapter provides data persistence services for the market-data-simulator-go component, following Clean Architecture principles. It exposes domain-driven APIs for price feed management, historical data storage, and scenario orchestration while abstracting database implementation details.

## Architecture Compliance

**Clean Architecture**:
- Exposes business domain concepts, not database artifacts
- Provides market-data-specific APIs tailored to pricing and scenario needs
- Maintains complete separation from market-data-simulator business logic
- Uses shared infrastructure with logical namespace isolation

**Domain Focus**:
- Real-time price feeds (BTC/USD, ETH/USD, USDT pairs)
- Historical OHLCV data with time-series optimization
- Market microstructure (bid/ask spreads, volumes)
- Scenario definitions and chaos injection configurations
- Price manipulation state and correlation tracking

## Data Requirements

### Real-Time Market Data
- **Price Feeds**: Real-time price distribution for BTC/USD, ETH/USD, USDT pairs
- **Market Microstructure**: Bid/ask spreads, volumes, market depth
- **Price History**: Historical price storage with time-series optimization
- **Volatility Modeling**: Simple volatility tracking and simulation
- **Feed Performance**: Metrics for feed latency and reliability

### Scenario Management
- **Scenario Definitions**: Chaos injection and market manipulation scenarios
- **Active State**: Current scenario execution state and parameters
- **Price Manipulation**: Coordinated price manipulation tracking
- **Correlation Data**: Cross-asset correlation and scenario impact tracking

### Storage Patterns
- **Redis**: Real-time prices (ring buffers), active scenario state, feed cache
- **PostgreSQL**: Historical OHLCV, scenario definitions, market metadata

## API Design Principles

### Domain-Driven APIs
The adapter exposes market data and pricing concepts, not database implementation:

**Good Examples**:
```go
PublishPriceFeed(symbol, price) -> FeedResult
GetMarketData(symbol, timeRange) -> OHLCV[]
CreateScenario(scenarioSpec) -> ScenarioID
ExecuteMarketManipulation(scenario) -> ManipulationResult
GetMarketDepth(symbol) -> OrderBookDepth
```

**Avoid Database Artifacts**:
```go
// Don't expose these
GetPriceTable() -> []PriceRow
UpdateScenarioRecord(id, fields) -> bool
QueryHistoricalData(sql) -> ResultSet
```

## Technology Standards

### Database Conventions
- **PostgreSQL**: snake_case for tables, columns, functions
- **Redis**: kebab-case with `market:` namespace prefix
- **Go**: PascalCase for public APIs, camelCase for internal functions

### Performance Requirements
- **Real-time Feeds**: Sub-millisecond price feed publication
- **Historical Queries**: Efficient time-series queries for OHLCV data
- **Scenario Execution**: Low-latency scenario state management
- **Data Retention**: Optimized storage for high-frequency market data

## Integration Points

### Serves
- **Primary**: market-data-simulator-go
- **Integration**: Provides market data to risk-monitor-py, trading-system-engine-py, exchange-simulator-go

### Dependencies
- **Shared Infrastructure**: Single PostgreSQL + Redis instances
- **Protocol Buffers**: Via protobuf-schemas for market data definitions
- **Service Discovery**: Via orchestrator-docker configuration

## Supported Trading Pairs

### Initial Pairs
- **BTC/USD**: Bitcoin to US Dollar
- **ETH/USD**: Ethereum to US Dollar
- **BTC/USDT**: Bitcoin to Tether
- **ETH/USDT**: Ethereum to Tether

### Market Data Types
- **L1 Quotes**: Best bid/ask prices and volumes
- **OHLCV**: Open, High, Low, Close, Volume candles
- **Market Depth**: Order book depth and liquidity
- **Trade Ticks**: Individual trade executions

## Chaos Testing Integration

### Scenario Types
- **Price Manipulation**: Coordinated price movement scenarios
- **Feed Disruption**: Market data feed reliability testing
- **Volatility Injection**: Extreme volatility scenario simulation
- **Correlation Stress**: Cross-asset correlation breakdown testing

## Development Status

**Repository**: Active (Repository Created)
**Branch**: feature/TSE-0003.0-data-adapter-foundation
**Epic Progress**: TSE-0001.4 (Market Data Foundation) - Not Started

## Next Steps

1. **Component Configuration**: Add `.claude/` configuration for market-data-specific patterns
2. **Schema Design**: Design market_data schema in PostgreSQL with time-series optimization
3. **API Definition**: Define price feed and scenario management APIs
4. **Implementation**: Implement adapter with comprehensive testing
5. **Integration**: Connect with market-data-simulator-go component

## Configuration Management

- **Shared Configuration**: project-plan/.claude/ for global architecture patterns
- **Component Configuration**: .claude/ directory for market-data-specific settings (to be created)
- **Database Schema**: `market_data` namespace with time-series and scenario optimization

---

**Epic Context**: TSE-0001 Foundation Services & Infrastructure
**Last Updated**: 2025-09-18
**Architecture**: Clean Architecture with domain-driven data persistence