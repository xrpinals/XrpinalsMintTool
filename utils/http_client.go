package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	Timeout int64
}

type RpcReq struct {
	Id     int           `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type RpcResp struct {
	Id      int    `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}

func (h HttpClient) HttpPost(url string, rpcReq RpcReq) ([]byte, error) {
	c := &http.Client{
		Transport: &http.Transport{},
		Timeout:   time.Second * time.Duration(h.Timeout),
	}

	bs, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", url, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Set("Connection", "close")
	r.Close = true

	rs, err := c.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		_ = rs.Body.Close()
		return nil, err
	}

	_ = rs.Body.Close()
	return body, nil
}
