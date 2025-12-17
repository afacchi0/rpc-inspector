# rpc-inspector

`rpc-inspector` is a small, protocol-aware CLI tool for querying and inspecting blockchain RPC endpoints in a consistent way.

The tool currently supports:
- **Tendermint / CometBFT–based chains** (Cosmos SDK, etc.)
- **Ethereum JSON-RPC**
- **Bitcoin Core JSON-RPC**

---

## Features

- Single CLI interface for multiple blockchain RPC protocols
- Protocol-aware endpoints (`status`, `health`, `block`)
- Human-readable output by default
- Raw JSON output with `--json`
- Clean separation of concerns:
  - RPC clients return data only
  - Output formatting handled centrally
- Designed to mirror how real infra tooling works

---

## Installation

Clone the repository and build locally:

```bash
git clone https://github.com/afacchi0/rpc-inspector.git
cd rpc-inspector
go build ./cmd/rpccli
```

---
## Usage

Syntax

```bash
./rpccli --type <protocol> --rpc <rpc-url> --endpoint <endpoint>
```

Example

```bash
./rpccli --type tendermint --endpoint block --rpc https://rpc.cosmos.directory/cosmoshub
=== Tendermint Block ===
Height: 28904229
Hash:   5F8166A0F99F2617626758080A4A04851361B8280147EF80E22AE7471DE8B57F
Time:   2025-12-17T22:16:36.125948183Z

./rpccli   --type ethereum   --endpoint status  --rpc https://eth.llamarpc.com
=== Ethereum Status ===
Client:       rpc-proxy
Syncing:      false
Latest Block: 24035259

./rpccli   --type bitcoin   --rpc http://127.0.0.1:8332   --rpcuser rpcuser   --rpcpass rpcpass   --endpoint status
=== Bitcoin Status ===
Chain:         main
Blocks:        515676
Headers:       928321
Verification:  0.288481716672891
Pruned:        true
Difficulty:    3.462542391191563e+12
```

> **Note**  
The Bitcoin examples above assume a **locally running Bitcoin Core node** (`bitcoind`) started via Docker. Bitcoin Core JSON-RPC requires authentication and does not provide public, unauthenticated endpoints.