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
 * @file eth-getcode_test.go
 * @authors:
 *   Junjie Chen <chuckjunjchen@gmail.com>
 * @date 2018
 */

package test

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/cellcycle/go-web3/eth/block"

	web3 "github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

func TestEthGetcode(t *testing.T) {

	content, err := ioutil.ReadFile("../resources/simple-token.json")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	type TruffleContract struct {
		Abi              string `json:"abi"`
		Bytecode         string `json:"bytecode"`
		DeployedBytecode string `json:"deployedBytecode"`
	}

	var unmarshalResponse TruffleContract

	json.Unmarshal(content, &unmarshalResponse)

	var connection = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))
	bytecode := unmarshalResponse.Bytecode
	deployedBytecode := unmarshalResponse.DeployedBytecode

	contract, err := connection.Eth.NewContract(unmarshalResponse.Abi)

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
	transaction.Gas = big.NewInt(4000000)
	hash, err := contract.Deploy(transaction, bytecode, nil)

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

	address := receipt.ContractAddress
	code, err := connection.Eth.GetCode(address, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if deployedBytecode != code {
		t.Error("Contract code not expected")
		t.FailNow()
	}
}
