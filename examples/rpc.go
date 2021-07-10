package main

import (
	"fmt"
	"math/big"

	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
	"github.com/cellcycle/go-web3/web3/eth"
)

func main() {
	web3Client := web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

	coinbase, err := web3Client.Eth.GetCoinbase()
	if err != nil {
		fmt.Println(err)
	}

	newAccount, err := web3Client.Personal.NewAccount("test")
	if err != nil {
		fmt.Println(err)
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = coinbase
	transaction.To = newAccount

	value := big.NewInt(0).Mul(big.NewInt(1), big.NewInt(1e18))

	transaction.Value = value

	estimation, err := web3Client.Eth.EstimateGas(transaction)
	if err != nil {
		fmt.Println(err)
	}
	transaction.Gas = estimation

	hash, err := web3Client.Eth.SendTransaction(transaction)
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

	newBalance, err := web3Client.Eth.GetBalance(newAccount, eth.LATEST)
	if newBalance.Cmp(value) == 0 {
		fmt.Printf("transaction: %s, amount: %d, address: %s\n", hash, newBalance.Int64(), newAccount)
	} else {
		fmt.Println("error sending transaction")
	}

}
