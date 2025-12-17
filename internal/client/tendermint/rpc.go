package tendermint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type StatusResponse struct {
	Result struct {
		NodeInfo struct {
			ID      string `json:"id"`
			Network string `json:"network"`
			Version string `json:"version"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHash   string `json:"latest_block_hash"`
			LatestBlockHeight string `json:"latest_block_height"`
			LatestBlockTime   string `json:"latest_block_time"`
		} `json:"sync_info"`
	} `json:"result"`
}

type BlockResponse struct {
	Result struct {
		BlockID struct {
			Hash string `json:"hash"`
		} `json:"block_id"`
		Block struct {
			Header struct {
				Height string `json:"height"`
				Time   string `json:"time"`
			} `json:"header"`
		} `json:"block"`
	} `json:"result"`
}

func New(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) do(endpoint string, out any) error {
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)

	resp, err := c.client.Get(url)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(
			"RPC returned %s for endpoint %s — this usually means the node has pruned the block or does not support historical queries",
			resp.Status,
			url,
		)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) Status() (any, error) {
	var result StatusResponse
	err := c.do("status", &result)
	return &result, err
}

func (c *Client) Block(height int64) (any, error) {
	end := "block"
	if height > 0 {
		end = fmt.Sprintf("block?height=%d", height)
	}
	var result BlockResponse
	err := c.do(end, &result)
	return &result, err
}

func (c *Client) Health() (any, error) {
	out := make(map[string]any)
	err := c.do("health", &out)
	return out, err
}
