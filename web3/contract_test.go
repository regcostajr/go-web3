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
 * @file contract.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2018
 */

package web3

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/complex/types"
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

func TestABIEncoding(t *testing.T) {
	abi, err := ioutil.ReadFile("../test/resources/contracts/abi_encoding/ABIEncoding_sol_ABIEncoding.abi")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bytecode, err := ioutil.ReadFile("../test/resources/contracts/abi_encoding/ABIEncoding_sol_ABIEncoding.bin")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var connection = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))
	contract, err := connection.Eth.NewContract(string(abi))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	coinbase, err := connection.Eth.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction.From = coinbase
	transaction.Gas = big.NewInt(4700000)
	transaction.Data = types.ComplexString(bytecode)

	argsArray := make([]*big.Int, 0)
	argsArray = append(argsArray, big.NewInt(10), big.NewInt(25))

	hash, err := contract.Deploy(transaction, string(bytecode), argsArray)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		receipt, err = connection.Eth.GetTransactionReceipt(hash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Contract Address: ", receipt.ContractAddress)
	transaction.To = receipt.ContractAddress

	result, err := contract.Call(transaction, "uint_dynamic_array", big.NewInt(0))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	firstDecimal, err := result.ToBigInt()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if firstDecimal.Cmp(big.NewInt(10)) != 0 {
		t.Error("First decimal on deploy is not 10")
	}

	result, err = contract.Call(transaction, "uint_dynamic_array", big.NewInt(1))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	secondDecimal, err := result.ToBigInt()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if secondDecimal.Cmp(big.NewInt(25)) != 0 {
		t.Error("Second decimal on deploy is not 25")
	}

}

func TestERC20Contract(t *testing.T) {

	content, err := ioutil.ReadFile("../test/resources/simple-token.json")

	type TruffleContract struct {
		Abi      string `json:"abi"`
		Bytecode string `json:"bytecode"`
	}

	var unmarshalResponse TruffleContract

	json.Unmarshal(content, &unmarshalResponse)

	var connection = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))
	bytecode := unmarshalResponse.Bytecode
	contract, err := connection.Eth.NewContract(unmarshalResponse.Abi)

	transaction := new(dto.TransactionParameters)
	coinbase, err := connection.Eth.GetCoinbase()
	transaction.From = coinbase
	transaction.Gas = big.NewInt(4000000)

	hash, err := contract.Deploy(transaction, bytecode)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		receipt, err = connection.Eth.GetTransactionReceipt(hash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Contract Address: ", receipt.ContractAddress)

	transaction.To = receipt.ContractAddress

	result, err := contract.Call(transaction, "name")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result != nil && err == nil {
		name, _ := result.ToComplexString()
		if name.ToString() != "SimpleToken" {
			t.Errorf(fmt.Sprintf("Name not expected; [Expected %s | Got %s]", "SimpleToken", name.ToString()))
			t.FailNow()
		}
	}

	result, err = contract.Call(transaction, "symbol")
	if result != nil && err == nil {
		symbol, _ := result.ToComplexString()
		if symbol.ToString() != "SIM" {
			t.Errorf("Symbol not expected")
			t.FailNow()
		}
	}

	result, err = contract.Call(transaction, "decimals")
	if result != nil && err == nil {
		decimals, _ := result.ToBigInt()
		if decimals.Int64() != 18 {
			t.Errorf("Decimals not expected")
			t.FailNow()
		}
	}

	bigInt, _ := new(big.Int).SetString("00000000000000000000000000000000000000000000021e19e0c9bab2400000", 16)

	result, err = contract.Call(transaction, "totalSupply")
	if result != nil && err == nil {
		total, _ := result.ToBigInt()
		if total.Cmp(bigInt) != 0 {
			t.Errorf("Total not expected")
			t.FailNow()
		}
	}

	result, err = contract.Call(transaction, "balanceOf", coinbase)
	if result != nil && err == nil {
		balance, _ := result.ToBigInt()
		if balance.Cmp(bigInt) != 0 {
			t.Errorf("Balance not expected")
			t.FailNow()
		}
	}

	hash, err = contract.Send(transaction, "approve", coinbase, big.NewInt(10))
	if err != nil {
		t.Errorf("Can't send approve transaction")
		t.FailNow()
	}

	t.Log(hash)

	reallyBigInt, _ := big.NewInt(0).SetString("20000000000000000000000000000000000000000000000000000000000000000", 16)
	_, err = contract.Send(transaction, "approve", coinbase, reallyBigInt)
	if err == nil {
		t.Errorf("Can't send approve transaction")
		t.FailNow()
	}

}
