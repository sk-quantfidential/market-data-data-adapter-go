# Market Data Adapter

Data adapter service for market data in the Trading Ecosystem.

## Status

🚧 **In Development** - Foundation phase

## Purpose

Provides data access and persistence layer for market data services, implementing the adapter pattern to separate business logic from data concerns.

## Architecture

- **Pattern**: Clean Architecture with Adapter pattern
- **Language**: Go 1.24+
- **Dependencies**: Protocol Buffers, gRPC

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## Testing

```bash
go test ./...
```

## Related Services

- **market-data-simulator-go**: Market data generation and distribution service

## Project Context

Part of Epic TSE-0001: Foundation Services & Infrastructure
