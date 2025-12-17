package bitcoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	rpcURL  string
	user    string
	pass    string
	timeout time.Duration
	http    *http.Client
}

func New(rpcURL, user, pass string, timeout time.Duration) *Client {
	return &Client{
		rpcURL:  rpcURL,
		user:    user,
		pass:    pass,
		timeout: timeout,
		http:    &http.Client{Timeout: timeout},
	}
}

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      string `json:"id"`
}

type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *Client) call(method string, params []any, out any) error {
	body, err := json.Marshal(rpcRequest{
		JSONRPC: "1.0",
		Method:  method,
		Params:  params,
		ID:      "rpc-inspector",
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.rpcURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.user, c.pass)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Error != nil {
		return fmt.Errorf("bitcoin rpc error: %s", r.Error.Message)
	}

	return json.Unmarshal(r.Result, out)
}

func (c *Client) Status() (any, error) {
	var info map[string]any
	if err := c.call("getblockchaininfo", nil, &info); err != nil {
		return nil, err
	}

	return map[string]any{
		"chain":        info["chain"],
		"blocks":       info["blocks"],
		"headers":      info["headers"],
		"verification": info["verificationprogress"],
		"pruned":       info["pruned"],
		"difficulty":   info["difficulty"],
	}, nil
}

func (c *Client) Block(height int64) (any, error) {
	var heightResp int64
	if err := c.call("getblockcount", nil, &heightResp); err != nil {
		return nil, err
	}

	return map[string]int64{
		"blockHeight": heightResp,
	}, nil
}

func (c *Client) Health() (any, error) {
	var info map[string]any
	if err := c.call("getnetworkinfo", nil, &info); err != nil {
		return nil, err
	}

	return map[string]string{
		"status": "healthy",
	}, nil
}
