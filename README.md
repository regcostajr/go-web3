# Ethereum Go Client

[![Build Status](https://travis-ci.org/cellcycle/go-web3.svg?branch=master)](https://travis-ci.org/cellcycle/go-web3)

This is a Ethereum compatible RPC client

## Status

Under active development, use at your own risk

## Why should I use it?

This package aims to make things easier while handling contracts and subscriptions.

You probably want(must) use the Geth RPC client for production: [ethclient](https://github.com/ethereum/go-ethereum/tree/master/ethclient)

## Usage

#### Simple RPC call

```go
web3Client := web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))
coinbase, err := web3Client.Eth.GetCoinbase()
```

#### Contract Call

```go
web3Client := ...
contractAddress := ...

contract, _ := web3Client.Eth.NewContract(abi, contractAddress)

result, _ = contract.Functions("balanceOf", 0).Call(contractAddress)
balance = result.ToDataChunks().DecodeInt(0)
```

## Examples

- [Subscription](examples/subscription.go)
- [RPC](examples/rpc.go)
- [Deploying a contract](examples/deploy.go)

## Installation

```bash
go get -u github.com/cellcycle/go-web3
```

### Requirements and Dependencies

* go ^1.16

This package uses `go mod` to handle the dependencies, just running the following
command should make them available into your environment

```bash
go mod tidy
go mod vendor
```

## Contribute!

#### Before a Pull Request:
- If it is an issue please open an issue first so it can be discussed
- Make sure your implementation have been well tested and you wrote/change a test for it
- Don't change the import path to your github username
- run `go fmt` for all your changes.
- run `go test -v ./...`

#### After a Pull Request:
- Make sure the Travis tests are passing.

## Testing

The tests require a running Geth node using the development mode:

```bash
./geth --dev --ws --ws.origins="*" --rpc --rpcapi admin,debug,web3,eth,txpool,personal,clique,miner,net --mine --allow-insecure-unlock
```

Some tests also require access to Infura so if you need to test them please add the file `infura.conf` to the root folder of this project containing your Infura key

Full test:

```bash
go test -v ./test/...
```

Individual test:
```bash
go test -v test/filename.go
```

## License

Package go-web3 is licensed under the [GPLv3](https://www.gnu.org/licenses/gpl-3.0.en.html) License.
