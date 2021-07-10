package test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

var eth *web3.Web3 = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

func TestIsSyncing(t *testing.T) {
	_, err := eth.Eth.IsSyncing()

	if err != nil {
		t.Error("wrong response for non syncing node", err)
		t.Fail()
	}

	// TODO: add a proper mocked request response here
}

func TestGetCoinbase(t *testing.T) {
	coinbase, err := eth.Eth.GetCoinbase()

	if err != nil || !strings.HasPrefix(coinbase, "0x") || len(coinbase) < 40 {
		t.Error("wrong response for get coinbase address", err)
		t.Fail()
	}
}

func TestIsMining(t *testing.T) {
	isMining, err := eth.Eth.IsMining()

	if err != nil || !isMining {
		t.Error("wrong response for mining node", err)
		t.Fail()
	}
}

func TestGetHashrate(t *testing.T) {
	hashrate, err := eth.Eth.GetHashRate()

	if err != nil || hashrate < 0 {
		t.Error("error getting hashrate", err)
		t.Fail()
	}
}

func TestGetGasPrice(t *testing.T) {
	gasPrice, err := eth.Eth.GetGasPrice()

	if err != nil {
		t.Error("error getting node gas price", err)
		t.Fail()
	}

	if gasPrice.Int64() <= 0 {
		t.Error("gas price must be bigger than zero")
		t.Fail()
	}
}

func TestListAccounts(t *testing.T) {
	account, err := eth.Personal.NewAccount("test")

	if err != nil {
		t.Error("error in new account request: ", err)
		t.Fail()
	}

	accounts, err := eth.Eth.ListAccounts()

	if err != nil {
		t.Error("error in list accounts request: ", err)
		t.Fail()
	}

	matched := false
	for _, v := range accounts {
		if v == account {
			matched = true
			break
		}
	}

	if !matched {
		t.Error("new address not found on list accounts request")
		t.Fail()
	}
}

func TestEthSendTransaction(t *testing.T) {
	coinbase, err := eth.Eth.GetCoinbase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = coinbase
	transaction.To = coinbase

	value := big.NewInt(0).Mul(big.NewInt(500), big.NewInt(1e18))

	transaction.Value = value
	transaction.Data = "p2p transaction"

	estimation, err := eth.Eth.EstimateGas(transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction.Gas = estimation

	hash, err := eth.Eth.SendTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tx, err := eth.Eth.GetTransactionByHash(hash)

	if err != nil {
		t.Error("transaction not sent: ", err)
		t.FailNow()
	}

	if value.Cmp(tx.Value) != 0 {
		t.Error("transaction value differs from the value sent")
		t.FailNow()
	}

	if coinbase != tx.To || coinbase != tx.From {
		t.Error("transaction addresses differs from coinbase")
		t.FailNow()
	}

}
