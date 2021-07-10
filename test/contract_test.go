package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

var contractTest *web3.Web3 = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

func TestABIEncoding(t *testing.T) {
	abi, err := ioutil.ReadFile("../test/resources/contract/abi_encoding/ABIEncoding_sol_ABIEncoding.abi")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bytecode, err := ioutil.ReadFile("../test/resources/contract/abi_encoding/ABIEncoding_sol_ABIEncoding.bin")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	contract, err := contractTest.Eth.NewContract(string(abi))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	coinbase, err := contractTest.Eth.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction.From = coinbase

	argsArray := make([]int64, 0)
	argsArray = append(argsArray, 10, 25)

	transaction, err = contract.PrepareTransaction(transaction, contract.Constructor(0), bytecode, argsArray)
	estimation, err := eth.Eth.EstimateGas(transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction.Gas = estimation

	hash, err := contract.Deploy(transaction, bytecode, argsArray)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		receipt, err = contractTest.Eth.GetTransactionReceipt(hash)
	}

	if err != nil || !receipt.Status {
		t.Error("contract deployment unsuccessful", err)
		t.FailNow()
	}

	transaction.To = receipt.ContractAddress

	result, err := contract.Functions("uint_dynamic_array", 0).Call(transaction, big.NewInt(0))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	firstDecimal, err := result.ToInt()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if firstDecimal != 10 {
		t.Error("First decimal on deploy is not 10")
		t.FailNow()
	}

	result, err = contract.Functions("uint_dynamic_array", 0).Call(transaction, big.NewInt(1))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	secondDecimal, err := result.ToInt()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if secondDecimal != 25 {
		t.Error("Second decimal on deploy is not 25")
		t.FailNow()
	}

	hash, err = contract.Send(transaction, contract.Functions("testString", 0), "string", []string{"string"}, [2]string{"string", "string"})

	receipt = nil
	for receipt == nil {
		receipt, err = contractTest.Eth.GetTransactionReceipt(hash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	log := contractTest.Sha3("TestString(string,string[],string[2])")
	if receipt.Logs[0].Topics[0] != log {
		t.Log("topics differs from sha3 conversion")
		t.Fail()
	}

	byteArray := make([][]byte, 0)
	byteArray = append(byteArray, []byte("L1"))

	byteFixedArray := make([][]byte, 0)
	byteFixedArray = append(byteFixedArray, []byte("L1"))
	byteFixedArray = append(byteFixedArray, []byte("L2"))

	byteFixed := make([][]byte, 0)
	byteFixed = append(byteFixed, []byte("dsacsacdas"))
	byteFixed = append(byteFixed, []byte("dsacsacdas"))

	hash, err = contract.Send(transaction, contract.Functions("testBytes", 0), []byte("bytes"), byteArray, byteFixedArray, byteFixed)

	receipt = nil
	for receipt == nil {
		receipt, err = contractTest.Eth.GetTransactionReceipt(hash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	log = contractTest.Sha3("TestBytes(bytes,bytes[],bytes[2],bytes10[2])")
	if receipt.Logs[0].Topics[0] != log {
		t.Log("topics differs from sha3 conversion")
		t.Fail()
	}
}

func TestERC20Contract(t *testing.T) {

	content, err := ioutil.ReadFile("../test/resources/contract/simple-token.json")

	type TruffleContract struct {
		Abi      string `json:"abi"`
		Bytecode string `json:"bytecode"`
	}

	var unmarshalResponse TruffleContract

	json.Unmarshal(content, &unmarshalResponse)

	bytecode := unmarshalResponse.Bytecode
	contract, err := contractTest.Eth.NewContract(unmarshalResponse.Abi)

	transaction := new(dto.TransactionParameters)
	coinbase, err := contractTest.Eth.GetCoinbase()
	transaction.From = coinbase

	transaction, err = contract.PrepareTransaction(transaction, contract.Constructor(0), []byte(bytecode))
	estimation, err := eth.Eth.EstimateGas(transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction.Gas = estimation

	hash, err := contract.Deploy(transaction, []byte(bytecode))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		receipt, err = contractTest.Eth.GetTransactionReceipt(hash)
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction.To = receipt.ContractAddress

	result, err := contract.Functions("name", 0).Call(transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chunks := result.ToDataChunks()
	name, _ := contract.DecodeString(chunks[2])
	if name != "SimpleToken" {
		t.Errorf(fmt.Sprintf("Name not expected; [Expected %s | Got %s]", "SimpleToken", name))
		t.FailNow()
	}

	result, err = contract.Functions("symbol", 0).Call(transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chunks = result.ToDataChunks()
	symbol, _ := contract.DecodeString(chunks[2])
	if symbol != "SIM" {
		t.Errorf("Symbol not expected")
		t.FailNow()
	}

	result, err = contract.Functions("decimals", 0).Call(transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chunks = result.ToDataChunks()
	decimals := contract.DecodeInt(chunks[0])
	if decimals.Int64() != 18 {
		t.Errorf("Decimals not expected")
		t.FailNow()
	}

	totalSupply := big.NewInt(0).Mul(big.NewInt(10000), big.NewInt(1e18))

	result, err = contract.Functions("totalSupply", 0).Call(transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chunks = result.ToDataChunks()
	total := contract.DecodeInt(chunks[0])
	if total.Cmp(totalSupply) != 0 {
		t.Errorf("Total not expected")
		t.FailNow()
	}

	result, err = contract.Functions("balanceOf", 0).Call(transaction, coinbase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chunks = result.ToDataChunks()
	balance := contract.DecodeInt(chunks[0])
	if balance.Cmp(totalSupply) != 0 {
		t.Errorf("Balance not expected")
		t.FailNow()
	}

	newAddress, err := contractTest.Personal.NewAccount("test")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	approveBalance := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18))
	hash, err = contract.Functions("approve", 0).Send(transaction, newAddress, approveBalance)
	if err != nil {
		t.Errorf("Can't send approve transaction")
		t.FailNow()
	}

	receipt = nil
	for receipt == nil {
		receipt, err = contractTest.Eth.GetTransactionReceipt(hash)
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	log := contractTest.Sha3("Approval(address,address,uint256)")
	if receipt.Logs[0].Topics[0] != log {
		t.Log("topics differs from sha3 conversion")
		t.Fail()
	}

	if contract.DecodeAddress(receipt.Logs[0].Topics[1]) != coinbase {
		t.Log("owner address differs from the topic")
		t.Fail()
	}

	if contract.DecodeAddress(receipt.Logs[0].Topics[2]) != newAddress {
		t.Log("new address differs from the topic")
		t.Fail()
	}

	result, err = contract.Functions("allowance", 0).Call(transaction, coinbase, newAddress)
	chunks = result.ToDataChunks()
	allowance := contract.DecodeInt(chunks[0])
	if allowance.Cmp(approveBalance) != 0 {
		t.Error("Allowance not expected")
		t.FailNow()
	}

	ethTransfer := new(dto.TransactionParameters)
	ethTransfer.From = coinbase
	ethTransfer.To = newAddress

	ethTransfer.Value = big.NewInt(0).Mul(big.NewInt(1), big.NewInt(1e18))
	ethTransfer.Gas = big.NewInt(40000)

	hash, err = contractTest.Eth.SendTransaction(ethTransfer)
	if err != nil {
		t.Error("Can't send fee for the new address", err)
		t.FailNow()
	}

	receipt = nil
	for receipt == nil {
		receipt, err = contractTest.Eth.GetTransactionReceipt(hash)
	}

	transaction.From = newAddress
	sucess, err := contractTest.Personal.UnlockAccount(newAddress, "test", 100)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !sucess {
		t.Error("can't unlock new address")
		t.FailNow()
	}

	hash, err = contract.Functions("transferFrom", 0).Send(transaction, coinbase, newAddress, approveBalance)
	if err != nil {
		t.Error("Can't send transfer from transaction", err)
		t.FailNow()
	}

	receipt = nil
	for receipt == nil {
		receipt, err = contractTest.Eth.GetTransactionReceipt(hash)
	}

	result, err = contract.Functions("balanceOf", 0).Call(transaction, newAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chunks = result.ToDataChunks()
	balance = contract.DecodeInt(chunks[0])
	if balance.Cmp(approveBalance) != 0 {
		t.Errorf("Balance not expected")
		t.FailNow()
	}
}
