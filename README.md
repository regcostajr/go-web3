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
transaction.Data = types.String("p2p transaction")

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

[] blablabla

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
