package util

import (
	"encoding/json"
	"fmt"

	tm "github.com/afacchi0/rpc-inspector/internal/client/tendermint"
)

func Print(proto, endpoint string, v interface{}, jsonOutput bool) {
	if jsonOutput {
		out, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(out))
		return
	}

	switch proto {

	case "tendermint":
		switch endpoint {
		case "status":
			printTendermintStatus(v)
		case "health":
			printTendermintHealth(v)
		case "block":
			printTendermintBlock(v)
		default:
			fmt.Printf("%v\n", v)
		}

	case "ethereum":
		switch endpoint {
		case "status":
			printEthereumStatus(v.(map[string]interface{}))
		case "health":
			printEthereumHealth(v.(map[string]string))
		case "block":
			printEthereumBlock(v)
		default:
			fmt.Printf("%v\n", v)
		}

	case "bitcoin":
		switch endpoint {
		case "status":
			printBitcoinStatus(v)
		case "health":
			printBitcoinHealth(v)
		case "block":
			printBitcoinBlock(v)
		}

	default:
		fmt.Printf("%v\n", v)
	}
}

func printTendermintStatus(v interface{}) {
	s := v.(*tm.StatusResponse)
	fmt.Println("=== Tendermint Status ===")
	fmt.Println("Node ID:     ", s.Result.NodeInfo.ID)
	fmt.Println("Network:     ", s.Result.NodeInfo.Network)
	fmt.Println("Version:     ", s.Result.NodeInfo.Version)
	fmt.Println("Latest Block:", s.Result.SyncInfo.LatestBlockHeight)
	fmt.Println("Block Time:  ", s.Result.SyncInfo.LatestBlockTime)
}

func printTendermintHealth(v interface{}) {
	fmt.Println("=== Tendermint Health ===")
	fmt.Println("Status: healthy")
}

func printTendermintBlock(v interface{}) {
	b := v.(*tm.BlockResponse)
	fmt.Println("=== Tendermint Block ===")
	fmt.Println("Height:", b.Result.Block.Header.Height)
	fmt.Println("Hash:  ", b.Result.BlockID.Hash)
	fmt.Println("Time:  ", b.Result.Block.Header.Time)
}

func printEthereumStatus(m map[string]interface{}) {
	fmt.Println("=== Ethereum Status ===")
	fmt.Println("Client:      ", m["client"])
	fmt.Println("Syncing:     ", m["syncing"])
	fmt.Println("Latest Block:", m["latestBlock"])
}

func printEthereumHealth(m map[string]string) {
	fmt.Println("=== Ethereum Health ===")
	fmt.Println("Status:", m["status"])
}

func printEthereumBlock(v interface{}) {
	m := v.(map[string]uint64)
	fmt.Println("=== Ethereum block===")
	fmt.Println("Latest Block:", m["blockNumber"])
}

func printBitcoinStatus(v interface{}) {
	m := v.(map[string]interface{})
	fmt.Println("=== Bitcoin Status ===")
	fmt.Println("Chain:        ", m["chain"])
	fmt.Println("Blocks:       ", m["blocks"])
	fmt.Println("Headers:      ", m["headers"])
	fmt.Println("Verification: ", m["verification"])
	fmt.Println("Pruned:       ", m["pruned"])
	fmt.Println("Difficulty:   ", m["difficulty"])
}

func printBitcoinHealth(v interface{}) {
	m := v.(map[string]string)
	fmt.Println("=== Bitcoin Health ===")
	fmt.Println("Status:", m["status"])
}

func printBitcoinBlock(v interface{}) {
	m := v.(map[string]int64)
	fmt.Println("=== Bitcoin Block ===")
	fmt.Println("Height:", m["blockHeight"])
}
