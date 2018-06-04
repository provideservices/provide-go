package provide

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

// GetBlockNumber retrieves the latest block known to the JSON-RPC client
func GetBlockNumber(networkID, rpcURL string) *uint64 {
	params := make([]interface{}, 0)
	var resp = &EthereumJsonRpcResponse{}
	Log.Debugf("Attempting to fetch latest block number via JSON-RPC eth_blockNumber method")
	err := InvokeJsonRpcClient(networkID, rpcURL, "eth_blockNumber", params, &resp)
	if err != nil {
		Log.Warningf("Failed to invoke eth_blockNumber method via JSON-RPC; %s", err.Error())
		return nil
	}
	blockNumber, err := hexutil.DecodeBig(resp.Result.(string))
	if err != nil {
		return nil
	}
	_blockNumber := blockNumber.Uint64()
	return &_blockNumber
}

// GetChainConfig parses the cached network config mapped to the given
// `networkID`, if one exists; otherwise, the mainnet chain config is returned.
func GetChainConfig(networkID, rpcURL string) *params.ChainConfig {
	if cfg, ok := chainConfigs[networkID]; ok {
		return cfg
	}
	return params.MainnetChainConfig
}

// GetChainID retrieves the current chainID via JSON-RPC
func GetChainID(networkID, rpcURL string) *big.Int {
	ethClient, err := ResolveEthClient(networkID, rpcURL)
	if err != nil {
		Log.Warningf("Failed to read network id for *ethclient.Client instance: %s; %s", ethClient, err.Error())
		return nil
	}
	chainID, err := ethClient.NetworkID(context.TODO())
	if err != nil {
		Log.Warningf("Failed to read network id for *ethclient.Client instance: %s; %s", ethClient, err.Error())
		return nil
	}
	if chainID != nil {
		Log.Debugf("Received chain id from *ethclient.Client instance: %s", ethClient, chainID)
	}
	return chainID
}

// GetGasPrice returns the gas price
func GetGasPrice(networkID, rpcURL string) *string {
	params := make([]interface{}, 0)
	var resp = &EthereumJsonRpcResponse{}
	Log.Debugf("Attempting to fetch gas price via JSON-RPC eth_gasPrice method")
	err := InvokeJsonRpcClient(networkID, rpcURL, "eth_gasPrice", params, &resp)
	if err != nil {
		Log.Warningf("Failed to invoke eth_gasPrice method via JSON-RPC; %s", err.Error())
		return nil
	}
	return stringOrNil(resp.Result.(string))
}

// GetLatestBlock retrieves the best block known to the JSON-RPC client
func GetLatestBlock(networkID, rpcURL string) (uint64, error) {
	status, err := GetNetworkStatus(networkID, rpcURL)
	if err != nil {
		return 0, err
	}
	return status.Block, nil
}

// GetNativeBalance retrieves a wallet's native currency balance
func GetNativeBalance(networkID, rpcURL, addr string) (*big.Int, error) {
	client, err := DialJsonRpc(networkID, rpcURL)
	if err != nil {
		return nil, err
	}
	return client.BalanceAt(context.TODO(), common.HexToAddress(addr), nil)
}

// GetNetworkStatus retrieves current metadata from the JSON-RPC client;
// returned struct includes block height, chainID, number of connected peers,
// protocol version, and syncing state.
func GetNetworkStatus(networkID, rpcURL string) (*NetworkStatus, error) {
	ethClient, err := ResolveEthClient(networkID, rpcURL)
	if err != nil || rpcURL == "" {
		meta := map[string]interface{}{
			"error": nil,
		}
		if err != nil {
			Log.Warningf("Failed to dial JSON-RPC host: %s; %s", rpcURL, err.Error())
			meta["error"] = err.Error()
		} else if rpcURL == "" {
			meta["error"] = errors.New("No 'full-node' JSON-RPC URL configured or resolvable")
		}
		return &NetworkStatus{
			State: stringOrNil("configuring"),
			Meta:  meta,
		}, nil
	}

	syncProgress, err := GetSyncProgress(ethClient)
	if err != nil {
		Log.Warningf("Failed to read sync progress using JSON-RPC host; %s", err.Error())
		return nil, err
	}
	var state string
	var block uint64   // current block; will be less than height while syncing in progress
	var height *uint64 // total number of blocks
	chainID := GetChainID(networkID, rpcURL)
	peers := GetPeerCount(networkID, rpcURL)
	protocolVersion := GetProtocolVersion(networkID, rpcURL)
	var syncing = false
	if syncProgress == nil {
		state = "synced"
		hdr, err := ethClient.HeaderByNumber(context.TODO(), nil)
		if err != nil && hdr == nil {
			Log.Warningf("Failed to read latest block header for using JSON-RPC host; %s", err.Error())
			var jsonRpcResponse = &EthereumJsonRpcResponse{}
			err = InvokeJsonRpcClient(networkID, rpcURL, "eth_getBlockByNumber", []interface{}{"latest", true}, &jsonRpcResponse)
			if err != nil {
				Log.Warningf("Failed to read latest block header for using JSON-RPC host; %s", err.Error())
				err = InvokeJsonRpcClient(networkID, rpcURL, "eth_getBlockByNumber", []interface{}{"earliest", true}, &jsonRpcResponse)
				if err != nil {
					Log.Warningf("Failed to read earliest block header for using JSON-RPC host; %s", err.Error())
					return nil, err
				}
			}
			if jsonRpcResponse.Result != nil {
				Log.Debugf("Got JSON-RPC response; %s", jsonRpcResponse.Result)
			}
		}
		block = hdr.Number.Uint64()
	} else {
		block = syncProgress.CurrentBlock
		height = &syncProgress.HighestBlock
		syncing = true
	}
	return &NetworkStatus{
		Block:           block,
		Height:          height,
		ChainID:         chainID,
		PeerCount:       peers,
		ProtocolVersion: protocolVersion,
		State:           stringOrNil(state),
		Syncing:         syncing,
		Meta:            map[string]interface{}{},
	}, nil
}

// GetPeerCount returns the number of peers currently connected to the JSON-RPC client
func GetPeerCount(networkID, rpcURL string) uint64 {
	var peerCount uint64
	params := make([]interface{}, 0)
	var resp = &EthereumJsonRpcResponse{}
	Log.Debugf("Attempting to fetch peer count via net_peerCount method via JSON-RPC")
	err := InvokeJsonRpcClient(networkID, rpcURL, "net_peerCount", params, &resp)
	if err != nil {
		Log.Debugf("Attempting to fetch peer count via parity_netPeers method via JSON-RPC")
		err := InvokeJsonRpcClient(networkID, rpcURL, "parity_netPeers", params, &resp)
		Log.Warningf("Failed to invoke parity_netPeers method via JSON-RPC; %s", err.Error())
		return 0
	}
	if peerCountStr, ok := resp.Result.(string); ok {
		peerCount, err = hexutil.DecodeUint64(peerCountStr)
		if err != nil {
			return 0
		}
	}
	return peerCount
}

// GetProtocolVersion returns the JSON-RPC client protocol version
func GetProtocolVersion(networkID, rpcURL string) *string {
	params := make([]interface{}, 0)
	var resp = &EthereumJsonRpcResponse{}
	Log.Debugf("Attempting to fetch protocol version via JSON-RPC eth_protocolVersion method")
	err := InvokeJsonRpcClient(networkID, rpcURL, "eth_protocolVersion", params, &resp)
	if err != nil {
		Log.Debugf("Attempting to fetch protocol version via JSON-RPC net_version method")
		err := InvokeJsonRpcClient(networkID, rpcURL, "net_version", params, &resp)

		Log.Warningf("Failed to invoke eth_protocolVersion method via JSON-RPC; %s", err.Error())
		return nil
	}
	return stringOrNil(resp.Result.(string))
}

// GetCode retrieves the code stored at the named address in the given scope;
// scope can be a block number, latest, earliest or pending
func GetCode(networkID, rpcURL, addr, scope string) (*string, error) {
	params := make([]interface{}, 0)
	params = append(params, addr)
	params = append(params, scope)
	var resp = &EthereumJsonRpcResponse{}
	Log.Debugf("Attempting to fetch code from %s via eth_getCode JSON-RPC method", addr)
	err := InvokeJsonRpcClient(networkID, rpcURL, "eth_getCode", params, &resp)
	if err != nil {
		Log.Warningf("Failed to invoke eth_getCode method via JSON-RPC; %s", err.Error())
		return nil, err
	}
	return stringOrNil(resp.Result.(string)), nil
}

// GetSyncProgress retrieves the status of the current network sync
func GetSyncProgress(client *ethclient.Client) (*ethereum.SyncProgress, error) {
	progress, err := client.SyncProgress(context.TODO())
	if err != nil {
		Log.Warningf("Failed to read sync progress for *ethclient.Client instance: %s; %s", client, err.Error())
		return nil, err
	}
	if progress != nil {
		Log.Debugf("Latest synced block reported by *ethclient.Client instance: %v [of %v]", client, progress.CurrentBlock, progress.HighestBlock)
	}
	return progress, nil
}

// GetTokenBalance retrieves a token balance for a specific token contract and network address
func GetTokenBalance(networkID, rpcURL, tokenAddr, addr string, contractABI interface{}) (*big.Int, error) {
	var balance *big.Int
	abi, err := parseContractABI(contractABI)
	if err != nil {
		return nil, err
	}
	client, err := DialJsonRpc(networkID, rpcURL)
	gasPrice, _ := client.SuggestGasPrice(context.TODO())
	to := common.HexToAddress(tokenAddr)
	msg := ethereum.CallMsg{
		From:     common.HexToAddress(addr),
		To:       &to,
		Gas:      0,
		GasPrice: gasPrice,
		Value:    nil,
		Data:     common.FromHex(HashFunctionSelector("balanceOf(address)")),
	}
	result, _ := client.CallContract(context.TODO(), msg, nil)
	if method, ok := abi.Methods["balanceOf"]; ok {
		method.Outputs.Unpack(&balance, result)
		if balance != nil {
			symbol, _ := GetTokenSymbol(networkID, rpcURL, addr, tokenAddr, contractABI)
			Log.Debugf("Read %s token balance (%v) from token contract address: %s", symbol, balance, addr)
		}
	} else {
		Log.Warningf("Unable to read balance of unsupported token contract address: %s", tokenAddr)
	}
	return balance, nil
}

// GetTokenSymbol attempts to retrieve the symbol of a token presumed to be deployed at the given token contract address
func GetTokenSymbol(networkID, rpcURL, from, tokenAddr string, contractABI interface{}) (*string, error) {
	client, err := DialJsonRpc(networkID, rpcURL)
	if err != nil {
		return nil, err
	}
	_abi, err := parseContractABI(contractABI)
	if err != nil {
		return nil, err
	}
	to := common.HexToAddress(tokenAddr)
	msg := ethereum.CallMsg{
		From:     common.HexToAddress(from),
		To:       &to,
		Gas:      0,
		GasPrice: big.NewInt(0),
		Value:    nil,
		Data:     common.FromHex(HashFunctionSelector("symbol()")),
	}
	result, _ := client.CallContract(context.TODO(), msg, nil)
	var symbol string
	if method, ok := _abi.Methods["symbol"]; ok {
		err = method.Outputs.Unpack(&symbol, result)
		if err != nil {
			Log.Warningf("Failed to read token symbol from deployed token contract %s; %s", tokenAddr, err.Error())
		}
	}
	return stringOrNil(symbol), nil
}

// TraceTx returns the VM traces; requires parity JSON-RPC client and the node must
// be configured with `--fat-db on --tracing on --pruning archive`
func TraceTx(networkID, rpcURL string, hash *string) (interface{}, error) {
	var addr = *hash
	if !strings.HasPrefix(addr, "0x") {
		addr = fmt.Sprintf("0x%s", addr)
	}
	params := make([]interface{}, 0)
	params = append(params, addr)
	var result = &EthereumTxTraceResponse{}
	Log.Debugf("Attempting to trace tx via trace_transaction method via JSON-RPC; tx hash: %s", addr)
	err := InvokeJsonRpcClient(networkID, rpcURL, "trace_transaction", params, &result)
	if err != nil {
		Log.Warningf("Failed to invoke trace_transaction method via JSON-RPC; %s", err.Error())
		return nil, err
	}
	return result, nil
}

// GetTxReceipt retrieves the full transaction receipt via JSON-RPC given the transaction hash
func GetTxReceipt(networkID, rpcURL, txHash, from string) (*types.Receipt, error) {
	var err error
	var receipt *types.Receipt
	client, err := DialJsonRpc(networkID, rpcURL)
	// FIXME-- make sure 0-prefixed and non-prefixed hashes work... txHash := fmt.Sprintf("0x%s", *t.Hash)
	// FIXME-- set a timeout on the following code that currently blocks util the tx receipt is retrieved:
	Log.Debugf("Attempting to retrieve tx receipt for broadcast tx: %s", txHash)
	err = ethereum.NotFound
	for receipt == nil && err == ethereum.NotFound {
		Log.Debugf("Retrieving broadcast tx receipt for tx: %s", txHash)
		receipt, err = client.TransactionReceipt(context.TODO(), common.HexToHash(txHash))
	}
	return receipt, err
}