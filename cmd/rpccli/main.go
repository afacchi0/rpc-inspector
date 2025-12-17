package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/afacchi0/rpc-inspector/internal/client/bitcoin"
	"github.com/afacchi0/rpc-inspector/internal/client/ethereum"
	"github.com/afacchi0/rpc-inspector/internal/client/tendermint"
	"github.com/afacchi0/rpc-inspector/internal/util"
)

type ProtocolClient interface {
	Status() (any, error)
	Block(height int64) (any, error)
	Health() (any, error)
}

func main() {
	rpcURL := flag.String("rpc", "http://localhost:26657", "RPC endpoint")
	proto := flag.String("type", "tendermint", "RPC type: tendermint | ethereum | bitcoin")
	endpoint := flag.String("endpoint", "status", "RPC endpoint: status | health | block")
	height := flag.Int64("height", 0, "Block height for block endpoint")
	jsonOutput := flag.Bool("json", false, "Print raw JSON output")
	timeout := flag.Duration("timeout", 5*time.Second, "HTTP timeout")
	rpcUser := flag.String("rpcuser", "", "Bitcoin RPC username")
	rpcPass := flag.String("rpcpass", "", "Bitcoin RPC password")

	flag.Parse()

	var c ProtocolClient
	var data any
	var err error

	switch *proto {
	case "tendermint":
		c = tendermint.New(*rpcURL, *timeout)
		switch *endpoint {
		case "status":
			data, err = c.Status()
		case "block":
			data, err = c.Block(*height)
		case "health":
			data, err = c.Health()
		default:
			fmt.Println("Unknown Tendermint endpoint: " + *endpoint)
			os.Exit(1)
		}
	case "ethereum":
		c = ethereum.New(*rpcURL)
		switch *endpoint {
		case "status":
			data, err = c.Status()
		case "health":
			data, err = c.Health()
		case "block":
			data, err = c.Block(0)
		default:
			fmt.Println("Unknown Ethereum endpoint: " + *endpoint)
			os.Exit(1)
		}
	case "bitcoin":
		c = bitcoin.New(*rpcURL, *rpcUser, *rpcPass, *timeout)
		switch *endpoint {
		case "status":
			data, err = c.Status()
		case "health":
			data, err = c.Health()
		case "block":
			data, err = c.Block(0)
		default:
			fmt.Println("Unknown Bitcoin endpoint: " + *endpoint)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown type:", *proto)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	util.Print(*proto, *endpoint, data, *jsonOutput)
}
