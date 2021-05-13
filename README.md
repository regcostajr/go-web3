# Ethereum Go Client

[![Build Status](https://travis-ci.org/cellcycle/go-web3.svg?branch=master)](https://travis-ci.org/cellcycle/go-web3)

This is a Ethereum compatible Go Client

## Status

This package is currently under active development. It is not yet stable and there are some RPC methods left to implement (some to fix) and documentation to be done.

## Usage

#### Deploying a contract

```go

bytecode := ... #contract bytecode
abi := ... #contract abi

var connection = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))
contract, err := connection.Eth.NewContract(abi)

transaction := new(dto.TransactionParameters)
coinbase, err := connection.Eth.GetCoinbase()
transaction.From = coinbase
transaction.Gas = big.NewInt(4000000)

hash, err := contract.Deploy(transaction, bytecode, nil)

fmt.Println(hash)

```

#### Using contract public functions

```go

result, err = contract.Call(transaction, "balanceOf", coinbase)
if result != nil && err == nil {
	balance, _ := result.ToComplexIntResponse()
	fmt.Println(balance.ToBigInt())
}

```

#### Using contract payable functions

```go

hash, err = contract.Send(transaction, "approve", coinbase, 10)

```

#### Using RPC commands

GetBalance

```go

balance, err := connection.Eth.GetBalance(coinbase, block.LATEST)

```

SendTransaction

```go

transaction := new(dto.TransactionParameters)
transaction.From = coinbase
transaction.To = coinbase
transaction.Value = big.NewInt(10)
transaction.Gas = big.NewInt(40000)
transaction.Data = types.ComplexString("p2p transaction")

txID, err := connection.Eth.SendTransaction(transaction)

```


## Contribute!

#### Before a Pull Request:
- Create at least one test for your implementation.
- Don't change the import path to your github username.
- run `go fmt` for all your changes.
- run `go test -v ./...`

#### After a Pull Request:
- Please use the travis log if an error occurs.

### In Progress = ![](https://placehold.it/15/FFFF00/000000?text=+)
### Partially implemented = ![](https://placehold.it/15/008080/000000?text=+)

TODO List

- [x] web3_clientVersion
- [x] web3_sha3
- [x] net_version
- [x] net_peerCount
- [x] net_listening
- [x] eth_syncing
- [x] eth_coinbase
- [x] eth_mining
- [x] eth_hashrate
- [x] eth_gasPrice
- [x] eth_accounts
- [x] eth_blockNumber
- [x] eth_getBalance
- [x] eth_getStorageAt (deprecated)
- [x] eth_getTransactionCount
- [x] eth_getBlockTransactionCountByHash
- [x] eth_getBlockTransactionCountByNumber
- [x] eth_getUncleCountByBlockHash
- [x] eth_getUncleCountByBlockNumber
- [x] eth_getCode
- [x] eth_sign
- [x] eth_sendTransaction
- [ ] eth_sendRawTransaction
- [x] eth_call
- [x] eth_estimateGas
- [x] eth_getBlockByHash
- [x] eth_getBlockByNumber
- [x] eth_getTransactionByHash
- [x] eth_getTransactionByBlockHashAndIndex
- [x] eth_getTransactionByBlockNumberAndIndex
- [x] eth_getTransactionReceipt
- [ ] eth_getUncleByBlockHashAndIndex
- [ ] eth_getUncleByBlockNumberAndIndex
- [ ] eth_getCompilers
- [ ] eth_compileLLL
- [x] eth_compileSolidity (deprecated)
- [ ] eth_compileSerpent
- [ ] eth_newFilter
- [ ] eth_newBlockFilter
- [ ] eth_newPendingTransactionFilter
- [ ] eth_uninstallFilter
- [ ] eth_getFilterChanges
- [ ] eth_getFilterLogs
- [ ] eth_getLogs
- [ ] eth_getWork
- [ ] eth_submitWork
- [ ] eth_submitHashrate
- [ ] db_putString
- [ ] db_getString
- [ ] db_putHex
- [ ] db_getHex
- [ ] shh_post
- [ ] shh_version
- [ ] shh_newIdentity
- [ ] shh_hasIdentity
- [ ] shh_newGroup
- [ ] shh_addToGroup
- [ ] shh_newFilter
- [ ] shh_uninstallFilter
- [ ] shh_getFilterChanges
- [ ] shh_getMessages
- [x] personal_listAccounts
- [x] personal_newAccount
- [x] personal_sendTransaction
- [x] personal_unlockAccount

## Installation

### go get

```bash
go get -u github.com/cellcycle/go-web3
```

### glide

```bash
glide get github.com/cellcycle/go-web3
```

### Requirements

* go ^1.8.3
* golang.org/x/net

## Testing

Node running in dev mode:

```bash
geth --dev --shh --ws --wsorigins="*" --rpc --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3 --mine
```

Full test:

```bash
go test -v ./test/...
```

Individual test:
```bash
go test -v test/modulename/filename.go
```

## License

Package go-web3 is licensed under the [GPLv3](https://www.gnu.org/licenses/gpl-3.0.en.html) License.
