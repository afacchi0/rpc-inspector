package ethereum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	rpcURL  string
	timeout time.Duration
	http    *http.Client
}

func New(rpcURL string) *Client {
	t := 5 * time.Second
	return &Client{
		rpcURL:  rpcURL,
		timeout: t,
		http:    &http.Client{Timeout: t},
	}
}

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type rpcEnvelope struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *Client) call(method string, out any) error {
	body, _ := json.Marshal(rpcRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  []any{},
		ID:      1,
	})

	req, _ := http.NewRequest("POST", c.rpcURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var env rpcEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return err
	}
	if env.Error != nil {
		return fmt.Errorf(env.Error.Message)
	}

	return json.Unmarshal(env.Result, out)
}

func parseHexUint64(s string) (uint64, error) {
	return strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 64)
}

func (c *Client) Status() (any, error) {
	var clientVersion string
	if err := c.call("web3_clientVersion", &clientVersion); err != nil {
		return nil, err
	}

	// eth_syncing returns false OR an object
	var syncingResp any
	if err := c.call("eth_syncing", &syncingResp); err != nil {
		return nil, err
	}
	isSyncing := syncingResp != false

	var blockHex string
	if err := c.call("eth_blockNumber", &blockHex); err != nil {
		return nil, err
	}

	blockNum, err := parseHexUint64(blockHex)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"client":      clientVersion,
		"syncing":     isSyncing,
		"latestBlock": blockNum,
	}, nil
}

func (c *Client) Health() (any, error) {
	var listening bool
	if err := c.call("net_listening", &listening); err != nil {
		return nil, err
	}

	if !listening {
		return nil, fmt.Errorf("ethereum rpc unhealthy: net_listening=false")
	}

	return map[string]string{
		"status": "healthy",
	}, nil
}

func (c *Client) Block(height int64) (any, error) {
	var hex string
	if err := c.call("eth_blockNumber", &hex); err != nil {
		return nil, err
	}

	n, err := parseHexUint64(hex)
	if err != nil {
		return nil, err
	}

	return map[string]uint64{"blockNumber": n}, nil
}
