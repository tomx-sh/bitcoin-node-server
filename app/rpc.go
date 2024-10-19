package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type RpcResponse struct {
	Result interface{} `json:"result"`
	Error  *RpcError   `json:"error"`
	Id     interface{} `json:"id"`
}

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Function to perform an RPC call
func Rpc(method string, params []interface{}) (*RpcResponse, error) {
	var RPCUsername = os.Getenv("RPC_USERNAME")
	var RPCPassword = os.Getenv("RPC_PASSWORD")
	var RPCURL = os.Getenv("RPC_URL")

	if RPCUsername == "" || RPCPassword == "" || RPCURL == "" {
		return &RpcResponse{Error: &RpcError{Code: -1, Message: "RPC credentials or URL not set"}}, nil
	}

	payload := map[string]interface{}{
		"jsonrpc": "1.0",
		"id":      "curltest",
		"method":  method,
		"params":  params,
	}

	// Marshal the payload to JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create the request
	req, err := http.NewRequest("POST", RPCURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "text/plain")
	req.SetBasicAuth(RPCUsername, RPCPassword)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var rpcResponse RpcResponse
	err = json.NewDecoder(resp.Body).Decode(&rpcResponse)
	if err != nil {
		return nil, err
	}

	return &rpcResponse, nil
}
