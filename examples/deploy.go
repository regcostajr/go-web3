package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

func main() {
	// Just importing the contract bytecode and ABI
	content, err := ioutil.ReadFile("test/resources/contract/simple-token.json")

	type TruffleContract struct {
		Abi      string `json:"abi"`
		Bytecode string `json:"bytecode"`
	}

	var unmarshalResponse TruffleContract

	json.Unmarshal(content, &unmarshalResponse)

	bytecode := unmarshalResponse.Bytecode

	// deployment starts here
	web3Client := web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))
	contract, err := web3Client.Eth.NewContract(unmarshalResponse.Abi, "")

	transaction := new(dto.TransactionParameters)
	coinbase, err := web3Client.Eth.GetCoinbase()
	transaction.From = coinbase

	// this is required to estimate the transaction gas
	// if you are not going to estimate it you can use it
	// fixed as: transaction.Gas = big.NewInt(400000)
	transaction, err = contract.PrepareTransaction(transaction, contract.Constructor(0), []byte(bytecode))
	estimation, err := web3Client.Eth.EstimateGas(transaction)
	if err != nil {
		fmt.Println(err)
	}
	transaction.Gas = estimation

	hash, err := contract.Deploy(transaction, []byte(bytecode))

	if err != nil {
		fmt.Println(err)
	}

	var receipt *dto.TransactionReceipt

	// keep trying until get the receipt, you most probably
	// DON'T want to do this in production
	for receipt == nil {
		receipt, err = web3Client.Eth.GetTransactionReceipt(hash)
	}
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(receipt.ContractAddress)

}
