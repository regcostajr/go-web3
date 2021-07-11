/********************************************************************************
   This file is part of go-web3.
   go-web3 is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-web3 is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-web3.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

/**
 * @file eth.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package eth

import (
	"fmt"
	"math/big"

	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

type Eth struct {
	provider providers.ProviderInterface
}

// NewEth - creates a new eth instance
func NewEth(provider providers.ProviderInterface) *Eth {
	eth := new(Eth)
	eth.provider = provider
	return eth
}

// Reference: https://eth.wiki/json-rpc/API#eth_syncing
func (eth *Eth) IsSyncing() (*dto.SyncingResponse, error) {
	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_syncing", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToSyncingResponse()
}

// Reference: https://eth.wiki/json-rpc/API#eth_coinbase
func (eth *Eth) GetCoinbase() (string, error) {
	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_coinbase", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

// Reference: https://eth.wiki/json-rpc/API#eth_mining
func (eth *Eth) IsMining() (bool, error) {
	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_mining", nil)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

// Reference: https://eth.wiki/json-rpc/API#eth_hashrate
func (eth *Eth) GetHashRate() (int64, error) {
	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_hashrate", nil)

	if err != nil {
		return -1, err
	}

	return pointer.ToInt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_gasprice
func (eth *Eth) GetGasPrice() (*big.Int, error) {
	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_gasPrice", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_accounts
func (eth *Eth) ListAccounts() ([]string, error) {
	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_accounts", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()
}

// Reference: https://eth.wiki/json-rpc/API#eth_blocknumber
func (eth *Eth) GetBlockNumber() (int64, error) {
	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_blockNumber", nil)

	if err != nil {
		return -1, err
	}

	return pointer.ToInt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_getbalance
func (eth *Eth) GetBalance(address string, defaultBlockParameter string) (*big.Int, error) {
	params := make([]string, 2)
	params[0] = address
	params[1] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getBalance", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_gettransactionaccount
func (eth *Eth) GetTransactionCount(address string, defaultBlockParameter string) (int64, error) {
	params := make([]string, 2)
	params[0] = address
	params[1] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getTransactionCount", params)

	if err != nil {
		return -1, err
	}

	return pointer.ToInt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_getstorageat
//TODO Default block number must be an object
func (eth *Eth) GetStorageAt(address string, position int64, defaultBlockParameter string) (string, error) {
	params := make([]string, 3)
	params[0] = address
	params[1] = fmt.Sprintf("0x%x", position)
	params[2] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getStorageAt", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

// Reference: https://eth.wiki/json-rpc/API#eth_estimategas
func (eth *Eth) EstimateGas(transaction *dto.TransactionParameters) (*big.Int, error) {
	params := make([]*dto.RequestTransactionParameters, 1)

	params[0] = transaction.Transform()

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(&pointer, "eth_estimateGas", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_gettransactionbyhash
func (eth *Eth) GetTransactionByHash(hash string) (*dto.TransactionResponse, error) {
	params := make([]string, 1)
	params[0] = hash

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getTransactionByHash", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()

}

// Reference: https://eth.wiki/json-rpc/API#eth_getTransactionByBlockNumberAndIndex
func (eth *Eth) GetTransactionByBlockHashAndIndex(hash string, index int64) (*dto.TransactionResponse, error) {
	params := make([]string, 2)
	params[0] = hash
	params[1] = fmt.Sprintf("0x%x", index)

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getTransactionByBlockHashAndIndex", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()
}

// Reference: https://eth.wiki/json-rpc/API#eth_getTransactionByBlockNumberAndIndex
func (eth *Eth) GetTransactionByBlockNumberAndIndex(blockIndex int64, index int64) (*dto.TransactionResponse, error) {
	params := make([]string, 2)
	params[0] = fmt.Sprintf("0x%x", blockIndex)
	params[1] = fmt.Sprintf("0x%x", index)

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getTransactionByBlockNumberAndIndex", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()
}

// Reference: https://eth.wiki/json-rpc/API#eth_sendtransaction
func (eth *Eth) SendTransaction(transaction *dto.TransactionParameters) (string, error) {
	params := make([]*dto.RequestTransactionParameters, 1)
	params[0] = transaction.Transform()

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(&pointer, "eth_sendTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

// Reference: https://eth.wiki/json-rpc/API#eth_signtransaction
func (eth *Eth) SignTransaction(transaction *dto.TransactionParameters) (string, error) {
	params := make([]*dto.RequestTransactionParameters, 1)
	params[0] = transaction.Transform()

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(&pointer, "eth_signTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

// Reference: https://eth.wiki/json-rpc/API#eth_sendRawTransaction
func (eth *Eth) SendRawTransaction(signedTransaction string) (string, error) {
	params := make([]string, 1)
	params[0] = signedTransaction

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(&pointer, "eth_sendRawTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

// Reference: https://eth.wiki/json-rpc/API#eth_call
func (eth *Eth) Call(transaction *dto.TransactionParameters) (*dto.RequestResult, error) {
	params := make([]interface{}, 2)
	params[0] = transaction.Transform()
	params[1] = LATEST

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(&pointer, "eth_call", params)

	if err != nil {
		return nil, err
	}

	//TODO Check return response
	return pointer, err
}

// Reference: https://eth.wiki/json-rpc/API#eth_compilesolidity
func (eth *Eth) CompileSolidity(sourceCode string) (string, error) {
	params := make([]string, 1)
	params[0] = sourceCode

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_compileSolidity", params)

	if err != nil {
		return "", err
	}

	// TODO - Convert to object
	return pointer.ToString()
}

// Reference: https://eth.wiki/json-rpc/API#eth_gettransactionreceipt
func (eth *Eth) GetTransactionReceipt(hash string) (*dto.TransactionReceipt, error) {
	params := make([]string, 1)
	params[0] = hash

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getTransactionReceipt", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionReceipt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_getblockbynumber
func (eth *Eth) GetBlockByNumber(blockNumber int64, transactionDetails bool) (*dto.Block, error) {
	params := make([]interface{}, 2)
	params[0] = fmt.Sprintf("0x%x", blockNumber)
	params[1] = transactionDetails

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getBlockByNumber", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBlock()
}

// Reference: https://eth.wiki/json-rpc/API#eth_getblocktransactioncountbyhash
func (eth *Eth) GetBlockTransactionCountByHash(hash string) (int64, error) {
	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getBlockTransactionCountByHash", []string{hash})

	if err != nil {
		return -1, err
	}

	return pointer.ToInt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_getblocktransactioncountbynumber
func (eth *Eth) GetBlockTransactionCountByNumber(defaultBlockParameter string) (int64, error) {
	params := make([]string, 1)
	params[0] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getBlockTransactionCountByNumber", params)

	if err != nil {
		return -1, err
	}

	return pointer.ToInt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_getblockbyhash
func (eth *Eth) GetBlockByHash(hash string, transactionDetails bool) (*dto.Block, error) {
	params := make([]interface{}, 2)
	params[0] = hash
	params[1] = transactionDetails

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getBlockByHash", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBlock()
}

// Reference: https://eth.wiki/json-rpc/API#eth_getunclecountbyblockhash
func (eth *Eth) GetUncleCountByBlockHash(hash string) (int64, error) {
	params := make([]string, 1)
	params[0] = hash

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getUncleCountByBlockHash", params)

	if err != nil {
		return -1, err
	}

	return pointer.ToInt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_getunclecountbyblocknumber
func (eth *Eth) GetUncleCountByBlockNumber(quantity int64) (int64, error) {
	params := make([]string, 1)
	params[0] = fmt.Sprintf("0x%x", quantity)

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getUncleCountByBlockNumber", params)

	if err != nil {
		return -1, err
	}

	return pointer.ToInt()
}

// Reference: https://eth.wiki/json-rpc/API#eth_getcode
func (eth *Eth) GetCode(address string, defaultBlockParameter string) (string, error) {

	params := make([]string, 2)
	params[0] = address
	params[1] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(pointer, "eth_getCode", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}
