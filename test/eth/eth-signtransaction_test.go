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
 * @file eth-signtransaction_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */
package test

import (
	"testing"

	"fmt"
	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/complex/types"
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
	"math/big"
)

func TestEthSignTransaction(t *testing.T) {

	var connection = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

	coinbase, err := connection.Eth.GetCoinbase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.Nonce = big.NewInt(5)
	transaction.From = coinbase
	transaction.To = coinbase
	transaction.Value = big.NewInt(0).Mul(big.NewInt(500), big.NewInt(1E18))
	transaction.Gas = big.NewInt(40000)
	transaction.GasPrice = big.NewInt(1E9)
	transaction.Data = types.ComplexString("p2p transaction")

	txID, err := connection.Eth.SignTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txID.Transaction.Nonce.Cmp(transaction.Nonce) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", 5, txID.Transaction.Nonce.Uint64()))
		t.FailNow()
	}

	if txID.Transaction.To != coinbase {
		t.Errorf(fmt.Sprintf("Expected %s | Got: %s", coinbase, txID.Transaction.To))
		t.FailNow()
	}

	if txID.Transaction.Value.Cmp(transaction.Value) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.Value.Uint64(), txID.Transaction.Value.Uint64()))
		t.FailNow()
	}

	if txID.Transaction.Gas.Cmp(transaction.Gas) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.Gas.Uint64(), txID.Transaction.Gas.Uint64()))
		t.FailNow()
	}
	if txID.Transaction.GasPrice.Cmp(transaction.GasPrice) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.GasPrice.Uint64(), txID.Transaction.GasPrice.Uint64()))
		t.FailNow()
	}
}
